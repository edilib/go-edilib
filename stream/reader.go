package stream

import (
	iparser "github.com/edilib/go-edilib/internal/stream"
	"github.com/edilib/go-edilib/stream/types"
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
