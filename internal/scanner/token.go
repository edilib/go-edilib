package scanner

import (
	"fmt"
	"github.com/shopspring/decimal"
	"math/big"
)

type ScannerTokenType string

const (
	UNA_SEGMENT                      ScannerTokenType = "UNA_SEGMENT"
	REPETITION_SEPERATOR             ScannerTokenType = "REPETITION_SEPERATOR"
	COMPONENT_DATA_ELEMENT_SEPERATOR ScannerTokenType = "COMPONENT_DATA_ELEMENT_SEPERATOR" // default :
	DATA_ELEMENT_SEPERATOR           ScannerTokenType = "DATA_ELEMENT_SEPERATOR"           // default +
	SEGMENT_TERMINATOR               ScannerTokenType = "SEGMENT_TERMINATOR"               // default '
	VALUE                            ScannerTokenType = "VALUE"
	EOF                              ScannerTokenType = "EOF"
	ERROR                            ScannerTokenType = "ERROR"
)

type ScannerToken struct {
	tType        ScannerTokenType
	stringValue  string
	integerValue *big.Int
	decimalValue *decimal.Decimal
	line         int
	column       int
	err          error
}

func (t ScannerTokenType) Name() string {
	return string(t)
}

func (t *ScannerToken) String() string {
	return fmt.Sprintf("type=%s,value=%s at <unknown>@%s", t.tType.Name(), t.StringValue(), t.Pos())
}

func (t *ScannerToken) Type() ScannerTokenType {
	return t.tType
}

func (t *ScannerToken) Pos() string {
	return fmt.Sprintf("%d:%d", t.line, t.column)
}

func (t *ScannerToken) StringValue() string {
	return t.stringValue
}

func (t *ScannerToken) IntegerValue() *big.Int {
	return t.integerValue
}

func (t *ScannerToken) DecimalValue() *decimal.Decimal {
	return t.decimalValue
}
func (t *ScannerToken) Error() error {
	return t.err
}
