package scanner

import "fmt"

type EDITokenType string

const (
	UNA_SEGMENT                      EDITokenType = "UNA_SEGMENT"
	SEGMENT_TAG                      EDITokenType = "SEGMENT_TAG"
	REPETITION_SEPERATOR             EDITokenType = "REPETITION_SEPERATOR"
	COMPONENT_DATA_ELEMENT_SEPERATOR EDITokenType = "COMPONENT_DATA_ELEMENT_SEPERATOR" // default :
	DATA_ELEMENT_SEPERATOR           EDITokenType = "DATA_ELEMENT_SEPERATOR"           // default +
	SEGMENT_TERMINATOR               EDITokenType = "SEGMENT_TERMINATOR"               // default '
	VALUE                            EDITokenType = "VALUE"
	EOF                              EDITokenType = "EOF"
	ERROR                            EDITokenType = "ERROR"
)

type ScannerToken struct {
	tType EDITokenType
	value string
	pos   int
	err   error
}

func (t EDITokenType) Name() string {
	return string(t)
}

func (t *ScannerToken) String() string {
	return fmt.Sprintf("type=%s,value=%s at <unknown>:0:%d", t.tType.Name(), t.Value(), t.pos)
}

func (t *ScannerToken) Type() EDITokenType {
	return t.tType
}

func (t *ScannerToken) Pos() int {
	return t.pos
}

func (t *ScannerToken) Value() string {
	return t.value
}

func (t *ScannerToken) Error() error {
	return t.err
}
