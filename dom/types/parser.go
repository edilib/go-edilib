package types

type Parser interface {
	Read() (*EDISegment, error)
}
