package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// FormField represents a single form input
type FormField struct {
	Label       string
	Placeholder string
	Value       string
	IsPassword  bool
}

// FormModel handles multi-field form input
type FormModel struct {
	Title      string
	Fields     []FormField
	inputs     []textinput.Model
	focusIndex int
	submitted  bool
	cancelled  bool
	width      int
	height     int
}

// FormSubmitMsg is sent when the form is submitted
type FormSubmitMsg struct {
	Values map[string]string
}

// FormCancelMsg is sent when the form is cancelled
type FormCancelMsg struct{}

// NewFormModel creates a new form with the given fields
func NewFormModel(title string, fields []FormField) FormModel {
	inputs := make([]textinput.Model, len(fields))

	for i, field := range fields {
		t := textinput.New()
		t.Placeholder = field.Placeholder
		t.CharLimit = 256

		if field.Value != "" {
			t.SetValue(field.Value)
		}

		if field.IsPassword {
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
		}

		if i == 0 {
			t.Focus()
		}

		inputs[i] = t
	}

	return FormModel{
		Title:  title,
		Fields: fields,
		inputs: inputs,
	}
}

func (m FormModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m FormModel) Update(msg tea.Msg) (FormModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.cancelled = true
			return m, func() tea.Msg { return FormCancelMsg{} }

		case "tab", "down":
			m.focusIndex++
			if m.focusIndex >= len(m.inputs) {
				m.focusIndex = 0
			}
			return m, m.updateFocus()

		case "shift+tab", "up":
			m.focusIndex--
			if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs) - 1
			}
			return m, m.updateFocus()

		case "enter":
			if m.focusIndex == len(m.inputs)-1 {
				// Submit on last field
				m.submitted = true
				values := make(map[string]string)
				for i, field := range m.Fields {
					val := m.inputs[i].Value()
					if val == "" {
						val = field.Placeholder
					}
					values[field.Label] = val
				}
				return m, func() tea.Msg { return FormSubmitMsg{Values: values} }
			}
			// Move to next field
			m.focusIndex++
			return m, m.updateFocus()
		}
	}

	// Update the focused input
	var cmd tea.Cmd
	m.inputs[m.focusIndex], cmd = m.inputs[m.focusIndex].Update(msg)
	return m, cmd
}

func (m FormModel) updateFocus() tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))
	for i := range m.inputs {
		if i == m.focusIndex {
			cmds[i] = m.inputs[i].Focus()
		} else {
			m.inputs[i].Blur()
		}
	}
	return tea.Batch(cmds...)
}

func (m FormModel) View() string {
	var b strings.Builder

	// Breadcrumb
	b.WriteString(BreadcrumbStyle.Render("Menu"))
	b.WriteString(BreadcrumbStyle.Render(" > "))
	b.WriteString(BreadcrumbActiveStyle.Render(m.Title))
	b.WriteString("\n\n")

	// Title
	b.WriteString(SubtitleStyle.Render(m.Title))
	b.WriteString("\n\n")

	// Form fields
	for i, field := range m.Fields {
		// Label
		label := InputLabelStyle.Render(field.Label + ":")
		if i == m.focusIndex {
			label = InputLabelStyle.Copy().
				Foreground(ColorAccent).
				Render("❯ " + field.Label + ":")
		}
		b.WriteString(label)
		b.WriteString("\n")

		// Input
		b.WriteString("  " + m.inputs[i].View())
		b.WriteString("\n\n")
	}

	// Progress indicator
	progress := fmt.Sprintf("  Field %d/%d", m.focusIndex+1, len(m.inputs))
	b.WriteString(InfoStyle.Render(progress))
	b.WriteString("\n\n")

	// Help
	b.WriteString(HelpStyle.Render("tab/↓ next • shift+tab/↑ prev • enter submit (on last field) • esc cancel"))

	return b.String()
}

// GetValues returns the current form values
func (m FormModel) GetValues() map[string]string {
	values := make(map[string]string)
	for i, field := range m.Fields {
		val := m.inputs[i].Value()
		if val == "" {
			val = field.Placeholder
		}
		values[field.Label] = val
	}
	return values
}
