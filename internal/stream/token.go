package stream

import (
	"github.com/edilib/go-edilib/stream/types"
)

type ScannerToken struct {
	tType types.EDITokenType
	value string
	pos   int
	err   error
}

func (t *ScannerToken) Type() types.EDITokenType {
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
