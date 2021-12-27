package stream

import (
	iparser "github.com/edilib/go-edi/internal/stream"
	"github.com/edilib/go-edi/stream/types"
	"io"
)

type EDITokenReader interface {
	Peek() (types.EDIToken, error)
	Read() (types.EDIToken, error)
	ReadAll() ([]types.EDIToken, error)
	Close() error
}

func NewEDITokenReader(reader io.Reader) EDITokenReader {
	return EDITokenReader(iparser.NewParser(reader))
}
