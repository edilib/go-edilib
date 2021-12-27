package stream

import (
	"fmt"
	"github.com/edilib/go-edi/stream/types"
	"io"
)

type ParserState string

const (
	PS_INITIAL      ParserState = "INITIAL"
	PS_SEGMENT_SEEN ParserState = "SEGMENT_SEEN"
	PS_IN_SEGMENT   ParserState = "IN_SEGMENT"
)

type Parser struct {
	scanner *Scanner
	state   ParserState
	tokens  []ScannerToken
}

func NewParser(reader io.Reader) *Parser {
	return &Parser{scanner: NewScanner(reader), state: PS_INITIAL, tokens: []ScannerToken{}}
}

func (p *Parser) Peek() (types.EDIToken, error) {
	err := p.tryFill()
	if err != nil {
		return nil, err
	}

	token := p.tokens[0]

	return types.EDIToken(&token), nil
}

func (p *Parser) Read() (types.EDIToken, error) {
	err := p.tryFill()
	if err != nil {
		return nil, err
	}

	token := p.tokens[0]
	p.tokens = p.tokens[1:]

	return types.EDIToken(&token), nil
}

func (p *Parser) tryFill() error {
	next, err := p.scanner.Next()
	if err != nil {
		return err
	}

	switch p.state {
	case PS_INITIAL:
		if next.Type() == types.EOF {
			next = ScannerToken{tType: types.ERROR, value: next.value, pos: next.pos, err: io.EOF}
		} else if next.Type() == types.UNA_SEGMENT {
			p.state = PS_SEGMENT_SEEN
		} else if next.Type() == types.VALUE {
			next = ScannerToken{tType: types.SEGMENT_TAG, value: next.value, pos: next.pos, err: nil}
			p.state = PS_IN_SEGMENT
		} else {
			next = ScannerToken{tType: types.ERROR, value: next.value, pos: next.pos, err: fmt.Errorf("unexpected input: %s", next.tType)}
		}
	case PS_SEGMENT_SEEN:
		if next.Type() == types.EOF {
			next = ScannerToken{tType: types.EOF, value: next.value, pos: next.pos, err: nil}
		} else if next.Type() == types.UNA_SEGMENT {
			next = ScannerToken{tType: types.ERROR, value: next.value, pos: next.pos, err: fmt.Errorf("una not allowed in between: %s", next.tType)}
		} else if next.Type() == types.VALUE {
			next = ScannerToken{tType: types.SEGMENT_TAG, value: next.value, pos: next.pos, err: nil}
			p.state = PS_IN_SEGMENT
		} else {
			next = ScannerToken{tType: types.ERROR, value: next.value, pos: next.pos, err: fmt.Errorf("unexpected input: %s", next.tType)}
		}
	case PS_IN_SEGMENT:
		if next.Type() == types.EOF {
			next = ScannerToken{tType: types.ERROR, value: next.value, pos: next.pos, err: fmt.Errorf("eof in segment")}
		} else if next.Type() == types.UNA_SEGMENT {
			next = ScannerToken{tType: types.ERROR, value: next.value, pos: next.pos, err: fmt.Errorf("una not allowed in between: %s", next.tType)}
		} else if next.Type() == types.SEGMENT_TAG {
			next = ScannerToken{tType: types.ERROR, value: next.value, pos: next.pos, err: fmt.Errorf("segment tag not allowed in between: %s", next.tType)}
		} else if next.Type() == types.SEGMENT_TERMINATOR {
			p.state = PS_SEGMENT_SEEN
		} else {
			// ok
		}
	default:
		return fmt.Errorf("invalid state")
	}

	p.tokens = append(p.tokens, next)

	return nil
}

func (s *Parser) ReadAll() ([]types.EDIToken, error) {
	var tokens []types.EDIToken
	for {
		token, err := s.Read()
		if err != nil {
			return nil, err
		}

		tokens = append(tokens, token)

		if token.Type() == types.EOF || token.Type() == types.ERROR {
			break
		}
	}

	return tokens, nil
}

func (p *Parser) Close() error {
	return p.scanner.Close()
}
