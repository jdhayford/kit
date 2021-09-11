package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// type client chan<- string // an outgoing message channel
// var (
// 	messages = make(chan string) // all incoming client messages
// 	entering = make(chan client) // entering clients
// 	leaving  = make(chan client) // leaving clients
// )

var ctx context.Context
var cancel context.CancelFunc
var cmd *exec.Cmd

func runCommand(args []string, errChan chan bool) {
	firstArg := args[0]
	otherArgs := []string{}
	if len(args) > 1 {
		otherArgs = args[1:]
	}

	ctx, cancel = context.WithCancel(context.Background())

	cmd = exec.CommandContext(ctx, firstArg, otherArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		panic(err)
	}

	go func() {
		cmdErr := cmd.Wait()
		// TODO: Based on config, either exit here or nah
		if cmdErr != nil && cmdErr.Error() != "signal: killed" {
			// errChan <- true
			fmt.Printf("\nCommand exited with error, hit enter to try again\n")
		}
	}()
}

func main() {
	if len(os.Args) == 1 {
		fmt.Println("Usage: stopgo [some command]")
		return
	}

	errChan := make(chan bool)

	commandArgs := os.Args[1:]

	runCommand(commandArgs, errChan)

	go func() {
		consolescanner := bufio.NewScanner(os.Stdin)

		for consolescanner.Scan() {
			input := consolescanner.Text()
			if len(input) == 0 {
				cancel()
				fmt.Printf("> %v\n", strings.Join(commandArgs, " "))
				runCommand(commandArgs, errChan)
			}
		}

		if err := consolescanner.Err(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}()

	<-errChan
}
