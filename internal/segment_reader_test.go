package internal

import (
	"github.com/edilib/go-edilib/edifact/types"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"math/big"
	"strings"
	"testing"
)

func TestReadsSingleEmptySegment(t *testing.T) {
	rd := strings.NewReader("UNB+'")
	p := NewSegmentReader(rd, types.UnEdifactFormat())
	all, err := p.ReadAll()
	if err != nil {
		t.Fatal(err)
		return
	}

	assert.Equal(t, []types.Segment{
		{Tag: types.Tag{
			Values: types.CompositeValue{Values: []types.SimpleValue{{StringValue: "UNB"}}},
		}, DataElements: []types.Value{
			types.SimpleValue{StringValue: "", IntegerValue: nil, DecimalValue: nil},
		}}}, all)
}

func TestReadsSingleSegmentWithSingleEmptySimpleDataElement(t *testing.T) {
	rd := strings.NewReader("UNB+'")
	p := NewSegmentReader(rd, types.UnEdifactFormat())
	all, err := p.ReadAll()
	if err != nil {
		t.Fatal(err)
		return
	}

	assert.Equal(t, []types.Segment{
		{Tag: types.Tag{
			Values: types.CompositeValue{Values: []types.SimpleValue{{StringValue: "UNB"}}},
		}, DataElements: []types.Value{
			types.SimpleValue{StringValue: "", IntegerValue: nil, DecimalValue: nil},
		}}}, all)
}

func TestReadsSingleSegmentWithTwoEmptySimpleDataElements(t *testing.T) {
	rd := strings.NewReader("UNB++'")
	p := NewSegmentReader(rd, types.UnEdifactFormat())
	all, err := p.ReadAll()
	if err != nil {
		t.Fatal(err)
		return
	}

	assert.Equal(t, []types.Segment{
		{Tag: types.Tag{
			Values: types.CompositeValue{Values: []types.SimpleValue{{StringValue: "UNB"}}},
		}, DataElements: []types.Value{
			types.SimpleValue{StringValue: "", IntegerValue: nil, DecimalValue: nil},
			types.SimpleValue{StringValue: "", IntegerValue: nil, DecimalValue: nil},
		}}}, all)
}

func TestReadsTwoSegmentWithTwoRepetitionDataElements(t *testing.T) {
	rd := strings.NewReader("UNB+X*X2+Y*Y2'")
	p := NewSegmentReader(rd, types.UnEdifactFormat())
	all, err := p.ReadAll()
	if err != nil {
		t.Fatal(err)
		return
	}

	assert.Equal(t, []types.Segment{
		{Tag: types.Tag{
			Values: types.CompositeValue{Values: []types.SimpleValue{{StringValue: "UNB"}}},
		}, DataElements: []types.Value{
			types.RepetitionValue{Values: []types.SimpleValue{{StringValue: "X", IntegerValue: nil, DecimalValue: nil}, {StringValue: "X2", IntegerValue: nil, DecimalValue: nil}}},
			types.RepetitionValue{Values: []types.SimpleValue{{StringValue: "Y", IntegerValue: nil, DecimalValue: nil}, {StringValue: "Y2", IntegerValue: nil, DecimalValue: nil}}},
		}},
	}, all)
}

func TestReadsSingleSegmentWithDecimalDataElement(t *testing.T) {
	rd := strings.NewReader("UNB+123.456'")
	p := NewSegmentReader(rd, types.UnEdifactFormat())
	all, err := p.ReadAll()
	if err != nil {
		t.Fatal(err)
		return
	}

	decimalValue := decimal.RequireFromString("123.456")
	assert.Equal(t, []types.Segment{
		{Tag: types.Tag{
			Values: types.CompositeValue{Values: []types.SimpleValue{{StringValue: "UNB"}}},
		}, DataElements: []types.Value{
			types.SimpleValue{StringValue: "123.456", IntegerValue: nil, DecimalValue: &decimalValue},
		}},
	}, all)
}

func TestReadsSingleSegmentWithIntegerDataElement(t *testing.T) {
	rd := strings.NewReader("UNB+123'")
	p := NewSegmentReader(rd, types.UnEdifactFormat())
	all, err := p.ReadAll()
	if err != nil {
		t.Fatal(err)
		return
	}

	decimalValue := decimal.RequireFromString("123")
	integerValue := big.NewInt(123)
	assert.Equal(t, []types.Segment{
		{Tag: types.Tag{
			Values: types.CompositeValue{Values: []types.SimpleValue{{StringValue: "UNB"}}},
		}, DataElements: []types.Value{
			types.SimpleValue{StringValue: "123", IntegerValue: integerValue, DecimalValue: &decimalValue},
		}},
	}, all)
}

func TestUnaNotAllowedWithX12(t *testing.T) {
	rd := strings.NewReader("UNA...")
	p := NewSegmentReader(rd, types.X12Format())
	_, err := p.ReadAll()

	assert.EqualError(t, err, "una segment not allowed at 0:0")
}
