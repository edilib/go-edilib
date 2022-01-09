package main

import (
	"fmt"
	"github.com/edilib/go-edilib/internal/edi2json"
	"os"
	"path/filepath"
)

func main() {
	fqProg := os.Args[0]
	prog := filepath.Base(fqProg)
	err := edi2json.Run(prog, os.Args[1:])
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed: %s}n", err.Error())
		os.Exit(1)
	}
}
