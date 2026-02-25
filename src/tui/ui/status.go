package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"superflowsql-tui/core"
)

// StatusModel displays service status information
type StatusModel struct {
	services   []core.ServiceStatus
	spinner    spinner.Model
	loading    bool
	err        error
	width      int
	height     int
	projectDir string
}

// StatusLoadedMsg is sent when status data is fetched
type StatusLoadedMsg struct {
	Services []core.ServiceStatus
	Err      error
}

// NewStatusModel creates a new status view
func NewStatusModel(projectDir string) StatusModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(ColorAccent)

	return StatusModel{
		spinner:    s,
		loading:    true,
		projectDir: projectDir,
	}
}

func (m StatusModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, m.fetchStatus())
}

func (m StatusModel) fetchStatus() tea.Cmd {
	return func() tea.Msg {
		services, err := core.GetStackStatus(m.projectDir)
		return StatusLoadedMsg{Services: services, Err: err}
	}
}

func (m StatusModel) Update(msg tea.Msg) (StatusModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case StatusLoadedMsg:
		m.loading = false
		m.services = msg.Services
		m.err = msg.Err

	case tea.KeyMsg:
		switch msg.String() {
		case "r":
			m.loading = true
			m.err = nil
			return m, tea.Batch(m.spinner.Tick, m.fetchStatus())
		case "esc", "q":
			return m, func() tea.Msg { return FormCancelMsg{} }
		}

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m StatusModel) View() string {
	var b strings.Builder

	// Breadcrumb
	b.WriteString(BreadcrumbStyle.Render("Menu"))
	b.WriteString(BreadcrumbStyle.Render(" > "))
	b.WriteString(BreadcrumbActiveStyle.Render("Stack Status"))
	b.WriteString("\n\n")

	b.WriteString(SubtitleStyle.Render("📊 Stack Status"))
	b.WriteString("\n\n")

	if m.loading {
		b.WriteString(fmt.Sprintf("  %s Fetching service status...\n", m.spinner.View()))
		return b.String()
	}

	if m.err != nil {
		b.WriteString(ErrorStyle.Render("  ✗ Error: " + m.err.Error()))
		b.WriteString("\n\n")
		b.WriteString(WarningStyle.Render("  Make sure Docker is running and you're in a SuperFlowSQL project directory."))
		b.WriteString("\n\n")
		b.WriteString(HelpStyle.Render("r refresh • esc back"))
		return b.String()
	}

	if len(m.services) == 0 {
		b.WriteString(WarningStyle.Render("  No services found."))
		b.WriteString("\n\n")
		b.WriteString(HelpStyle.Render("r refresh • esc back"))
		return b.String()
	}

	// Docker status
	if core.IsDockerRunning() {
		b.WriteString(SuccessStyle.Render("  ● Docker is running"))
	} else {
		b.WriteString(ErrorStyle.Render("  ○ Docker is not running"))
	}
	b.WriteString("\n\n")

	// Table header
	header := fmt.Sprintf("  %-30s %-15s %-12s %s",
		TableHeaderStyle.Render("SERVICE"),
		TableHeaderStyle.Render("STATUS"),
		TableHeaderStyle.Render("HEALTH"),
		TableHeaderStyle.Render("PORTS"),
	)
	b.WriteString(header)
	b.WriteString("\n")
	b.WriteString("  " + strings.Repeat("─", 75))
	b.WriteString("\n")

	// Service rows
	for _, svc := range m.services {
		statusIcon := "○"
		statusStyle := ErrorStyle
		switch svc.Status {
		case "running":
			statusIcon = "●"
			statusStyle = SuccessStyle
		case "stopped":
			statusIcon = "○"
			statusStyle = ErrorStyle
		case "not found":
			statusIcon = "?"
			statusStyle = WarningStyle
		}

		healthStr := svc.Health
		healthStyle := InfoStyle
		switch svc.Health {
		case "healthy":
			healthStyle = SuccessStyle
		case "unhealthy":
			healthStyle = ErrorStyle
		case "starting":
			healthStyle = WarningStyle
		}

		row := fmt.Sprintf("  %-30s %-15s %-12s %s",
			TableCellStyle.Render(svc.Name),
			statusStyle.Render(statusIcon+" "+svc.Status),
			healthStyle.Render(healthStr),
			TableCellStyle.Render(svc.Port),
		)
		b.WriteString(row)
		b.WriteString("\n")
	}

	// Timestamp
	b.WriteString("\n")
	b.WriteString(InfoStyle.Render(fmt.Sprintf("  Last updated: %s", time.Now().Format("15:04:05"))))
	b.WriteString("\n\n")

	// Help
	b.WriteString(HelpStyle.Render("r refresh • esc back"))

	return b.String()
}
