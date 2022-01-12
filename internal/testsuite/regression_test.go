package testsuite

import (
	"github.com/edilib/go-edilib/edifact"
	"github.com/edilib/go-edilib/edifact/types"
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
	p := edifact.NewSegmentReader(file, types.UnEdifactFormat())
	segments, err := p.ReadAll()
	if err != nil {
		t.Fatal(err)
		return
	}

	assert.Len(t, segments, 31)
}

func TestParsesHIPAA(t *testing.T) {
	file, err := os.Open("hipaa1.txt")
	if err != nil {
		t.Fatal(err)
		return
	}
	p := edifact.NewSegmentReader(file, types.X12Format())
	segments, err := p.ReadAll()
	if err != nil {
		t.Fatal(err)
		return
	}

	assert.Len(t, segments, 75)
}
