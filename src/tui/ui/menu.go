package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// View represents different views/screens in the app
type View int

const (
	ViewMenu View = iota
	ViewInitProject
	ViewStartStack
	ViewStopStack
	ViewStackStatus
	ViewAddPipeline
	ViewConfig
	ViewLogs
)

// MenuItem represents a menu entry
type MenuItem struct {
	Title       string
	Description string
	View        View
}

// MenuModel is the main menu bubbletea model
type MenuModel struct {
	items    []MenuItem
	cursor   int
	width    int
	height   int
	quitting bool
}

// NewMenuModel creates a new main menu
func NewMenuModel() MenuModel {
	items := []MenuItem{
		{Title: "⚡ Init Project", Description: "Scaffold a new SuperFlowSQL project with all config files", View: ViewInitProject},
		{Title: "▶  Start Stack", Description: "Start all Docker services (Postgres, Airflow, Superset, PgAdmin)", View: ViewStartStack},
		{Title: "⏹  Stop Stack", Description: "Stop and remove all running containers", View: ViewStopStack},
		{Title: "📊 Stack Status", Description: "View the health and status of all services", View: ViewStackStatus},
		{Title: "➕ Add Pipeline", Description: "Create a new Airflow DAG from templates", View: ViewAddPipeline},
		{Title: "⚙  Configuration", Description: "View and edit project configuration (.env)", View: ViewConfig},
	}

	return MenuModel{
		items:  items,
		cursor: 0,
	}
}

func (m MenuModel) Init() tea.Cmd {
	return nil
}

// MenuSelectionMsg is sent when a menu item is selected
type MenuSelectionMsg struct {
	View View
}

func (m MenuModel) Update(msg tea.Msg) (MenuModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.items)-1 {
				m.cursor++
			}
		case "enter", " ":
			return m, func() tea.Msg {
				return MenuSelectionMsg{View: m.items[m.cursor].View}
			}
		}
	}
	return m, nil
}

func (m MenuModel) View() string {
	if m.quitting {
		return ""
	}

	var b strings.Builder

	// Logo
	b.WriteString(Logo())
	b.WriteString("\n\n")

	// Title
	b.WriteString(TitleStyle.Render(" SuperFlowSQL "))
	b.WriteString("  ")
	b.WriteString(InfoStyle.Render("Data Orchestration CLI"))
	b.WriteString("\n\n")

	// Menu items
	for i, item := range m.items {
		if i == m.cursor {
			pointer := lipgloss.NewStyle().
				Foreground(ColorAccent).
				Bold(true).
				Render("❯ ")
			title := MenuItemSelectedStyle.Render(item.Title)
			b.WriteString(fmt.Sprintf("%s%s\n", pointer, title))
			b.WriteString(MenuItemDescStyle.Render(item.Description))
		} else {
			b.WriteString(MenuItemStyle.Render(fmt.Sprintf("  %s", item.Title)))
			b.WriteString("\n")
			b.WriteString(MenuItemDescStyle.Render(item.Description))
		}
		b.WriteString("\n\n")
	}

	// Help
	b.WriteString(HelpStyle.Render("↑/↓ navigate • enter select • q quit"))

	return b.String()
}
