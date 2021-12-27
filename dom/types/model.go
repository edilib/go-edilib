package types

type EDIComposite struct {
	Elements []interface{}
}

type EDISegment struct {
	Tag      string
	Elements []interface{}
}
