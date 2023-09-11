package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"os"
	"strings"
)

func builtinHeaders() []string {
	return []string{"problem", "steps", "facts", "questions", "theories"}
}

type row struct {
	indent  int
	content string
	header  string
}

type model struct {
	title    string // TODO: Display title at the top
	items    map[string][]row
	viewport viewport.Model
	textarea textarea.Model
	ready    bool
}

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
		textarea: ta,
		items: map[string][]row{
			"questions": {row{indent: 0, content: "An idea!", header: "questions"}},
			"steps": {
				row{indent: 1, content: "Step 1: Do something", header: "steps"},
				row{indent: 2, content: "Step 1.1: And then some", header: "steps"},
				row{indent: 1, content: "Step 2: Continue", header: "steps"},
			},
		},
	}

	return initModel
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmds []tea.Cmd

		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	switch msg := message.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyEnter:
			newRow := makeNewRow(m)
			headerRows := m.items[newRow.header]
			headerRows = append(headerRows, newRow)
			m.viewport.SetContent(renderContent(m))
			m.textarea.Reset()
			m.viewport.GotoBottom()
		}
	case tea.WindowSizeMsg:
		if !m.ready {
			m.viewport = viewport.New(msg.Width, msg.Height-lipgloss.Height(m.textarea.Value()))
			m.viewport.YPosition = 0
			m.viewport.SetContent(renderContent(m))
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - m.textarea.Height()
		}
	}

	// Handle keyboard and mouse events in the viewport
	m.textarea, tiCmd = m.textarea.Update(message)
	m.viewport, vpCmd = m.viewport.Update(message)
	cmds = append(cmds, tiCmd, vpCmd)

	return m, tea.Batch(cmds...)
}

func makeNewRow(m model) row {
	rn := []rune(m.textarea.Value())
	hIndex := -1
	for i := 0; i < len(rn); i++ {
		if rn[i] == '#' {
			hIndex = i
			break
		}
	}
	var (
		header  []rune
		content string
	)
	if hIndex > -1 {
		header = rn[hIndex:]
	} else {
		// No header
		content = m.textarea.Value()
	}
	newRow := row{
		indent:  0,
		content: content,
		header:  string(header),
	}
	return newRow
}

func renderContent(m model) string {
	s := ""
	for j := 0; j < len(builtinHeaders()); j++ {
		header := builtinHeaders()[j]
		for i := 0; i < len(m.items[header]); i++ {
			row := m.items[header][i]
			indent := strings.Repeat("  ", row.indent)
			s += fmt.Sprintf("%s%s %s\n", indent, row.content, row.header)
		}
	}

	return s
}

func (m model) View() string {
	if m.ready {
		return fmt.Sprintf(
			"%s\n%s",
			m.viewport.View(),
			m.textarea.View(),
		)
	}
	return "\n Initializing..."
}

func main() {
	if _, err := tea.NewProgram(initialModel(), tea.WithAltScreen()).Run(); err != nil {
		fmt.Printf("Error! %v", err)
		os.Exit(1)
	}
}
