package kit

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/erikgeiser/promptkit/selection"
	"github.com/logrusorgru/aurora"
)

type ArgSelectModel struct {
	generatedItems bool
	argument       KitArgument
	choices        []string
	selectedChoice string
	selection      selection.Model
	err            error
}

func newArgSelectModel(kitArg KitArgument) *ArgSelectModel {
	return &ArgSelectModel{argument: kitArg}
}

var argSelectTemp = `
{{- if .Prompt -}}
  {{ Bold .Prompt }}
{{ end -}}
{{ if .IsFiltered }}
  {{- print .FilterPrompt " " .FilterInput }}
{{ end }}

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
`

func (s *ArgSelectModel) formatSelectedChoiceStyle(c *selection.Choice) string {
	var (
		hasMatch bool
		pre      string
		match    string
		post     string
	)
	choice, _ := c.Value.(string)

	if len(s.argument.OptionRegex) > 0 {
		reg := regexp.MustCompile(s.argument.OptionRegex)
		loc := reg.FindIndex([]byte(choice))
		if loc != nil {
			hasMatch = true
			pre = choice[0:loc[0]]
			match = choice[loc[0]:loc[1]]
			post = choice[loc[1]:]
		}
	} else {
		return aurora.Underline(choice).Cyan().String()
	}

	if hasMatch {
		return aurora.Underline(pre).String() + aurora.Underline(match).Cyan().String() + aurora.Underline(post).String()
	} else {
		return choice
	}
}

func (s *ArgSelectModel) Init() tea.Cmd {
	items := s.argument.Options

	if len(items) == 0 && len(s.argument.OptionCommand) > 0 {
		rawOut := runCommandSilent(s.argument.OptionCommand)
		items = strings.Split(rawOut, "\n")
		s.generatedItems = true
	}

	sel := selection.New(fmt.Sprintf("Select value for %v argument", s.argument.Name),
		selection.Choices(items))

	sel.Template = argSelectTemp
	sel.Filter = nil
	sel.FilterPrompt = ""
	sel.SelectedChoiceStyle = s.formatSelectedChoiceStyle
	sel.UnselectedChoiceStyle = func(c *selection.Choice) string {
		return c.Value.(string)
	}

	s.selection = *selection.NewModel(sel)

	return s.selection.Init()
}

func (s *ArgSelectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return s, nil
	}

	switch {
	case keyMsg.String() == "enter":
		c, err := s.selection.Value()
		if err != nil {
			s.err = err
			return s, tea.Quit
		}

		s.selectedChoice = c.Value.(string)
		return s, tea.Quit
	case keyMsg.String() == "esc":
		return s, tea.Quit
	default:
		_, cmd := s.selection.Update(msg)

		return s, cmd
	}
}

func (s *ArgSelectModel) View() string {
	if s.err != nil {
		return fmt.Sprintf("Error: %v", s.err)
	}

	if s.selectedChoice != "" {
		return ""
	}

	var b strings.Builder

	b.WriteString(s.selection.View())
	return b.String()
}

func RunArgSelectPrompt(kitArg KitArgument) (string, error) {
	model := newArgSelectModel(kitArg)

	p := tea.NewProgram(model)

	returnedModel, err := p.StartReturningModel()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v", err)

		os.Exit(1)
	}

	finalModel := returnedModel.(*ArgSelectModel)
	choice := finalModel.selectedChoice
	if finalModel.generatedItems && (len(finalModel.argument.OptionRegex) > 0) {
		reg := regexp.MustCompile(kitArg.OptionRegex)
		match := reg.Find([]byte(choice))
		if match == nil {
			fmt.Println("Unable to match selected value using option pattern")
			os.Exit(1)
		}
		choice = string(match)
	}

	fmt.Println(choice)
	return choice, nil
}
