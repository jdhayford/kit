package kit

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/erikgeiser/promptkit/selection"
	"github.com/logrusorgru/aurora"
	"github.com/muesli/termenv"
)

type CmdSelectModel struct {
	commands        []KitCommand
	selectedCommand *KitCommand
	selection       selection.Model
	err             error
}

func newCmdSelectModel(commands []KitCommand) *CmdSelectModel {
	return &CmdSelectModel{commands: commands}
}

var temp = `
{{- if .Prompt -}}
  {{ Bold .Prompt }}
{{ end -}}
{{ if .IsFiltered }}
  {{- print .FilterPrompt " " .FilterInput }}
{{ end }}

{{- if and .FilterInput (eq (len .Choices) 0) -}}
	{{- Faint "    No matching choices\n" }}
{{- end -}}

{{- range  $i, $choice := .Choices }}
  {{- if IsScrollUpHintPosition $i }}
    {{- "⇡ " -}}
  {{- else if IsScrollDownHintPosition $i -}}
    {{- "⇣ " -}} 
  {{- else -}}
    {{- "  " -}}
  {{- end -}} 

  {{- if eq $.SelectedIndex $i }}
   {{- print ("> ") (Selected $choice) "\n" }}
  {{- else }}
    {{- print "  " (Unselected $choice) "\n" }}
  {{- end }}

  {{- end}}

{{- if gt (len .Choices) 0 }}
{{- with (index .Choices $.SelectedIndex) }}
{{- with .Value }}
{{ Faint "---------" }} {{ Foreground "#63FFFF" .Alias }} {{ Faint "---------" }}
{{ Faint "Command"}} - {{ Faint .Command }}
{{ Faint "Description"}} - {{ Faint .Description }}
{{- end }}
{{- end }}
{{- end }}
`

func customFilter(filter string, choice *selection.Choice) bool {
	commandChoice, _ := choice.Value.(KitCommand)
	command := strings.ToLower(strings.Trim(commandChoice.Alias, " "))
	filter = strings.ToLower(strings.Trim(filter, " "))
	return strings.HasPrefix(command, filter)
}

func (s *CmdSelectModel) Init() tea.Cmd {
	sel := selection.New("",
		selection.Choices(s.commands))
	sel.Template = temp
	sel.Filter = customFilter
	sel.FilterPrompt = ""
	sel.SelectedChoiceStyle = func(c *selection.Choice) string {
		choice, _ := c.Value.(KitCommand)
		return aurora.Cyan(choice.Alias).String() + termenv.String(" - "+choice.Command).Faint().String()
	}
	sel.UnselectedChoiceStyle = func(c *selection.Choice) string {
		choice, _ := c.Value.(KitCommand)
		return choice.Alias
	}

	s.selection = *selection.NewModel(sel)

	return s.selection.Init()
}

func (s *CmdSelectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return s, nil
	}

	switch {
	case keyMsg.String() == "enter":
		c, err := s.selection.Value()
		if err != nil {
			if c == nil {
				return s, nil
			}
			s.err = err
			return s, tea.Quit
		}

		command, _ := c.Value.(KitCommand)
		s.selectedCommand = &command
		return s, tea.Quit
	case keyMsg.String() == "esc":
		return s, tea.Quit
	default:
		_, cmd := s.selection.Update(msg)

		return s, cmd
	}
}

func (s *CmdSelectModel) View() string {
	if s.err != nil {
		return fmt.Sprintf("Error: %v", s.err)
	}

	if s.selectedCommand != nil {
		return ""
	}

	var b strings.Builder

	b.WriteString(s.selection.View())

	return b.String()
}

func RunCommandSelectPrompt(commands []KitCommand) KitCommand {
	model := newCmdSelectModel(commands)

	p := tea.NewProgram(model)

	returnedModel, err := p.StartReturningModel()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v", err)

		os.Exit(1)
	}

	return *returnedModel.(*CmdSelectModel).selectedCommand
}
