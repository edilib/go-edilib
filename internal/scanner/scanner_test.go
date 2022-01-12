package scanner

import (
	"bytes"
	"github.com/edilib/go-edilib/edifact/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestScansEOF(t *testing.T) {

	s := NewScanner(bytes.NewReader([]byte{}), types.UnEdifactFormat())
	tokens, err := s.All()
	if err != nil {
		t.Fatal(err)
		return
	}

	assert.Len(t, tokens, 1)
	assert.Equal(t, ScannerToken{tType: EOF, value: "", line: 0, column: 0, err: nil}, tokens[0])
}

func TestScansSegmentTerminator(t *testing.T) {

	s := NewScanner(bytes.NewReader([]byte("'")), types.UnEdifactFormat())
	tokens, err := s.All()
	if err != nil {
		t.Fatal(err)
		return
	}

	assert.Len(t, tokens, 2)
	assert.Equal(t, ScannerToken{tType: SEGMENT_TERMINATOR, value: "'", line: 0, column: 0, err: nil}, tokens[0])
	assert.Equal(t, ScannerToken{tType: EOF, value: "", line: 0, column: 1, err: nil}, tokens[1])
}

func TestScansValue(t *testing.T) {

	s := NewScanner(bytes.NewReader([]byte("ABC")), types.UnEdifactFormat())
	tokens, err := s.All()
	if err != nil {
		t.Fatal(err)
		return
	}

	assert.Len(t, tokens, 2)
	assert.Equal(t, ScannerToken{tType: VALUE, value: "ABC", line: 0, column: 0, err: nil}, tokens[0])
	assert.Equal(t, ScannerToken{tType: EOF, value: "", line: 0, column: 3, err: nil}, tokens[1])
}

func TestScansValueWithReleaseChar(t *testing.T) {

	s := NewScanner(bytes.NewReader([]byte("ABC?'")), types.UnEdifactFormat())
	tokens, err := s.All()
	if err != nil {
		t.Fatal(err)
		return
	}

	assert.Len(t, tokens, 2)
	assert.Equal(t, ScannerToken{tType: VALUE, value: "ABC'", line: 0, column: 0, err: nil}, tokens[0])
	assert.Equal(t, ScannerToken{tType: EOF, value: "", line: 0, column: 5, err: nil}, tokens[1])
}

func TestScansUNASegment(t *testing.T) {

	s := NewScanner(bytes.NewReader([]byte("UNA:+.? !ABC!")), types.UnEdifactFormat())
	tokens, err := s.All()
	if err != nil {
		t.Fatal(err)
		return
	}

	assert.Len(t, tokens, 4)
	assert.Equal(t, ScannerToken{tType: UNA_SEGMENT, value: "UNA:+.? !", line: 0, column: 0, err: nil}, tokens[0])
	assert.Equal(t, ScannerToken{tType: VALUE, value: "ABC", line: 0, column: 9, err: nil}, tokens[1])
	assert.Equal(t, ScannerToken{tType: SEGMENT_TERMINATOR, value: "!", line: 0, column: 12, err: nil}, tokens[2])
	assert.Equal(t, ScannerToken{tType: EOF, value: "", line: 0, column: 13, err: nil}, tokens[3])
}

func TestScansUNASegmentSingleTimeOnly(t *testing.T) {

	s := NewScanner(bytes.NewReader([]byte("UNA:+.? !UNA:+.? !")), types.UnEdifactFormat())
	tokens, err := s.All()
	if err != nil {
		t.Fatal(err)
		return
	}

	assert.Len(t, tokens, 7)
	assert.Equal(t, ScannerToken{tType: UNA_SEGMENT, value: "UNA:+.? !", line: 0, column: 0, err: nil}, tokens[0])
	assert.Equal(t, ScannerToken{tType: VALUE, value: "UNA", line: 0, column: 9, err: nil}, tokens[1])
	assert.Equal(t, ScannerToken{tType: COMPONENT_DATA_ELEMENT_SEPERATOR, value: ":", line: 0, column: 12, err: nil}, tokens[2])
	assert.Equal(t, ScannerToken{tType: DATA_ELEMENT_SEPERATOR, value: "+", line: 0, column: 13, err: nil}, tokens[3])
	assert.Equal(t, ScannerToken{tType: VALUE, value: ". ", line: 0, column: 14, err: nil}, tokens[4])
	assert.Equal(t, ScannerToken{tType: SEGMENT_TERMINATOR, value: "!", line: 0, column: 17, err: nil}, tokens[5])
	assert.Equal(t, ScannerToken{tType: EOF, value: "", line: 0, column: 18, err: nil}, tokens[6])
}
