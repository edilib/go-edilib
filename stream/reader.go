package stream

import (
	iparser "github.com/cbuschka/go-edi/internal/stream"
	"github.com/cbuschka/go-edi/stream/types"
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
