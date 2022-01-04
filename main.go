package main

import (
	"kit/kit"
	"os"
)

func main() {
	var args []string
	if len(os.Args) > 2 {
		args = os.Args[2:]
	} else {
		args = os.Args[1:]
	}

	kitFilePath, err := kit.FindKitFile()
	if err != nil {
		panic(err)
	}

	targetKit := kit.ParseKitFile(kitFilePath)
	targetKit.Run(args)
}
