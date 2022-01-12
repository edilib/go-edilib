package types

import (
	"github.com/shopspring/decimal"
	"math/big"
)

// Value is an value in the edifact context. It is used for Tag names, Tag components,
// and DataElement values.
type Value interface {
}

// SimpleValue is a simple value, e.g. a string.
type SimpleValue struct {
	StringValue  string
	DecimalValue *decimal.Decimal
	IntegerValue *big.Int
}

// CompositeValue represents an edifact value with component values. Also
// used for segment Tag objects.
type CompositeValue struct {
	Values []SimpleValue
}

// RepetitionValue is a repetition of values within a single data element.
type RepetitionValue struct {
	Values []SimpleValue
}

// Tag is the beginning of an Edifact Segment.
type Tag struct {
	Values CompositeValue
}

// Segment is a Edifact Segment that consists of a Tag and DataElements.
type Segment struct {
	Tag          Tag
	DataElements []Value
}

// Name returns the name of a segment Tag.
func (t Tag) Name() string {
	return t.Values.Values[0].StringValue
}
