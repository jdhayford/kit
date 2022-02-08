package main

import (
	"kit/kit"
	"os"
	"strings"
)

func main() {
	var args []string
	if len(os.Args) > 2 {
		args = os.Args[2:]
	} else {
		args = os.Args[1:]
	}

	firstArg := ""
	if len(args) > 0 {
		firstArg = args[0]
	}

	switch {
	case strings.HasPrefix(firstArg, "_"):
		return
	case strings.HasPrefix(firstArg, "@"):
		kit.RunUserStrategy(args)
	default:
		kit.RunContextStrategy(args)
	}
}
