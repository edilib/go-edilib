package edi2jsonln

import (
	"fmt"
	"github.com/edilib/go-edi/stream"
	"github.com/edilib/go-edi/stream/types"
	"os"
)

func Run(prog string, args []string) error {

	if len(args) != 1 {
		_, _ = fmt.Fprintf(os.Stderr, "usage: %s edifile.txt\n", prog)
		return nil
	}

	file, err := os.Open(args[0])
	if err != nil {
		return err
	}

	rd := stream.NewEDITokenReader(file)
	for {
		token, err := rd.Read()
		if err != nil {
			return err
		}

		if token.Type() == types.EOF {
			_, _ = fmt.Fprintf(os.Stderr, "%s: %s @ position %d\n", token.Type(), token.Value(), token.Pos())
			return nil
		} else if token.Type() == types.ERROR {
			return token.Error()
		} else {
			_, _ = fmt.Fprintf(os.Stderr, "%s: %s @ position %d\n", token.Type(), token.Value(), token.Pos())
		}
	}
}
