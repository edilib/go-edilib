package scanner

import (
	"fmt"
	"io"
)

type state int

const (
	INITIAL state = iota
	INITIAL_DATA_SEEN
	IN_VALUE
	IN_VALUE_RELEASE_SEEN
)

type Scanner struct {
	state  state
	scanRd *ScannerReader
	tokens []ScannerToken
	format Format
}

func NewScanner(reader io.Reader) *Scanner {
	return &Scanner{state: INITIAL, scanRd: NewScannerReader(reader), tokens: []ScannerToken{}, format: DefaultFormat()}
}

func (s *Scanner) Follows(tType EDITokenType) (bool, error) {
	token, err := s.Peek()
	if err != nil {
		return false, err
	}

	return token.tType == tType, nil
}

func (s *Scanner) Peek() (ScannerToken, error) {
	err := s.fill()
	if err != nil {
		return ScannerToken{}, err
	}

	token := s.tokens[0]

	return token, nil
}

func (s *Scanner) All() ([]ScannerToken, error) {
	var tokens []ScannerToken
	for {
		token, err := s.Next()
		if err != nil {
			return nil, err
		}

		tokens = append(tokens, token)

		if token.tType == EOF {
			break
		}
	}

	return tokens, nil
}

func (s *Scanner) Consume(tType EDITokenType) (ScannerToken, error) {
	token, err := s.Next()
	if err != nil {
		return ScannerToken{}, err
	}

	if token.tType != tType {
		return ScannerToken{}, fmt.Errorf("expected %s but was %v", tType.Name(), token.String())
	}

	return token, nil
}

func (s *Scanner) Next() (ScannerToken, error) {
	err := s.fill()
	if err != nil {
		return ScannerToken{}, err
	}

	token := s.tokens[0]
	s.tokens = s.tokens[1:]

	return token, nil
}

func (s *Scanner) follows(str []rune) (bool, error) {

	for i := 0; i < len(str); i++ {
		r, err := s.scanRd.PeekRune(i)
		if r != str[i] || err == io.EOF {
			return false, nil
		}

		if err != nil {
			return false, err
		}
	}

	return true, nil
}

func (s *Scanner) fill() error {

	pos := s.scanRd.Position()
	if s.state == INITIAL {
		unaFollows, err := s.follows([]rune("UNA"))
		if err != nil {
			return err
		}
		if unaFollows {
			value, err := s.scanRd.ReadNRunes(9)
			if err != nil {
				return err
			}

			s.tokens = append(s.tokens, ScannerToken{tType: UNA_SEGMENT, value: string(value), pos: pos, err: err})
			s.format = Format{componentDataElementSeperator: value[3], dataElementSeperator: value[4], decimalMark: value[5], releaseCharacter: value[6], repetitionSeperator: value[7], segmentTerminator: value[8]}
			s.state = INITIAL_DATA_SEEN
			return nil
		}
	}

	var buf []rune
	for {
		b, err := s.scanRd.PeekRune(0)
		if err != nil && err != io.EOF {
			s.tokens = append(s.tokens, ScannerToken{tType: ERROR, value: "", pos: pos, err: err})
			return nil
		}

		switch s.state {
		case INITIAL, INITIAL_DATA_SEEN:
			if err == io.EOF {
				s.tokens = append(s.tokens, ScannerToken{tType: EOF, value: "", pos: pos, err: nil})
				return nil
			} else if b == s.format.releaseCharacter {
				_, _ = s.scanRd.ReadRune()
				s.state = IN_VALUE_RELEASE_SEEN
			} else if b == s.format.segmentTerminator {
				b, _ := s.scanRd.ReadRune()
				s.tokens = append(s.tokens, ScannerToken{tType: SEGMENT_TERMINATOR, value: string(b), pos: pos, err: nil})
				return nil
			} else if b == s.format.componentDataElementSeperator {
				b, _ := s.scanRd.ReadRune()
				s.tokens = append(s.tokens, ScannerToken{tType: COMPONENT_DATA_ELEMENT_SEPERATOR, value: string(b), pos: pos, err: nil})
				return nil
			} else if b == s.format.dataElementSeperator {
				b, _ := s.scanRd.ReadRune()
				s.tokens = append(s.tokens, ScannerToken{tType: DATA_ELEMENT_SEPERATOR, value: string(b), pos: pos, err: nil})
				return nil
			} else if b == s.format.repetitionSeperator {
				b, _ := s.scanRd.ReadRune()
				s.tokens = append(s.tokens, ScannerToken{tType: REPETITION_SEPERATOR, value: string(b), pos: pos, err: fmt.Errorf("eof after release char")})
				return nil
			} else {
				b, _ := s.scanRd.ReadRune()
				s.state = IN_VALUE
				buf = append(buf, b)
			}
		case IN_VALUE:
			if err == io.EOF ||
				b == s.format.dataElementSeperator ||
				b == s.format.componentDataElementSeperator ||
				b == s.format.segmentTerminator ||
				b == s.format.repetitionSeperator {
				s.state = INITIAL
				s.tokens = append(s.tokens, ScannerToken{tType: VALUE, value: string(buf), pos: pos, err: nil})
				return nil
			} else if b == s.format.releaseCharacter {
				_, _ = s.scanRd.ReadRune()
				s.state = IN_VALUE_RELEASE_SEEN
			} else {
				_, _ = s.scanRd.ReadRune()
				buf = append(buf, b)
			}
		case IN_VALUE_RELEASE_SEEN:
			if err == io.EOF {
				s.tokens = append(s.tokens, ScannerToken{tType: ERROR, value: "", pos: pos, err: fmt.Errorf("eof after release char")})
				return nil
			} else {
				_, _ = s.scanRd.ReadRune()
				s.state = IN_VALUE
				buf = append(buf, b)
			}
		default:
			return fmt.Errorf("invalid state")
		}
	}
}

func (s *Scanner) Close() error {
	return s.Close()
}
