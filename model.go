package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"strings"
)

type row struct {
	Indent  int
	Content string
	Header  string
}

type model struct {
	Title    string // TODO: Display Title at the top
	Items    map[string][]row
	Viewport viewport.Model
	Textarea textarea.Model
	Ready    bool
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
			headerRows := m.Items[newRow.Header]
			headerRows = append(headerRows, newRow)
			m.Items[newRow.Header] = headerRows
			m.Viewport.SetContent(renderContent(m))
			m.Textarea.Reset()
			m.Viewport.GotoBottom()
		}
	case tea.WindowSizeMsg:
		if !m.Ready {
			m.Viewport = viewport.New(msg.Width, msg.Height-lipgloss.Height(m.Textarea.Value()))
			m.Viewport.YPosition = 0
			m.Viewport.SetContent(renderContent(m))
			m.Ready = true
		} else {
			m.Viewport.Width = msg.Width
			m.Viewport.Height = msg.Height - m.Textarea.Height()
		}
	}

	// Handle keyboard and mouse events in the viewport
	m.Textarea, tiCmd = m.Textarea.Update(message)
	m.Viewport, vpCmd = m.Viewport.Update(message)
	cmds = append(cmds, tiCmd, vpCmd)

	return m, tea.Batch(cmds...)
}

func makeNewRow(m model) row {
	rn := []rune(m.Textarea.Value())
	hIndex := -1
	for i := 0; i < len(rn); i++ {
		if rn[i] == '#' {
			hIndex = i
			break
		}
	}
	var (
		header  string
		content string
	)
	if hIndex > -1 {
		header = string(rn[hIndex:])
		content = string(rn[0:hIndex])
	} else {
		// No header
		header = "facts"
		content = m.Textarea.Value()
	}
	newRow := row{
		Indent:  0,
		Content: content,
		Header:  header,
	}
	return newRow
}

func renderContent(m model) string {
	s := ""
	for j := 0; j < len(builtinHeaders()); j++ {
		header := builtinHeaders()[j]
		s += fmt.Sprintf("## %s\n\n", header)
		for i := 0; i < len(m.Items[header]); i++ {
			row := m.Items[header][i]
			indent := strings.Repeat("  ", row.Indent)
			s += fmt.Sprintf("%s- %s\n", indent, row.Content)
		}
		s += "\n"
	}

	return s
}

func (m model) View() string {
	if m.Ready {
		return fmt.Sprintf(
			"%s\n%s",
			m.Viewport.View(),
			m.Textarea.View(),
		)
	}
	return "\n Initializing..."
}

func builtinHeaders() []string {
	return []string{"problem", "steps", "facts", "questions", "theories"}
}
