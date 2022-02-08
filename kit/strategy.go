package kit

import (
	"fmt"
	"os"
	"strings"
)

func RunContextStrategy(args []string) {
	kitFilePath, err := FindKitFile()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	targetKit, err := ParseKitFile(kitFilePath)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	targetKit.Run(args)
}

func RunUserStrategy(args []string) {
	if len(args[0]) < 1 {
		fmt.Println("Invalid user kit name")
		return
	}

	var test string
	fmt.Println(test)

	kitName := strings.TrimPrefix(args[0], "@")
	userKit, err := FindUserKit(kitName)
	if err != nil {
		fmt.Println(NoMatchingKitError{}.error)
		return
	}

	var kitArgs []string
	if len(args) > 1 {
		kitArgs = args[1:]
	}

	userKit.Run(kitArgs)
	// targetKit.Run(args)
}
