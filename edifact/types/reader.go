package types

// SegmentReader is the interface to reading Edifact types.Segment slices.
type SegmentReader interface {
	// HasNext checks if there are more Segment values available in the data stream.
	HasNext() (bool, error)
	// Next returns the next Segment or an error.
	Next() (Segment, error)
	// ReadAll reads all Segment values left in the data stream.
	ReadAll() ([]Segment, error)
	// Close closes the reader an the underlying reader if it is closeable too.
	Close() error
}
