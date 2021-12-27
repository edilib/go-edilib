package dom

import (
	"fmt"
	"github.com/edilib/go-edi/dom/types"
	"github.com/edilib/go-edi/stream"
	streamTypes "github.com/edilib/go-edi/stream/types"
	"io"
)

type Parser struct {
	rd stream.EDITokenReader
}

func NewParser(reader io.Reader) *Parser {
	return &Parser{rd: stream.NewEDITokenReader(reader)}
}

func (p *Parser) ReadAll() ([]types.EDISegment, error) {
	var segments []types.EDISegment
	for {
		segment, err := p.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}
		segments = append(segments, segment)
	}

	return segments, nil
}

func (p *Parser) Read() (types.EDISegment, error) {

	var segment *types.EDISegment
	var composite *types.EDIComposite
	for {
		token, err := p.rd.Read()
		if err != nil {
			return types.EDISegment{}, err
		}

		if token.Type() == streamTypes.UNA_SEGMENT {
			return types.EDISegment{Tag: "UNA"}, nil
		} else if token.Type() == streamTypes.SEGMENT_TAG {
			segment = &types.EDISegment{Tag: token.Value(), Elements: []interface{}{}}
		} else if token.Type() == streamTypes.SEGMENT_TERMINATOR {
			return *segment, nil
		} else if token.Type() == streamTypes.DATA_ELEMENT_SEPERATOR {
			composite = nil
		} else if token.Type() == streamTypes.COMPONENT_DATA_ELEMENT_SEPERATOR {
			if composite == nil {
				composite = &types.EDIComposite{Elements: []interface{}{}}
				composite.Elements = append(composite.Elements, segment.Elements[len(segment.Elements)-1])
				segment.Elements[len(segment.Elements)-1] = composite
			}
		} else if token.Type() == streamTypes.REPETITION_SEPERATOR {
			return types.EDISegment{}, fmt.Errorf("repetition unsupported")
		} else if token.Type() == streamTypes.VALUE {
			if composite == nil {
				segment.Elements = append(segment.Elements, token.Value())
			} else {
				composite.Elements = append(composite.Elements, token.Value())
			}
		} else if token.Type() == streamTypes.EOF {
			return types.EDISegment{}, io.EOF
		} else {
			return types.EDISegment{}, fmt.Errorf("invalid input %s", token.Type())
		}
	}

}
