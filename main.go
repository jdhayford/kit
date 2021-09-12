package main

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/logrusorgru/aurora"
	"github.com/manifoldco/promptui"
	"gopkg.in/yaml.v3"
)

const kitFileName = "kit.yml"

var au aurora.Aurora

type KitCommand struct {
	Alias       string
	Command     string `yaml:"command"`
	Description string `yaml:"description"`
}

type Kit struct {
	Commands map[string]KitCommand `yaml:"commands"`
}

func newKit() Kit {
	return Kit{Commands: make(map[string]KitCommand)}
}

func (k Kit) GetCommands() []KitCommand {
	var commands []KitCommand
	for _, v := range k.Commands {
		commands = append(commands, v)
	}
	return commands
}

func (k Kit) GetCommand(index int) KitCommand {
	var commands []KitCommand
	for _, v := range k.Commands {
		commands = append(commands, v)
	}

	if index > len(commands)+1 {
		panic(errors.New("bad index for result"))
	}

	return commands[index]
}

var ctx context.Context
var cancel context.CancelFunc
var cmd *exec.Cmd

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func parseKitFile(filePath string) Kit {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		panic(errors.New("parseKitFile(): File doesnt exist"))
	}

	data, err := os.ReadFile(filePath)
	check(err)

	kit := newKit()
	yaml.Unmarshal(data, kit)
	check(err)

	// Assign key to command struct as Alias
	for k, v := range kit.Commands {
		v.Alias = k
		kit.Commands[k] = v
	}

	return kit
}

func findKitFile() (string, error) {
	filePath, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	for filePath != "/" {
		// Look for .kit file
		kitFilePath := path.Join(filePath, kitFileName)
		if _, err := os.Stat(kitFilePath); !os.IsNotExist(err) {
			return kitFilePath, nil
		}

		// Navigate to parent dir
		dirPath := path.Dir(filePath)
		filePath = dirPath
	}

	return "", errors.New("no .kit file found")
}

func promptCommandSelectionForKit(k Kit) KitCommand {
	kitCommands := k.GetCommands()

	templates := &promptui.SelectTemplates{
		Label:    `{{ "Select a command:" | faint }}`,
		Active:   "> {{ .Alias | cyan }} - {{ .Command | faint }}",
		Inactive: "  {{ .Alias | white | cyan }}",
		Selected: `{{ ">" | faint }} {{ .Alias | green }}`,
		Details: `
{{ "---------" | faint }} {{ .Alias | cyan }} {{ "---------" | faint }}
{{ "Command" | faint }} - {{ .Command | faint }}
{{ "Description" | faint }} - {{ .Description | faint }}`,
	}
	prompt := promptui.Select{
		Label:     "Kit",
		Templates: templates,
		Items:     kitCommands,
		HideHelp:  true,
		// HideSelected: true,
	}

	index, _, err := prompt.Run()
	if err != nil {
		fmt.Println(err)
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}

	return k.GetCommand(index)
}

func promptRunConfirmation(kitCommand KitCommand) error {
	fmt.Println(au.Gray(8-1, "$ ").String() + kitCommand.Command)

	promptTemplate := &promptui.PromptTemplates{
		Prompt: `Enter to run {{ "(any other key to cancel)" | faint }}`,
		Valid:  `Enter to run {{ "(any other key to cancel)" | faint }}`,
	}
	confirmPrompt := promptui.Prompt{
		Templates:   promptTemplate,
		HideEntered: true,
		Validate: func(s string) error {
			if len(s) != 0 {
				fmt.Println(" - Cancelling...")
				os.Exit(0)
			}
			return nil
		},
	}
	_, err := confirmPrompt.Run()
	if err != nil {
		return err
	}

	fmt.Println("\033[1A" + au.Green("$ ").String() + kitCommand.Command)
	return nil
}

func runCommand(command string) {
	args := strings.Split(command, " ")
	rest := strings.Join(args[1:], " ")

	ctx, cancel = context.WithCancel(context.Background())

	cmd = exec.CommandContext(ctx, args[0], rest)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		panic(err)
	}
}

var colors = flag.Bool("colors", true, "enable or disable colors")

func main() {
	au = aurora.NewAurora(*colors)
	kitFilePath, err := findKitFile()
	if err != nil {
		panic(err)
	}

	kit := parseKitFile(kitFilePath)
	kitCommand := promptCommandSelectionForKit(kit)

	if err := promptRunConfirmation(kitCommand); err != nil {
		os.Exit(1)
	}

	runCommand(kitCommand.Command)

	// Start goroutine to watch for sigint
	go func() {
		consolescanner := bufio.NewScanner(os.Stdin)
		if err := consolescanner.Err(); err != nil {
			os.Exit(1)
		}
	}()
}
