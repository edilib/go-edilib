package types

type SegmentReader interface {
	HasNext() (bool, error)
	Next() (Segment, error)
	ReadAll() ([]Segment, error)
}
