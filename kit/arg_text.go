package kit

import (
	"fmt"
	"os"

	"github.com/erikgeiser/promptkit/textinput"
	"github.com/logrusorgru/aurora"
)

var argTextTmpl = `
{{- Bold .Prompt }}
{{- if not .Valid }} {{- Foreground "1" (Bold "✘") }}
{{- else }} {{- Foreground "2" (Bold "✔") }}
{{- end }} {{ .Input -}}
`

func RunArgTextPrompt(kitArg KitArgument) (string, error) {
	var prefix string
	if kitArg.Required {
		prefix = aurora.Red("* ").String()
	}

	prompt := fmt.Sprintf("Select value for %v argument", kitArg.Name)
	if len(kitArg.Prompt) > 0 {
		prompt = kitArg.Prompt
	}

	input := textinput.New(prefix + prompt + "\n")
	input.Template = argTextTmpl
	input.ResultTemplate = ""
	input.Validate = func(s string) bool {
		if kitArg.Required {
			if len(s) == 0 {
				return false
			}
		}
		return true
	}

	name, err := input.RunPrompt()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	return name, nil
}
