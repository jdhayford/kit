package kit

import (
	"fmt"
	"os"

	"github.com/erikgeiser/promptkit/textinput"
)

var cmdConfirmTmpl = `
Enter to run {{ Faint "(any other key to cancel)"  }}
`

func RunCommandConfirmPrompt(command string) (string, error) {
	input := textinput.New("")
	input.Template = argTextTmpl
	input.ResultTemplate = ""

	name, err := input.RunPrompt()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	return name, nil
}
