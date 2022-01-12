package edi

import (
	"github.com/edilib/go-edilib/edifact"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestParsesINVRPT1(t *testing.T) {
	file, err := os.Open("invrpt1.txt")
	if err != nil {
		t.Fatal(err)
		return
	}
	p := edifact.NewSegmentReader(file)
	segments, err := p.ReadAll()
	if err != nil {
		t.Fatal(err)
		return
	}

	assert.Len(t, segments, 31)
}
