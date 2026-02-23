package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"superflowsql-tui/core"
)

// ActionModel handles stack start/stop operations with progress
type ActionModel struct {
	Title      string
	action     core.StackAction
	spinner    spinner.Model
	running    bool
	done       bool
	output     string
	err        error
	projectDir string
}

// ActionCompleteMsg is sent when the action finishes
type ActionCompleteMsg struct {
	Output string
	Err    error
}

// NewActionModel creates a new action view for start/stop
func NewActionModel(title string, action core.StackAction, projectDir string) ActionModel {
	s := spinner.New()
	s.Spinner = spinner.MiniDot
	s.Style = lipgloss.NewStyle().Foreground(ColorAccent)

	return ActionModel{
		Title:      title,
		action:     action,
		spinner:    s,
		running:    true,
		projectDir: projectDir,
	}
}

func (m ActionModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, m.runAction())
}

func (m ActionModel) runAction() tea.Cmd {
	return func() tea.Msg {
		var output string
		var err error

		switch m.action {
		case core.StackStart:
			output, err = core.StartStack(m.projectDir)
		case core.StackStop:
			output, err = core.StopStack(m.projectDir)
		}

		return ActionCompleteMsg{Output: output, Err: err}
	}
}

func (m ActionModel) Update(msg tea.Msg) (ActionModel, tea.Cmd) {
	switch msg := msg.(type) {
	case ActionCompleteMsg:
		m.running = false
		m.done = true
		m.output = msg.Output
		m.err = msg.Err

	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q", "enter":
			if m.done {
				return m, func() tea.Msg { return FormCancelMsg{} }
			}
		}

	case spinner.TickMsg:
		if m.running {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}
	}

	return m, nil
}

func (m ActionModel) View() string {
	var b strings.Builder

	// Breadcrumb
	b.WriteString(BreadcrumbStyle.Render("Menu"))
	b.WriteString(BreadcrumbStyle.Render(" > "))
	b.WriteString(BreadcrumbActiveStyle.Render(m.Title))
	b.WriteString("\n\n")

	b.WriteString(SubtitleStyle.Render(m.Title))
	b.WriteString("\n\n")

	if m.running {
		b.WriteString(fmt.Sprintf("  %s Running docker compose...\n", m.spinner.View()))
		b.WriteString("\n")
		b.WriteString(InfoStyle.Render("  This may take a few minutes for initial builds."))
		b.WriteString("\n")
		return b.String()
	}

	if m.err != nil {
		b.WriteString(ErrorStyle.Render("  ✗ Operation failed"))
		b.WriteString("\n\n")
		b.WriteString(ErrorStyle.Render("  " + m.err.Error()))
		b.WriteString("\n\n")
		if m.output != "" {
			b.WriteString(SectionHeaderStyle.Render("  Output:"))
			b.WriteString("\n")
			// Show last few lines of output
			lines := strings.Split(m.output, "\n")
			start := 0
			if len(lines) > 20 {
				start = len(lines) - 20
			}
			for _, line := range lines[start:] {
				b.WriteString("  " + InfoStyle.Render(line) + "\n")
			}
		}
	} else {
		b.WriteString(SuccessStyle.Render("  ✓ Operation completed successfully"))
		b.WriteString("\n\n")

		if m.action == core.StackStart {
			b.WriteString(SectionHeaderStyle.Render("  Service URLs:"))
			b.WriteString("\n")
			b.WriteString(InfoStyle.Render("  • Airflow:   http://localhost:8080"))
			b.WriteString("\n")
			b.WriteString(InfoStyle.Render("  • PgAdmin:   http://localhost:5050"))
			b.WriteString("\n")
			b.WriteString(InfoStyle.Render("  • Superset:  http://localhost:8088"))
			b.WriteString("\n")
			b.WriteString(InfoStyle.Render("  • Postgres:  localhost:5432"))
			b.WriteString("\n")
		}

		if m.output != "" {
			b.WriteString("\n")
			b.WriteString(SectionHeaderStyle.Render("  Output:"))
			b.WriteString("\n")
			lines := strings.Split(m.output, "\n")
			start := 0
			if len(lines) > 15 {
				start = len(lines) - 15
			}
			for _, line := range lines[start:] {
				if strings.TrimSpace(line) != "" {
					b.WriteString("  " + InfoStyle.Render(line) + "\n")
				}
			}
		}
	}

	b.WriteString("\n")
	b.WriteString(HelpStyle.Render("enter/esc back to menu"))

	return b.String()
}
