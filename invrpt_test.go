package edi

import (
	"github.com/cbuschka/go-edi/stream"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestScansINVRPT1(t *testing.T) {
	file, err := os.Open("invrpt1.txt")
	if err != nil {
		t.Fatal(err)
		return
	}
	p := stream.NewEDITokenReader(file)
	all, err := p.ReadAll()
	if err != nil {
		t.Fatal(err)
		return
	}

	assert.Len(t, all, 249)
}
