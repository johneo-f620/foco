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

type row struct {
	indent  int
	content string
	tags    []string
}

type model struct {
	title    string // TODO: Display title at the top
	items    map[int]row
	viewport viewport.Model
	textarea textarea.Model
	ready    bool
}

var builtinTags = []string{"problem", "steps", "facts", "questions"}

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
		items: map[int]row{
			0: {indent: 0, content: "An idea!", tags: []string{"idea", "new stuff"}},
			1: {indent: 1, content: "Step 1: Do something", tags: []string{"step"}},
			2: {indent: 2, content: "Step 1.1: And then some", tags: []string{"step"}},
			3: {indent: 1, content: "Step 2: Continue", tags: []string{"step"}},
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
			m.items[len(m.items)] = newRow
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
	tagStart := -1
	for i := 0; i < len(rn); i++ {
		if rn[i] == '#' {
			tagStart = i
			break
		}
	}
	var (
		tagsSection []rune
		tags        []string
		content     string
	)
	if tagStart > -1 {
		tagsSection = rn[tagStart:]
		tags = strings.Split(string(tagsSection), "#")
		content = string(rn[0 : tagStart-1])
	} else {
		// No tags
		tags = []string{}
		content = m.textarea.Value()
	}
	newRow := row{
		indent:  0,
		content: content,
		tags:    tags,
	}
	return newRow
}

func renderContent(m model) string {
	s := ""
	for i := 0; i < len(m.items); i++ {
		item := m.items[i]
		tags := []string{}
		for _, t := range item.tags {
			tags = append(tags, "#"+t)
		}
		indent := strings.Repeat("  ", item.indent)
		s += fmt.Sprintf("%s%s %s\n", indent, item.content, strings.Join(tags, " "))
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
