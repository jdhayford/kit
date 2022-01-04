package kit

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func runCommand(command string) {
	args := strings.Split(command, " ")
	ctx, cancel := context.WithCancel(context.Background())

	cmd := exec.CommandContext(ctx, args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		cancel()
		fmt.Println(err)
		panic(err)
	}
	cancel()
}

// func runCommandList(args []string) {
// 	rest := strings.Join(args[1:], " ")
// 	fmt.Println(rest)

// 	ctx, cancel := context.WithCancel(context.Background())

// 	cmd := exec.CommandContext(ctx, args[0], rest)
// 	cmd.Stdout = os.Stdout
// 	cmd.Stderr = os.Stderr

// 	if err := cmd.Run(); err != nil {
// 		cancel()
// 		fmt.Println(err)
// 		panic(err)
// 	}
// 	cancel()
// }
