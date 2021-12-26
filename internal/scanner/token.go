package scanner

type ScannerTokenType int

const (
	UNA_SEGMENT ScannerTokenType = iota
	REPETITION_SEPERATOR
	COMPONENT_DATA_ELEMENT_SEPERATOR // default :
	DATA_ELEMENT_SEPERATOR           // default +
	SEGMENT_TERMINATOR               // default '
	SEGMENT_TAG
	VALUE
	EOF
	ERROR
)

type ScannerToken struct {
	tType ScannerTokenType
	value string
	pos   int
	err   error
}
