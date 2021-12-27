package edi2json

import (
	"fmt"
	"github.com/cbuschka/go-jsonstream"
	"github.com/edilib/go-edi/dom"
	"github.com/edilib/go-edi/dom/types"
	"os"
)

func Run(prog string, args []string) error {

	if len(args) != 1 {
		_, _ = fmt.Fprintf(os.Stderr, "usage: %s edifile.txt\n", prog)
		return nil
	}

	file, err := os.Open(args[0])
	if err != nil {
		return err
	}

	parser := dom.NewParser(file)
	segments, err := parser.ReadAll()
	if err != nil {
		return err
	}

	wr := jsonstream.NewWriter(os.Stderr)
	wr.SetIndent("  ")
	_ = wr.StartObject()
	_ = wr.Key("segments")
	err = writeSegments(wr, segments)
	if err != nil {
		return err
	}
	_ = wr.EndObject()

	return nil
}

func writeSegments(wr jsonstream.Writer, segments []types.EDISegment) error {
	_ = wr.StartArray()
	for _, segment := range segments {
		err := writeSegment(wr, segment)
		if err != nil {
			return err
		}
	}
	_ = wr.EndArray()
	return nil
}

func writeSegment(wr jsonstream.Writer, segment types.EDISegment) error {
	_ = wr.StartObject()
	_ = wr.Key("tag")
	_ = wr.String(segment.Tag)
	err := wr.Key("values")
	if err != nil {
		return err
	}
	err = writeElements(wr, segment.Elements)
	if err != nil {
		return err
	}
	_ = wr.EndObject()
	return nil
}

func writeElements(wr jsonstream.Writer, elements []interface{}) error {
	_ = wr.StartArray()
	for _, element := range elements {
		composite, isComposite := element.(*types.EDIComposite)
		stringValue, isString := element.(string)
		boolValue, isBool := element.(bool)
		numberValue, isNumber := element.(int)
		if isComposite {
			err := writeElements(wr, composite.Elements)
			if err != nil {
				return err
			}
		} else if isNumber {
			_ = wr.Number(numberValue)
		} else if element == nil {
			_ = wr.Null()
		} else if isBool {
			_ = wr.Boolean(boolValue)
		} else if isString {
			_ = wr.String(stringValue)
		} else {
			return fmt.Errorf("unsupported value type")
		}

	}
	_ = wr.EndArray()
	return nil
}
