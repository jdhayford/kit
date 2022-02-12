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

func (kc KitCommand) GetLineCount() int {
	return strings.Count(kc.Command, "\n") + 1
}

type commandFormatOptions struct {
	highlightArg      string
	includePrefixes   bool
	highlightPrefixes bool
}

func (kc *KitCommand) PromptArguments() string {
	var lastArg string
	for _, arg := range kc.GetArguments() {
		lastArg = arg.Name
		opts := commandFormatOptions{highlightArg: arg.Name, includePrefixes: true}
		commandPreview := kc.FormatCommand(opts)
		fmt.Println(commandPreview)

		var val string
		var err error

		switch arg.Type {
		case "select":
			val, err = RunArgSelectPrompt(arg)
		case "text":
			val, err = RunArgTextPrompt(arg)
		}
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
	// lineCount := kc.GetLineCount()

	// Single line prefix
	if options.includePrefixes {
		firstPrefixContent := "$ "
		// if lineCount > 1 {
		// 	firstPrefixContent = "0 "
		// }

		if options.highlightPrefixes {
			command = aurora.Green(firstPrefixContent).String() + command
		} else {
			command = aurora.Gray(8-1, firstPrefixContent).String() + command
		}

		// // Multi-line prefixes
		// var tempLines []string
		// rawLines := strings.Split(command, "\n")
		// for i, line := range rawLines {
		// 	if i == 0 {
		// 		tempLines = append(tempLines, line)
		// 		continue
		// 	}
		// 	num := aurora.Gray(8-1, fmt.Sprint(i)).String()
		// 	tempLines = append(tempLines, num+" "+line)
		// }
		// command = strings.Join(tempLines, "\n")
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
	Name     string                `yaml: "name"`
	Commands map[string]KitCommand `yaml:"commands"`

	// Note: Ref's should always be assigned to any User Kits that are loaded in
	// 		 Lack of a Ref indicates that a Kit is being used in context on execution
	Ref *KitRef
}

func newKit() Kit {
	return Kit{Commands: make(map[string]KitCommand)}
}

func (k Kit) IsEmpty() bool {
	emptyName := k.Name == ""
	emptyCommands := len(k.Commands) == 0
	emptyRef := k.Ref == nil
	return emptyName && emptyCommands && emptyRef
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

// KitSet represents a collection of kits for execution, with at most one PrimaryKit, and the rest being OtherKits
// 	Determination of a PrimaryKit goes as follows:
// 		If there is only one kit, it will be the PrimaryKit.
// 		If there are multiple kits, and one is deemed a context kit, it will be the PrimaryKit
// 		Otherwise, there will be no PrimaryKit
type KitSet struct {
	PrimaryKit Kit
	OtherKits  []Kit
}

func (ks KitSet) HasPrimaryKit() bool {
	return !ks.PrimaryKit.IsEmpty()
}

func MakeKitSet(kits ...Kit) KitSet {
	kitSet := KitSet{}
	for _, kit := range kits {
		if kit.Ref == nil && kitSet.PrimaryKit.IsEmpty() {
			kitSet.PrimaryKit = kit
		} else {
			kitSet.OtherKits = append(kitSet.OtherKits, kit)
		}
	}
	return kitSet
}

func (ks KitSet) GetCommands() []KitCommand {
	kitCommands := ks.PrimaryKit.GetCommands()
	for _, otherKit := range ks.OtherKits {
		for _, command := range otherKit.GetCommands() {
			if ks.HasPrimaryKit() {
				command.Alias = "@" + otherKit.Ref.Alias + " " + command.Alias
			}
			kitCommands = append(kitCommands, command)
		}
	}
	return kitCommands
}

// GetCommandFromArgs attempts to match the supplied args to an existing command the PrimaryKit or OtherKits
func (ks KitSet) GetCommandFromArgs(args []string) KitCommand {
	firstArg := args[0]

	// Check PrimaryKit for match
	for _, c := range ks.PrimaryKit.Commands {
		if firstArg == c.Alias {
			return c
		}
	}

	// No matching command found on PrimaryKit, check OtherKits
	for _, otherKit := range ks.OtherKits {
		for _, c := range otherKit.Commands {
			if firstArg == c.Alias {
				return c
			}
		}
	}

	fmt.Println("No kit command found for:", firstArg)
	os.Exit(1)

	return KitCommand{}
}

func (ks KitSet) Run(args []string) {
	var command KitCommand
	if len(args) > 0 {
		command = ks.GetCommandFromArgs(args)
	} else {
		command = RunCommandSelectPrompt(ks.GetCommands())
	}

	command.PromptArguments()

	commandPreview := command.FormatCommand(commandFormatOptions{
		includePrefixes:   true,
		highlightPrefixes: true,
	})
	fmt.Println(commandPreview)

	// if err := promptRunConfirmation(commandStr); err != nil {
	// 	os.Exit(1)
	// }
	runCommand(command.Alias, command.GenerateCommand())

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
