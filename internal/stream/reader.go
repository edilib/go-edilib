package stream

import (
	"bufio"
	"io"
	"unicode/utf8"
)

type ScannerReader struct {
	ioRd       *bufio.Reader
	rdForClose io.Reader
	ringBuf    *RuneRingBuffer
	pos        int
}

func NewScannerReader(reader io.Reader) *ScannerReader {
	return &ScannerReader{rdForClose: reader, ioRd: bufio.NewReader(reader), ringBuf: NewRuneRingBuffer(128), pos: 0}
}

func (s *ScannerReader) Position() int {
	return s.pos
}

func (s *ScannerReader) PeekRune(i int) (rune, error) {
	err := s.tryFill(i + 1)
	if err != nil {
		return 0, err
	}

	if s.ringBuf.Size() < i+1 {
		return 0, io.EOF
	}

	r, err := s.ringBuf.Peek(i)
	return r, err
}

func (s *ScannerReader) tryFill(n int) error {
	if s.ringBuf.Size() >= n {
		return nil
	}

	for i := 0; i < n; i++ {
		r, _, err := s.ioRd.ReadRune()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		err = s.ringBuf.Add(r)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *ScannerReader) ReadRune() (rune, error) {
	err := s.tryFill(1)
	if err != nil {
		return 0, err
	}

	if s.ringBuf.Size() < 1 {
		return 0, io.EOF
	}

	r, err := s.ringBuf.Remove()
	s.pos = s.pos + utf8.RuneLen(r)

	return r, err
}

func (s *ScannerReader) Close() error {
	closer, isCloser := s.rdForClose.(io.ReadCloser)
	if isCloser {
		return closer.Close()
	}

	return nil
}

func (s *ScannerReader) PeekRunes(n int) ([]rune, error) {

	buf := []rune{}
	for i := 0; i < n; i++ {
		r, err := s.PeekRune(i)
		if err != nil {
			return nil, err
		}

		buf = append(buf, r)
	}

	return buf, nil

}

func (s *ScannerReader) ReadNRunes(n int) ([]rune, error) {

	buf := []rune{}
	for i := 0; i < n; i++ {
		r, err := s.ReadRune()
		if err != nil {
			return nil, err
		}

		buf = append(buf, r)
	}

	return buf, nil
}
