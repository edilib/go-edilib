package internal

import (
	"fmt"
	"github.com/edilib/go-edilib/edifact/types"
	"github.com/edilib/go-edilib/internal/scanner"
	"io"
)

const (
	INITIAL          = 0
	INITIAL_UNA_SEEN = 1
	IN_MESSAGE       = 2
	FINISHED         = 3

	IN_SIMPLE_VALUE     = 0
	IN_COMPOSITE_VALUE  = 1
	IN_REPETITION_VALUE = 2
)

type SegmentReader struct {
	scanner  *scanner.Scanner
	segments []types.Segment
	state    int
}

func NewSegmentReader(reader io.Reader, format types.Format) *SegmentReader {
	segmentReader := SegmentReader{scanner: scanner.NewScanner(reader, format), segments: []types.Segment{}, state: INITIAL}
	return &segmentReader
}

func (r *SegmentReader) HasNext() (bool, error) {
	err := r.fill()
	if err != nil {
		return false, err
	}

	return len(r.segments) > 0, nil
}

func (r *SegmentReader) Next() (types.Segment, error) {
	err := r.fill()
	if err != nil {
		return types.Segment{}, err
	}

	if len(r.segments) == 0 {
		return types.Segment{}, fmt.Errorf("no more segments")
	}

	next := r.segments[0]
	r.segments = r.segments[1:]
	return next, nil
}

func (r *SegmentReader) Close() error {
	if r.scanner != nil {
		err := r.scanner.Close()
		r.scanner = nil
		return err
	}

	return nil
}

func (r *SegmentReader) ReadAll() ([]types.Segment, error) {

	var segments []types.Segment
	for {
		hasNext, err := r.HasNext()
		if err != nil {
			return nil, err
		}

		if hasNext {
			token, err := r.Next()
			if err != nil {
				return nil, err
			}

			segments = append(segments, token)
		} else {
			break
		}
	}

	return segments, nil
}

func (r *SegmentReader) fill() error {
	if len(r.segments) > 0 {
		return nil
	}

	for {
		switch r.state {
		case INITIAL:
			err := r.processInitialState()
			if err != nil {
				return err
			}
		case INITIAL_UNA_SEEN:
			err := r.processInitialUnaSeenState()
			if err != nil {
				return err
			}
		case IN_MESSAGE:
			err := r.processInMessageState()
			if err != nil {
				return err
			}
		case FINISHED:
			return nil
		default:
			return fmt.Errorf("illegal state: %d", r.state)
		}
	}
}

func (r *SegmentReader) processInitialState() error {
	token, err := r.scanner.Peek()
	if err != nil {
		return err
	}

	switch token.Type() {
	case scanner.UNA_SEGMENT:
		err := r.readUnaSegment()
		if err != nil {
			return err
		}
		r.state = INITIAL_UNA_SEEN
		return nil
	case scanner.VALUE:
		r.state = IN_MESSAGE
		return nil
	case scanner.EOF:
		return fmt.Errorf("empty file")
	default:
		return r.unexpectedInput(token)
	}
}

func (r *SegmentReader) unexpectedInput(token scanner.ScannerToken) error {
	return fmt.Errorf("unexpected input %s, value=%s at %s", token.Type().Name(), token.Value(), token.Pos())
}

func (r *SegmentReader) readUnaSegment() error {
	_, err := r.scanner.Consume(scanner.UNA_SEGMENT)
	return err
}

func (r *SegmentReader) processInitialUnaSeenState() error {
	token, err := r.scanner.Peek()
	if err != nil {
		return err
	}

	switch token.Type() {
	case scanner.UNA_SEGMENT:
		return fmt.Errorf("duplicate una segment at %s", token.Pos())
	case scanner.VALUE:
		r.state = IN_MESSAGE
		return nil
	case scanner.EOF:
		return fmt.Errorf("no segments after una segment found")
	default:
		return r.unexpectedInput(token)
	}
}

func (r *SegmentReader) processInMessageState() error {
	token, err := r.scanner.Peek()
	if err != nil {
		return err
	}

	switch token.Type() {
	case scanner.VALUE:
		err := r.readSegment()
		if err != nil {
			return err
		}
		return nil
	case scanner.EOF:
		r.state = FINISHED
		return nil
	default:
		return r.unexpectedInput(token)
	}
}

func (r *SegmentReader) readSegment() error {
	token, err := r.scanner.Peek()
	if err != nil {
		return err
	}

	if token.Type() == scanner.EOF {
		return nil
	}

	tag, err := r.readTag()
	if err != nil {
		return err
	}

	dataElements, err := r.readDataElements()
	if err != nil {
		return err
	}

	err = r.readSegmentEnd()
	if err != nil {
		return err
	}

	r.segments = append(r.segments, types.Segment{Tag: tag, DataElements: dataElements})

	return nil
}

