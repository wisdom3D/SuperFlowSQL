package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"superflowsql-tui/core"
)

// ConfigViewModel displays and allows editing of project config
type ConfigViewModel struct {
	entries    []core.ConfigEntry
	cursor     int
	editing    bool
	editInput  string
	err        error
	projectDir string
	message    string
}

// NewConfigViewModel creates a new config view
func NewConfigViewModel(projectDir string) ConfigViewModel {
	entries, err := core.LoadConfig(projectDir)
	return ConfigViewModel{
		entries:    entries,
		projectDir: projectDir,
		err:        err,
	}
}

func (m ConfigViewModel) Init() tea.Cmd {
	return nil
}

func (m ConfigViewModel) Update(msg tea.Msg) (ConfigViewModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.editing {
			switch msg.String() {
			case "enter":
				// Save the edit
				if m.cursor < len(m.entries) {
					m.entries[m.cursor].Value = m.editInput
					err := core.SaveConfig(m.projectDir, m.entries)
					if err != nil {
						m.message = "Error saving: " + err.Error()
					} else {
						m.message = fmt.Sprintf("Updated %s", m.entries[m.cursor].Key)
					}
				}
				m.editing = false
				m.editInput = ""
			case "esc":
				m.editing = false
				m.editInput = ""
			case "backspace":
				if len(m.editInput) > 0 {
					m.editInput = m.editInput[:len(m.editInput)-1]
				}
			default:
				if len(msg.String()) == 1 || msg.String() == " " {
					m.editInput += msg.String()
				}
			}
			return m, nil
		}

		switch msg.String() {
		case "esc", "q":
			return m, func() tea.Msg { return FormCancelMsg{} }
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.entries)-1 {
				m.cursor++
			}
		case "enter", "e":
			if m.cursor < len(m.entries) {
				m.editing = true
				m.editInput = m.entries[m.cursor].Value
				m.message = ""
			}
		}
	}

	return m, nil
}

func (m ConfigViewModel) View() string {
	var b strings.Builder

	// Breadcrumb
	b.WriteString(BreadcrumbStyle.Render("Menu"))
	b.WriteString(BreadcrumbStyle.Render(" > "))
	b.WriteString(BreadcrumbActiveStyle.Render("Configuration"))
	b.WriteString("\n\n")

	b.WriteString(SubtitleStyle.Render("⚙  Configuration (.env)"))
	b.WriteString("\n\n")

	if m.err != nil {
		b.WriteString(ErrorStyle.Render("  ✗ " + m.err.Error()))
		b.WriteString("\n\n")
		b.WriteString(WarningStyle.Render("  Make sure you're in a SuperFlowSQL project directory."))
		b.WriteString("\n\n")
		b.WriteString(HelpStyle.Render("esc back"))
		return b.String()
	}

	if len(m.entries) == 0 {
		b.WriteString(WarningStyle.Render("  No configuration entries found."))
		b.WriteString("\n\n")
		b.WriteString(HelpStyle.Render("esc back"))
		return b.String()
	}

	// Group and display entries
	currentGroup := ""
	for i, entry := range m.entries {
		group := core.GetGroupPrefix(entry.Key)
		if group != currentGroup {
			currentGroup = group
			b.WriteString(SectionHeaderStyle.Render("  " + group))
			b.WriteString("\n")
		}

		if i == m.cursor {
			if m.editing {
				b.WriteString(fmt.Sprintf("  %s %s = ",
					InputLabelStyle.Copy().Foreground(ColorAccent).Render("❯"),
					InputLabelStyle.Render(entry.Key),
				))
				b.WriteString(InfoStyle.Render(m.editInput + "█"))
			} else {
				b.WriteString(fmt.Sprintf("  %s %s = %s",
					InputLabelStyle.Copy().Foreground(ColorAccent).Render("❯"),
					InputLabelStyle.Render(entry.Key),
					InfoStyle.Render(entry.Value),
				))
			}
		} else {
			b.WriteString(fmt.Sprintf("    %s = %s",
				TableCellStyle.Render(entry.Key),
				InfoStyle.Render(entry.Value),
			))
		}
		b.WriteString("\n")
	}

	// Status message
	if m.message != "" {
		b.WriteString("\n")
		b.WriteString(SuccessStyle.Render("  " + m.message))
	}

	b.WriteString("\n\n")
	if m.editing {
		b.WriteString(HelpStyle.Render("enter save • esc cancel edit"))
	} else {
		b.WriteString(HelpStyle.Render("↑/↓ navigate • enter/e edit • esc back"))
	}

	return b.String()
}
