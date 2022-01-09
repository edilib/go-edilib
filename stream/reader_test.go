package stream

import (
	"bytes"
	"github.com/edilib/go-edilib/stream/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestScansEOF(t *testing.T) {

	s := NewEDITokenReader(bytes.NewReader([]byte{}))
	tokens, err := s.ReadAll()
	if err != nil {
		t.Fatal(err)
		return
	}

	assert.Len(t, tokens, 1)
	assert.Equal(t, tokens[0].Type(), types.ERROR)
}

func TestScansSegmentWithValues(t *testing.T) {

	s := NewEDITokenReader(bytes.NewReader([]byte("UNA:+.? 'ABC+VALUE1:VALUE2?''")))
	tokens, err := s.ReadAll()
	if err != nil {
		t.Fatal(err)
		return
	}

	assert.Len(t, tokens, 8)
	assert.Equal(t, types.UNA_SEGMENT, tokens[0].Type())
	assert.Equal(t, types.SEGMENT_TAG, tokens[1].Type())
	assert.Equal(t, "ABC", tokens[1].Value())
	assert.Equal(t, types.DATA_ELEMENT_SEPERATOR, tokens[2].Type())
	assert.Equal(t, types.VALUE, tokens[3].Type())
	assert.Equal(t, "VALUE1", tokens[3].Value())
	assert.Equal(t, types.COMPONENT_DATA_ELEMENT_SEPERATOR, tokens[4].Type())
	assert.Equal(t, types.VALUE, tokens[5].Type())
	assert.Equal(t, "VALUE2'", tokens[5].Value())
	assert.Equal(t, types.SEGMENT_TERMINATOR, tokens[6].Type())
	assert.Equal(t, types.EOF, tokens[7].Type())
}
