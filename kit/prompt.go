package kit

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/manifoldco/promptui"
)

type argSelectItem struct {
	Original string
	HasMatch bool
	Pre      string
	Match    string
	Post     string
}

// bellSkipper implements an io.WriteCloser that skips the terminal bell
// character (ASCII code 7), and writes the rest to os.Stderr. It is used to
// replace readline.Stdout, that is the package used by promptui to display the
// prompts.
//
// This is a workaround for the bell issue documented in
// https://github.com/manifoldco/promptui/issues/49.
type bellSkipper struct{}

// Write implements an io.WriterCloser over os.Stderr, but it skips the terminal
// bell character.
func (bs *bellSkipper) Write(b []byte) (int, error) {
	const charBell = 7 // c.f. readline.CharBell
	if len(b) == 1 && b[0] == charBell {
		return 0, nil
	}
	return os.Stderr.Write(b)
}

// Close implements an io.WriterCloser over os.Stderr.
func (bs *bellSkipper) Close() error {
	return os.Stderr.Close()
}

func promptCommandSelectionForKit(k Kit) KitCommand {
	kitCommands := k.GetCommands()

	searcher := func(input string, index int) bool {
		kitCommand := kitCommands[index]
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		alias := strings.Replace(strings.ToLower(kitCommand.Alias), " ", "", -1)
		command := strings.Replace(strings.ToLower(kitCommand.Command), " ", "", -1)
		desc := strings.Replace(strings.ToLower(kitCommand.Description), " ", "", -1)

		var match bool
		for _, target := range []string{alias, command, desc} {
			if !match {
				match = strings.Contains(target, input)
			}
		}

		return match
	}

	templates := &promptui.SelectTemplates{
		// Label:    `{{ "Select a command:" | faint }}`,
		Active:   "> {{ .Alias | cyan }} - {{ .Command | faint }}",
		Inactive: "  {{ .Alias | white | cyan }}",
		Selected: `{{ ">" | faint }} {{ .Alias | cyan }}`,
		Details: `
{{ "---------" | faint }} {{ .Alias | cyan }} {{ "---------" | faint }}
{{ "Command" | faint }} - {{ .Command | faint }}
{{ "Description" | faint }} - {{ .Description | faint }}`,
	}
	prompt := promptui.Select{
		Label:             "Kit",
		Templates:         templates,
		Items:             kitCommands,
		HideHelp:          true,
		Searcher:          searcher,
		StartInSearchMode: true,
		Stdout:            &bellSkipper{},
		// HideSelected: true,
	}

	index, _, err := prompt.Run()
	if err != nil {
		fmt.Println(err)
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}

	return kitCommands[index]
}

func promptRunConfirmation(command string) error {
	// fmt.Println(aurora.Gray(8-1, "$ ").String() + command)

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

	// fmt.Println("\033[1A" + aurora.Green("$ ").String() + command)
	return nil
}

func promptSelectArgument(kitArg KitArgument) (string, error) {
	items := kitArg.Options

	shouldGenerateItems := len(items) == 0 && len(kitArg.OptionCommand) > 0
	if shouldGenerateItems {
		rawOut := runCommandSilent(kitArg.OptionCommand)
		items = strings.Split(rawOut, "\n")

	}

	var displayItems []argSelectItem
	for _, item := range items {
		displayItem := argSelectItem{Original: item}
		if len(kitArg.OptionRegex) > 0 {
			reg := regexp.MustCompile(kitArg.OptionRegex)
			loc := reg.FindIndex([]byte(item))
			if loc != nil {
				displayItem.HasMatch = true
				displayItem.Pre = item[0:loc[0]]
				displayItem.Match = item[loc[0]:loc[1]]
				displayItem.Post = item[loc[1] : len(item)-1]
			}
		}

		displayItems = append(displayItems, displayItem)
	}

	templates := &promptui.SelectTemplates{
		Inactive: "  {{ .Original }}",
		// Active:   promptui.IconSelect + " {{ .HasMatch | underline | faint }}",
		Active: promptui.IconSelect + " {{if .HasMatch}}{{ .Pre | underline }}{{ .Match | underline | cyan }}{{ .Post | underline }}{{else}}{{ .Original | underline }}{{end}}",
	}
	prompt := promptui.Select{
		Templates:    templates,
		Label:        fmt.Sprintf("Select value for %v argument", kitArg.Name),
		Items:        displayItems,
		HideHelp:     true,
		Stdout:       &bellSkipper{},
		HideSelected: true,
	}

	index, _, err := prompt.Run()
	if err != nil {
		return "", err
	}

	value := items[index]
	if shouldGenerateItems && len(kitArg.OptionRegex) > 0 {
		reg := regexp.MustCompile(kitArg.OptionRegex)
		match := reg.Find([]byte(value))
		if match == nil {
			fmt.Println("Unable to match selected value using option pattern")
			os.Exit(1)
		}
		value = string(match)
	}

	return value, nil
}

func promptArgument(kitArg KitArgument) (string, error) {
	if kitArg.Type == "select" {
		return promptSelectArgument(kitArg)
	}

	promptTemplate := &promptui.PromptTemplates{}

	var preLabel string
	if kitArg.Required {
		preLabel = "* "
	}
	confirmPrompt := promptui.Prompt{
		Label:       preLabel + kitArg.Prompt,
		Templates:   promptTemplate,
		HideEntered: true,
		Validate: func(s string) error {
			if kitArg.Required {
				if len(s) == 0 {
					return errors.New("required argument cannot be empty")
				}
			}
			return nil
		},
	}
	val, err := confirmPrompt.Run()
	return val, err
}
