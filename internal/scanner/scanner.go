package scanner

import (
	"fmt"
	"github.com/edilib/go-edilib/edifact/types"
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
	format types.Format
}

func NewScanner(reader io.Reader, format types.Format) *Scanner {
	return &Scanner{state: INITIAL, scanRd: NewScannerReader(reader), tokens: []ScannerToken{}, format: format}
}

func (s *Scanner) Follows(tType ScannerTokenType) (bool, error) {
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

func (s *Scanner) Consume(tType ScannerTokenType) (ScannerToken, error) {
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

	line, column := s.scanRd.Position()
	if s.state == INITIAL {
		unaFollows, err := s.follows([]rune("UNA"))
		if err != nil {
			return err
		}
		if unaFollows {
			if !s.format.UnaAllowed {
				return fmt.Errorf("una segment not allowed at %d:%d", line, column)
			}

			value, err := s.scanRd.ReadNRunes(9)
			if err != nil {
				return err
			}

			s.tokens = append(s.tokens, ScannerToken{tType: UNA_SEGMENT, value: string(value), line: line, column: column, err: err})
			s.format = types.Format{SkipNewLineAfterSegment: s.format.SkipNewLineAfterSegment, UnaAllowed: s.format.UnaAllowed, ComponentDataElementSeperator: value[3], DataElementSeperator: value[4], DecimalMark: value[5], ReleaseCharacter: value[6], RepetitionSeperator: value[7], SegmentTerminator: value[8]}
			s.state = INITIAL_DATA_SEEN
			return nil
		}
	}

	var buf []rune
	for {
		b, err := s.scanRd.PeekRune(0)
		if err != nil && err != io.EOF {
			s.tokens = append(s.tokens, ScannerToken{tType: ERROR, value: "", line: line, column: column, err: err})
			return nil
		}

		switch s.state {
		case INITIAL, INITIAL_DATA_SEEN:
			if err == io.EOF {
				s.tokens = append(s.tokens, ScannerToken{tType: EOF, value: "", line: line, column: column, err: nil})
				return nil
			} else if b == s.format.ReleaseCharacter {
				_, _ = s.scanRd.ReadRune()
				s.state = IN_VALUE_RELEASE_SEEN
			} else if b == s.format.SegmentTerminator {
				b, _ := s.scanRd.ReadRune()
				s.tokens = append(s.tokens, ScannerToken{tType: SEGMENT_TERMINATOR, value: string(b), line: line, column: column, err: nil})

				if s.format.SkipNewLineAfterSegment {
					next, err := s.scanRd.PeekRune(0)
					if err != nil {
						return err
					}
					if next == '\n' {
						_, err := s.scanRd.ReadRune()
						if err != nil {
							return err
						}
					}
				}
				return nil
			} else if b == s.format.ComponentDataElementSeperator {
				b, _ := s.scanRd.ReadRune()
				s.tokens = append(s.tokens, ScannerToken{tType: COMPONENT_DATA_ELEMENT_SEPERATOR, value: string(b), line: line, column: column, err: nil})
				return nil
			} else if b == s.format.DataElementSeperator {
				b, _ := s.scanRd.ReadRune()
				s.tokens = append(s.tokens, ScannerToken{tType: DATA_ELEMENT_SEPERATOR, value: string(b), line: line, column: column, err: nil})
				return nil
			} else if b == s.format.RepetitionSeperator {
				b, _ := s.scanRd.ReadRune()
				s.tokens = append(s.tokens, ScannerToken{tType: REPETITION_SEPERATOR, value: string(b), line: line, column: column, err: fmt.Errorf("eof after release char")})
				return nil
			} else {
				b, _ := s.scanRd.ReadRune()
				s.state = IN_VALUE
				buf = append(buf, b)
			}
		case IN_VALUE:
			if err == io.EOF ||
				b == s.format.DataElementSeperator ||
				b == s.format.ComponentDataElementSeperator ||
				b == s.format.RepetitionSeperator ||
				b == s.format.SegmentTerminator {
				s.state = INITIAL
				s.tokens = append(s.tokens, ScannerToken{tType: VALUE, value: string(buf), line: line, column: column, err: nil})
				return nil
			} else if b == s.format.ReleaseCharacter {
				_, _ = s.scanRd.ReadRune()
				s.state = IN_VALUE_RELEASE_SEEN
			} else {
				_, _ = s.scanRd.ReadRune()
				buf = append(buf, b)
			}
		case IN_VALUE_RELEASE_SEEN:
			if err == io.EOF {
				s.tokens = append(s.tokens, ScannerToken{tType: ERROR, value: "", line: line, column: column, err: fmt.Errorf("eof after release char")})
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
