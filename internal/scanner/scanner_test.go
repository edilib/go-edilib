package scanner

import (
	"bytes"
	"github.com/edilib/go-edilib/edifact/types"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"math/big"
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
	assert.Equal(t, ScannerToken{tType: EOF, stringValue: "", integerValue: nil, decimalValue: nil, line: 0, column: 0, err: nil}, tokens[0])
}

func TestScansSegmentTerminator(t *testing.T) {

	s := NewScanner(bytes.NewReader([]byte("'")), types.UnEdifactFormat())
	tokens, err := s.All()
	if err != nil {
		t.Fatal(err)
		return
	}

	assert.Len(t, tokens, 2)
	assert.Equal(t, ScannerToken{tType: SEGMENT_TERMINATOR, stringValue: "'", integerValue: nil, decimalValue: nil, line: 0, column: 0, err: nil}, tokens[0])
	assert.Equal(t, ScannerToken{tType: EOF, stringValue: "", integerValue: nil, decimalValue: nil, line: 0, column: 1, err: nil}, tokens[1])
}

func TestScansStringValue(t *testing.T) {

	s := NewScanner(bytes.NewReader([]byte("ABC")), types.UnEdifactFormat())
	tokens, err := s.All()
	if err != nil {
		t.Fatal(err)
		return
	}

	assert.Len(t, tokens, 2)
	assert.Equal(t, ScannerToken{tType: VALUE, stringValue: "ABC", integerValue: nil, decimalValue: nil, line: 0, column: 0, err: nil}, tokens[0])
	assert.Equal(t, ScannerToken{tType: EOF, stringValue: "", integerValue: nil, decimalValue: nil, line: 0, column: 3, err: nil}, tokens[1])
}

func TestScansIntegerValue(t *testing.T) {

	s := NewScanner(bytes.NewReader([]byte("100")), types.UnEdifactFormat())
	tokens, err := s.All()
	if err != nil {
		t.Fatal(err)
		return
	}

	assert.Len(t, tokens, 2)
	integerValue := big.NewInt(100)
	decimalValue, _ := decimal.NewFromString("100")
	assert.Equal(t, ScannerToken{tType: VALUE, stringValue: "100", integerValue: integerValue, decimalValue: &decimalValue, line: 0, column: 0, err: nil}, tokens[0])
	assert.Equal(t, ScannerToken{tType: EOF, stringValue: "", integerValue: nil, decimalValue: nil, line: 0, column: 3, err: nil}, tokens[1])
}

func TestScansDecimalValue(t *testing.T) {

	format := types.UnEdifactFormat()
	format.DecimalMark = ','
	s := NewScanner(bytes.NewReader([]byte("123,456")), format)
	tokens, err := s.All()
	if err != nil {
		t.Fatal(err)
		return
	}

	assert.Len(t, tokens, 2)
	decimalValue := decimal.RequireFromString("123.456")
	assert.Equal(t, ScannerToken{tType: VALUE, stringValue: "123,456", integerValue: nil, decimalValue: &decimalValue, line: 0, column: 0, err: nil}, tokens[0])
	assert.Equal(t, ScannerToken{tType: EOF, stringValue: "", integerValue: nil, decimalValue: nil, line: 0, column: 7, err: nil}, tokens[1])
}

func TestScansValueWithReleaseChar(t *testing.T) {

	s := NewScanner(bytes.NewReader([]byte("ABC?'")), types.UnEdifactFormat())
	tokens, err := s.All()
	if err != nil {
		t.Fatal(err)
		return
	}

	assert.Len(t, tokens, 2)
	assert.Equal(t, ScannerToken{tType: VALUE, stringValue: "ABC'", integerValue: nil, decimalValue: nil, line: 0, column: 0, err: nil}, tokens[0])
	assert.Equal(t, ScannerToken{tType: EOF, stringValue: "", integerValue: nil, decimalValue: nil, line: 0, column: 5, err: nil}, tokens[1])
}

func TestScansUNASegment(t *testing.T) {

	s := NewScanner(bytes.NewReader([]byte("UNA:+.? !ABC!")), types.UnEdifactFormat())
	tokens, err := s.All()
	if err != nil {
		t.Fatal(err)
		return
	}

	assert.Len(t, tokens, 4)
	assert.Equal(t, ScannerToken{tType: UNA_SEGMENT, stringValue: "UNA:+.? !", integerValue: nil, decimalValue: nil, line: 0, column: 0, err: nil}, tokens[0])
	assert.Equal(t, ScannerToken{tType: VALUE, stringValue: "ABC", integerValue: nil, decimalValue: nil, line: 0, column: 9, err: nil}, tokens[1])
	assert.Equal(t, ScannerToken{tType: SEGMENT_TERMINATOR, stringValue: "!", integerValue: nil, decimalValue: nil, line: 0, column: 12, err: nil}, tokens[2])
	assert.Equal(t, ScannerToken{tType: EOF, stringValue: "", integerValue: nil, decimalValue: nil, line: 0, column: 13, err: nil}, tokens[3])
}

func TestScansUNASegmentSingleTimeOnly(t *testing.T) {

	s := NewScanner(bytes.NewReader([]byte("UNA:+.? !UNA:+.? !")), types.UnEdifactFormat())
	tokens, err := s.All()
	if err != nil {
		t.Fatal(err)
		return
	}

	assert.Len(t, tokens, 7)
	assert.Equal(t, ScannerToken{tType: UNA_SEGMENT, stringValue: "UNA:+.? !", integerValue: nil, decimalValue: nil, line: 0, column: 0, err: nil}, tokens[0])
	assert.Equal(t, ScannerToken{tType: VALUE, stringValue: "UNA", integerValue: nil, decimalValue: nil, line: 0, column: 9, err: nil}, tokens[1])
	assert.Equal(t, ScannerToken{tType: COMPONENT_DATA_ELEMENT_SEPERATOR, stringValue: ":", integerValue: nil, decimalValue: nil, line: 0, column: 12, err: nil}, tokens[2])
	assert.Equal(t, ScannerToken{tType: DATA_ELEMENT_SEPERATOR, stringValue: "+", integerValue: nil, decimalValue: nil, line: 0, column: 13, err: nil}, tokens[3])
	assert.Equal(t, ScannerToken{tType: VALUE, stringValue: ". ", integerValue: nil, decimalValue: nil, line: 0, column: 14, err: nil}, tokens[4])
	assert.Equal(t, ScannerToken{tType: SEGMENT_TERMINATOR, stringValue: "!", integerValue: nil, decimalValue: nil, line: 0, column: 17, err: nil}, tokens[5])
	assert.Equal(t, ScannerToken{tType: EOF, stringValue: "", integerValue: nil, decimalValue: nil, line: 0, column: 18, err: nil}, tokens[6])
}
