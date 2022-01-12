package scanner

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestScansEOF(t *testing.T) {

	s := NewScanner(bytes.NewReader([]byte{}))
	tokens, err := s.All()
	if err != nil {
		t.Fatal(err)
		return
	}

	assert.Len(t, tokens, 1)
	assert.Equal(t, ScannerToken{tType: EOF, value: "", pos: 0, err: nil}, tokens[0])
}

func TestScansSegmentTerminator(t *testing.T) {

	s := NewScanner(bytes.NewReader([]byte("'")))
	tokens, err := s.All()
	if err != nil {
		t.Fatal(err)
		return
	}

	assert.Len(t, tokens, 2)
	assert.Equal(t, ScannerToken{tType: SEGMENT_TERMINATOR, value: "'", pos: 0, err: nil}, tokens[0])
	assert.Equal(t, ScannerToken{tType: EOF, value: "", pos: 1, err: nil}, tokens[1])
}

func TestScansValue(t *testing.T) {

	s := NewScanner(bytes.NewReader([]byte("ABC")))
	tokens, err := s.All()
	if err != nil {
		t.Fatal(err)
		return
	}

	assert.Len(t, tokens, 2)
	assert.Equal(t, ScannerToken{tType: VALUE, value: "ABC", pos: 0, err: nil}, tokens[0])
	assert.Equal(t, ScannerToken{tType: EOF, value: "", pos: 3, err: nil}, tokens[1])
}

func TestScansValueWithReleaseChar(t *testing.T) {

	s := NewScanner(bytes.NewReader([]byte("ABC?'")))
	tokens, err := s.All()
	if err != nil {
		t.Fatal(err)
		return
	}

	assert.Len(t, tokens, 2)
	assert.Equal(t, ScannerToken{tType: VALUE, value: "ABC'", pos: 0, err: nil}, tokens[0])
	assert.Equal(t, ScannerToken{tType: EOF, value: "", pos: 5, err: nil}, tokens[1])
}

func TestScansUNASegment(t *testing.T) {

	s := NewScanner(bytes.NewReader([]byte("UNA:+.? !ABC!")))
	tokens, err := s.All()
	if err != nil {
		t.Fatal(err)
		return
	}

	assert.Len(t, tokens, 4)
	assert.Equal(t, ScannerToken{tType: UNA_SEGMENT, value: "UNA:+.? !", pos: 0, err: nil}, tokens[0])
	assert.Equal(t, ScannerToken{tType: VALUE, value: "ABC", pos: 9, err: nil}, tokens[1])
	assert.Equal(t, ScannerToken{tType: SEGMENT_TERMINATOR, value: "!", pos: 12, err: nil}, tokens[2])
	assert.Equal(t, ScannerToken{tType: EOF, value: "", pos: 13, err: nil}, tokens[3])
}

func TestScansUNASegmentSingleTimeOnly(t *testing.T) {

	s := NewScanner(bytes.NewReader([]byte("UNA:+.? !UNA:+.? !")))
	tokens, err := s.All()
	if err != nil {
		t.Fatal(err)
		return
	}

	assert.Len(t, tokens, 7)
	assert.Equal(t, ScannerToken{tType: UNA_SEGMENT, value: "UNA:+.? !", pos: 0, err: nil}, tokens[0])
	assert.Equal(t, ScannerToken{tType: VALUE, value: "UNA", pos: 9, err: nil}, tokens[1])
	assert.Equal(t, ScannerToken{tType: COMPONENT_DATA_ELEMENT_SEPERATOR, value: ":", pos: 12, err: nil}, tokens[2])
	assert.Equal(t, ScannerToken{tType: DATA_ELEMENT_SEPERATOR, value: "+", pos: 13, err: nil}, tokens[3])
	assert.Equal(t, ScannerToken{tType: VALUE, value: ". ", pos: 14, err: nil}, tokens[4])
	assert.Equal(t, ScannerToken{tType: SEGMENT_TERMINATOR, value: "!", pos: 17, err: nil}, tokens[5])
	assert.Equal(t, ScannerToken{tType: EOF, value: "", pos: 18, err: nil}, tokens[6])
}
