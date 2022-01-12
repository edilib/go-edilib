package edifact

import (
	"github.com/edilib/go-edilib/edifact/types"
	"github.com/edilib/go-edilib/internal"
	"io"
)

// NewSegmentReader returns a new reader instance that can read EDI segments
// from io reader.
func NewSegmentReader(reader io.Reader, format types.Format) types.SegmentReader {
	return internal.NewSegmentReader(reader, format)
}
