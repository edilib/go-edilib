package internal

import (
	"github.com/edilib/go-edilib/edifact/types"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestReadsSingleEmptySegment(t *testing.T) {
	rd := strings.NewReader("UNB+'")
	p := NewSegmentReader(rd)
	all, err := p.ReadAll()
	if err != nil {
		t.Fatal(err)
		return
	}

	assert.Equal(t, []types.Segment{
		{Tag: types.Tag{
			Values: types.CompositeValue{Values: []types.SimpleValue{{StringValue: "UNB"}}},
		}, DataElements: []types.Value{
			types.SimpleValue{StringValue: ""},
		}}}, all)
}

func TestReadsSingleSegmentWithSingleEmptySimpleDataElement(t *testing.T) {
	rd := strings.NewReader("UNB+'")
	p := NewSegmentReader(rd)
	all, err := p.ReadAll()
	if err != nil {
		t.Fatal(err)
		return
	}

	assert.Equal(t, []types.Segment{
		{Tag: types.Tag{
			Values: types.CompositeValue{Values: []types.SimpleValue{{StringValue: "UNB"}}},
		}, DataElements: []types.Value{
			types.SimpleValue{StringValue: ""},
		}}}, all)
}

func TestReadsSingleSegmentWithTwoEmptySimpleDataElements(t *testing.T) {
	rd := strings.NewReader("UNB++'")
	p := NewSegmentReader(rd)
	all, err := p.ReadAll()
	if err != nil {
		t.Fatal(err)
		return
	}

	assert.Equal(t, []types.Segment{
		{Tag: types.Tag{
			Values: types.CompositeValue{Values: []types.SimpleValue{{StringValue: "UNB"}}},
		}, DataElements: []types.Value{
			types.SimpleValue{StringValue: ""},
			types.SimpleValue{StringValue: ""},
		}}}, all)
}
func TestReadsTwoSegmentWithTwoRepetitionDataElements(t *testing.T) {
	rd := strings.NewReader("UNB+X*X2+Y*Y2'")
	p := NewSegmentReader(rd)
	all, err := p.ReadAll()
	if err != nil {
		t.Fatal(err)
		return
	}

	assert.Equal(t, []types.Segment{
		{Tag: types.Tag{
			Values: types.CompositeValue{Values: []types.SimpleValue{{StringValue: "UNB"}}},
		}, DataElements: []types.Value{
			types.RepetitionValue{Values: []types.SimpleValue{{StringValue: "X"}, {StringValue: "X2"}}},
			types.RepetitionValue{Values: []types.SimpleValue{{StringValue: "Y"}, {StringValue: "Y2"}}},
		}},
	}, all)
}
