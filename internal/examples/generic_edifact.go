package examples

import (
	"fmt"
	"github.com/edilib/go-edilib/edifact"
	"github.com/edilib/go-edilib/edifact/types"
	"os"
)

func main() {
	file, _ := os.Open("edifact-file.txt")

	p := edifact.NewSegmentReader(file, types.UnEdifactFormat())
	segments, _ := p.ReadAll()

	fmt.Printf("%v", segments)
}
