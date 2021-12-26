package types

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

type EDIToken interface {
	Type() EDITokenType
	Pos() int
	Error() error
	Value() string
}
