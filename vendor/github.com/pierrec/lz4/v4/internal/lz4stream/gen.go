//+build ignore

package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/pierrec/lz4/v4/internal/lz4block"
	"github.com/pierrec/packer"
	"golang.org/x/tools/imports"
)

type DescriptorFlags struct {
	// FLG
	_                 [2]int
	ContentChecksum   [1]bool
	Size              [1]bool
	BlockChecksum     [1]bool
	BlockIndependence [1]bool
	Version           [2]uint16
	// BD
	_              [4]int
	BlockSizeIndex [3]lz4block.BlockSizeIndex
	_              [1]int
}

type DataBlockSize struct {
	size         [31]int
	Uncompressed bool
}

func main() {
	err := do()
	if err != nil {
		log.Fatal(err)
	}
}

func do() error {
	out, err := os.Create("frame_gen.go")
	if err != nil {
		return err
	}
	defer out.Close()

	pkg := "lz4stream"
	buf := new(bytes.Buffer)
	for i, t := range []interface{}{
		DescriptorFlags{}, DataBlockSize{},
	} {
		if i > 0 {
			pkg = ""
		}
		err := packer.GenPackedStruct(buf, &packer.Config{PkgName: pkg}, t)
		if err != nil {
			return fmt.Errorf("%T: %v", t, err)
		}
	}
	// Resolve imports.
	code, err := imports.Process("", buf.Bytes(), nil)
	if err != nil {
		// Output without imports.
		_, _ = io.Copy(out, buf)
		return err
	}
	_, err = io.Copy(out, bytes.NewReader(code))
	return err
}
