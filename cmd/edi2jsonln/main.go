package main

import (
	"fmt"
	"github.com/edilib/go-edilib/internal/edi2jsonln"
	"os"
	"path/filepath"
)

func main() {
	fqProg := os.Args[0]
	prog := filepath.Base(fqProg)
	err := edi2jsonln.Run(prog, os.Args[1:])
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed: %s}n", err.Error())
		os.Exit(1)
	}
}
