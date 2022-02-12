package kit

import (
	"fmt"
	"os"

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
