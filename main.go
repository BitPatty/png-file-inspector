package main

import (
	"fmt"
	"io/ioutil"
	"os"

	png "github.com/bitpatty/png-file-inspector/inspector"
)

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Printf("missing file path")
		return
	}

	bytes, err := ioutil.ReadFile(args[0])

	if err != nil {
		panic(err)
	}

	vr, err := png.Inspect(bytes)

	if err != nil {
		defer fmt.Println(err)
	}

	fmt.Printf("Valid Signature: %+v\n", vr.HasValidSignature)
	fmt.Printf("Leading IHDR: %+v\n", vr.HasLeadingIHDR)
	fmt.Printf("Data after IEND: %+v\n", vr.HasDataAfterIEND)

	for i := 0; i < len(vr.Chunks); i++ {
		fmt.Printf("%v: %+v\n", vr.Chunks[i].Header.Value.ToString(), vr.Chunks[i])
	}
}
