package kit

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"strings"
)

func runCommand(name string, command string) {
	// Write command to temp file in home kit
	// Defer deletion of temp file (after tested)
	// Execute bash against file

	// Create new temporary command file
	kitExecDir := GetOrMakeKitExecDir()
	fileName := path.Join(kitExecDir, name)
	// tmpCmdFile, err := ioutil.TempFile(kitDir, "cmd-")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer os.Remove(tmpCmdFile.Name())

	// Write user command content to command file
	tmpCmdBytes := []byte(command)
	err := os.WriteFile(fileName, tmpCmdBytes, 0644)
	check(err)

	// Todo: Make this dynamic?
	shell := "bash"

	ctx, cancel := context.WithCancel(context.Background())
	cmd := exec.CommandContext(ctx, shell, fileName)
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
