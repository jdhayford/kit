package kit

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/logrusorgru/aurora"
)

type KitArgument struct {
	Name          string   `yaml:"name"`
	Type          string   `yaml:"type"`
	Prompt        string   `yaml:"prompt"`
	Required      bool     `yaml:"required"`
	Options       []string `yaml:"options"`
	OptionCommand string   `yaml:"optionCommand"`
	OptionRegex   string   `yaml:"optionRegex"`
	// Might want to make value an interface for non string type args
	Value interface{}
}

type KitCommand struct {
	Alias       string
	Command     string `yaml:"command"`
	Description string `yaml:"description"`
	Arguments   map[string]KitArgument
}

func (kc KitCommand) GetArguments() []KitArgument {
	var arguments []KitArgument
	for _, v := range kc.Arguments {
		arguments = append(arguments, v)
	}
	return arguments
}

type commandFormatOptions struct {
	highlightArg string
	preview      bool
	execute      bool
}

func (kc *KitCommand) PromptArguments() string {
	var lastArg string
	for _, arg := range kc.GetArguments() {
		lastArg = arg.Name
		opts := commandFormatOptions{highlightArg: arg.Name, preview: true}
		commandPreview := kc.FormatCommand(opts)
		fmt.Println(commandPreview)
		val, err := promptArgument(arg)
		if err != nil {
			panic(err)
		}
		arg.Value = val
		kc.Arguments[arg.Name] = arg

		clearLastLine()
	}
	return lastArg
}

func (kc KitCommand) FormatCommand(options commandFormatOptions) string {
	// Parse command, looking for template spots, replace with arguments found
	command := kc.Command

	if options.execute {
		command = aurora.Green("$ ").String() + command
	} else if options.preview {
		command = aurora.Gray(8-1, "$ ").String() + command
	}

	for _, v := range kc.Arguments {
		placeholder := "@" + v.Name
		if v.Value == nil {
			if v.Name == options.highlightArg {
				focusedPlaceholder := aurora.Cyan(placeholder).String()
				command = strings.ReplaceAll(command, placeholder, focusedPlaceholder)
			}
			continue
		}
		// Make types actual types
		if v.Type == "text" || v.Type == "select" {
			textVal := v.Value.(string)
			if v.Name == options.highlightArg {
				textVal = aurora.Cyan(textVal).String()
			}
			command = strings.ReplaceAll(command, placeholder, textVal)
		}
	}
	return command
}

func (kc KitCommand) GenerateCommand() string {
	return kc.FormatCommand(commandFormatOptions{})
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

func (k Kit) GetCommandFromArgs(args []string) KitCommand {
	var kitCommand KitCommand
	firstArg := args[0]
	for _, c := range k.Commands {
		if firstArg == c.Alias {
			kitCommand = c
		}
	}

	if kitCommand.Alias == "" {
		fmt.Println("No kit command found for:", firstArg)
		os.Exit(1)
	}

	return kitCommand
}

func (k Kit) GetCommandFromPrompt() KitCommand {
	kitCommand := promptCommandSelectionForKit(k)

	lastArg := kitCommand.PromptArguments()

	commandPreview := kitCommand.FormatCommand(commandFormatOptions{
		highlightArg: lastArg,
		preview:      true,
	})
	fmt.Println(commandPreview)

	clearLastNLines(1)
	commandPreview = kitCommand.FormatCommand(commandFormatOptions{
		execute: true,
	})
	fmt.Println(commandPreview)
	return kitCommand
}

func (k Kit) Run(args []string) {
	var command KitCommand
	if len(args) > 0 {
		command = k.GetCommandFromArgs(args)
	} else {
		command = promptCommandSelectionForKit(k)
	}

	lastArg := command.PromptArguments()

	commandPreview := command.FormatCommand(commandFormatOptions{
		highlightArg: lastArg,
		preview:      true,
	})
	fmt.Println(commandPreview)

	clearLastNLines(1)
	commandPreview = command.FormatCommand(commandFormatOptions{
		execute: true,
	})
	fmt.Println(commandPreview)

	// if err := promptRunConfirmation(commandStr); err != nil {
	// 	os.Exit(1)
	// }
	runCommand(command.GenerateCommand())

	// home, err := os.UserHomeDir()
	// check(err)
	// historyPath := path.Join(home, ".zsh_history")
	// f, err := os.OpenFile(historyPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	// check(err)

	// defer f.Close()

	// _, err = f.WriteString(commandStr + "\n")
	// fmt.Println(err)
	// check(err)

	// runCommand(`/bin/bash -c "history -a"`)

	// Start goroutine to watch for sigint
	go func() {
		consolescanner := bufio.NewScanner(os.Stdin)
		if err := consolescanner.Err(); err != nil {
			os.Exit(1)
		}
	}()
}
