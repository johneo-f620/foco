package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"os"
)

func initialModel() model {
	ta := textarea.New()
	ta.Placeholder = "Enter an entry or a command"
	ta.Focus()

	ta.Prompt = "> "
	ta.CharLimit = 200
	ta.SetHeight(2)
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()

	ta.ShowLineNumbers = false

	ta.KeyMap.InsertNewline.SetEnabled(false)

	initModel := model{
		Textarea: ta,
		Items: map[string][]row{
			"questions": {row{Indent: 0, Content: "An idea!", Header: "questions"}},
			"steps": {
				row{Indent: 1, Content: "Step 1: Do something", Header: "steps"},
				row{Indent: 2, Content: "Step 1.1: And then some", Header: "steps"},
				row{Indent: 1, Content: "Step 2: Continue", Header: "steps"},
			},
		},
	}

	return initModel
}

func main() {
	if _, err := tea.NewProgram(initialModel(), tea.WithAltScreen()).Run(); err != nil {
		fmt.Printf("Error! %v", err)
		os.Exit(1)
	}
}
