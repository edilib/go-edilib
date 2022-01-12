package types

type Value interface {
}

type SimpleValue struct {
	StringValue string
}

type CompositeValue struct {
	Values []SimpleValue
}

type RepetitionValue struct {
	Values []SimpleValue
}

type Tag struct {
	Values CompositeValue
}

type Segment struct {
	Tag          Tag
	DataElements []Value
}

func (t Tag) Name() string {
	return t.Values.Values[0].StringValue
}
