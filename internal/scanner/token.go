package scanner

import "fmt"

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
	tType  ScannerTokenType
	value  string
	line   int
	column int
	err    error
}

func (t ScannerTokenType) Name() string {
	return string(t)
}

func (t *ScannerToken) String() string {
	return fmt.Sprintf("type=%s,value=%s at <unknown>@%s", t.tType.Name(), t.Value(), t.Pos())
}

func (t *ScannerToken) Type() ScannerTokenType {
	return t.tType
}

func (t *ScannerToken) Pos() string {
	return fmt.Sprintf("%d:%d", t.line, t.column)
}

func (t *ScannerToken) Value() string {
	return t.value
}

func (t *ScannerToken) Error() error {
	return t.err
}
