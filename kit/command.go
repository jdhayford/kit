package kit

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

func runCommand(command string) {
	args := strings.Split(command, " ")
	ctx, cancel := context.WithCancel(context.Background())

	cmd := exec.CommandContext(ctx, args[0], args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		cancel()
		fmt.Println(err)
		panic(err)
	}
	cancel()
}

func runCommandSilent(command string) string {
	args := strings.Split(command, " ")
	ctx, cancel := context.WithCancel(context.Background())
	var buff bytes.Buffer

	cmd := exec.CommandContext(ctx, args[0], args[1:]...)
	cmd.Stdout = &buff

	if err := cmd.Run(); err != nil {
		cancel()
		fmt.Println(err)
		panic(err)
	}
	cancel()

	bytes, _ := io.ReadAll(&buff)
	fmt.Println(string(bytes))
	return string(bytes)
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
