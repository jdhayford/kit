package kit

import (
	"fmt"
	"os"
	"strings"
)

func RunContextStrategy(args []string) {
	// Create KitSet from global kits and context kit (if context kit is found)
	kits := GetGlobalUserKits()

	contextKit, err := FindContextKit()
	if err != nil {
		if _, ok := err.(*NoContextKitFoundError); !ok {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	}

	if !contextKit.IsEmpty() {
		kits = append(kits, contextKit)
	}

	kitSet := MakeKitSet(kits...)

	kitSet.Run(args)
}

func RunUserStrategy(args []string) {
	if len(args[0]) < 1 {
		fmt.Println("Invalid user kit name")
		return
	}

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

	kitSet := MakeKitSet(userKit)
	kitSet.Run(kitArgs)
}