func (r *SegmentReader) readSegmentEnd() error {
	_, err := r.scanner.Consume(scanner.SEGMENT_TERMINATOR)
	return err

}

func (r *SegmentReader) readTag() (types.Tag, error) {
	values := []types.SimpleValue{}

	token, err := r.scanner.Consume(scanner.VALUE)
	if err != nil {
		return types.Tag{}, err
	}
	values = append(values, types.SimpleValue{StringValue: token.Value()})
	valueSeen := false

	token, err = r.scanner.Peek()
	if err != nil {
		return types.Tag{}, err
	}

	if token.Type() == scanner.COMPONENT_DATA_ELEMENT_SEPERATOR {
		_, err := r.scanner.Consume(token.Type())
		if err != nil {
			return types.Tag{}, err
		}

		for {
			token, err = r.scanner.Peek()
			if err != nil {
				return types.Tag{}, err
			}

			switch token.Type() {
			case scanner.VALUE:
				_, err := r.scanner.Consume(token.Type())
				if err != nil {
					return types.Tag{}, err
				}

				values = append(values, types.SimpleValue{StringValue: token.String()})
				valueSeen = true
			case scanner.DATA_ELEMENT_SEPERATOR, scanner.COMPONENT_DATA_ELEMENT_SEPERATOR:
				if !valueSeen {
					values = append(values, types.SimpleValue{StringValue: ""})
				}
				_, err := r.scanner.Consume(token.Type())
				if err != nil {
					return types.Tag{}, err
				}
				valueSeen = false
			case scanner.SEGMENT_TERMINATOR:
				if !valueSeen {
					values = append(values, types.SimpleValue{StringValue: ""})
				}
			default:
				return types.Tag{}, r.unexpectedInput(token)
			}
		}
	}

	return types.Tag{Values: types.CompositeValue{Values: values}}, nil
}

func (r *SegmentReader) readDataElements() ([]types.Value, error) {

	values := []types.Value{}

	for {
		token, err := r.scanner.Peek()
		if err != nil {
			return nil, err
		}

		switch token.Type() {
		case scanner.SEGMENT_TERMINATOR:
			return values, nil
		case scanner.DATA_ELEMENT_SEPERATOR:
			value, err := r.readDataElement()
			if err != nil {
				return nil, err
			}

			values = append(values, value)
		default:
			return nil, r.unexpectedInput(token)
		}
	}
}

func (r *SegmentReader) readDataElement() (types.Value, error) {
	values := []types.SimpleValue{}

	_, err := r.scanner.Consume(scanner.DATA_ELEMENT_SEPERATOR)
	if err != nil {
		return nil, err
	}

	state := IN_SIMPLE_VALUE
	valueSeen := false
	for {
		token, err := r.scanner.Peek()
		if err != nil {
			return nil, err
		}

		if token.Type() == scanner.VALUE {
			_, err := r.scanner.Consume(token.Type())
			if err != nil {
				return nil, err
			}

			values = append(values, types.SimpleValue{StringValue: token.Value()})
			valueSeen = true
		} else if (state == IN_SIMPLE_VALUE || state == IN_COMPOSITE_VALUE) && token.Type() == scanner.COMPONENT_DATA_ELEMENT_SEPERATOR {
			if !valueSeen {
				values = append(values, types.SimpleValue{StringValue: ""})
			}
			state = IN_COMPOSITE_VALUE
			valueSeen = false
			_, err := r.scanner.Consume(token.Type())
			if err != nil {
				return nil, err
			}
		} else if (state == IN_SIMPLE_VALUE || state == IN_REPETITION_VALUE) && token.Type() == scanner.REPETITION_SEPERATOR {
			if !valueSeen {
				values = append(values, types.SimpleValue{StringValue: ""})
			}
			state = IN_REPETITION_VALUE
			valueSeen = false
			_, err := r.scanner.Consume(token.Type())
			if err != nil {
				return nil, err
			}
		} else if token.Type() == scanner.DATA_ELEMENT_SEPERATOR || token.Type() == scanner.SEGMENT_TERMINATOR {
			if !valueSeen {
				values = append(values, types.SimpleValue{StringValue: ""})
			}
			switch state {
			case IN_SIMPLE_VALUE:
				return values[0], nil
			case IN_COMPOSITE_VALUE:
				return types.CompositeValue{Values: values}, nil
			case IN_REPETITION_VALUE:
				return types.RepetitionValue{Values: values}, nil
			default:
				return nil, fmt.Errorf("illegal state")
			}
		} else {
			return nil, r.unexpectedInput(token)
		}

	}
}
