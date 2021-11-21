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

	vr, err := png.Inspect(bytes, png.InspectOptions{
		AllowUnknownAncillaryChunks: false,
		AllowUnknownCriticalChunks:  false,
	})

	if err != nil {
		defer fmt.Println(err)
	}

	fmt.Printf("%+v\n", vr)
}
