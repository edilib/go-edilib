# go-edi - edifact reader written in golang

## Features

* reads generic edifact format
* UN/EDIFACT support
* X12 support

## Usage (generic edifact example)

edilib supports go modules.

```go
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
```

## License

Copyright (c) 2021 by [Cornelius Buschka](https://github.com/edilib).

[Apache License, Version 2.0](./license.txt)
