package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"superflowsql-tui/core"
)

// PipelineTypeItem represents a selectable pipeline type
type PipelineTypeItem struct {
	Title       string
	Description string
	Type        core.PipelineType
}

// PipelineSelectModel lets the user choose a pipeline type
type PipelineSelectModel struct {
	items  []PipelineTypeItem
	cursor int
}

// PipelineTypeSelectedMsg is sent when a pipeline type is chosen
type PipelineTypeSelectedMsg struct {
	Type core.PipelineType
}

// NewPipelineSelectModel creates the pipeline type selector
func NewPipelineSelectModel() PipelineSelectModel {
	return PipelineSelectModel{
		items: []PipelineTypeItem{
			{
				Title:       "Generic Pipeline",
				Description: "A basic pipeline with start → process → finish stages",
				Type:        core.PipelineGeneric,
			},
			{
				Title:       "Pandas Ingestion",
				Description: "Generate/transform data with Pandas and load into PostgreSQL",
				Type:        core.PipelinePandasIngestion,
			},
			{
				Title:       "API Ingestion",
				Description: "Fetch data from a REST API and load into PostgreSQL",
				Type:        core.PipelineAPIIngestion,
			},
			{
				Title:       "CSV Ingestion",
				Description: "Read data from CSV files and load into PostgreSQL",
				Type:        core.PipelineCSVIngestion,
			},
		},
	}
}

func (m PipelineSelectModel) Init() tea.Cmd { return nil }

func (m PipelineSelectModel) Update(msg tea.Msg) (PipelineSelectModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.items)-1 {
				m.cursor++
			}
		case "enter":
			return m, func() tea.Msg {
				return PipelineTypeSelectedMsg{Type: m.items[m.cursor].Type}
			}
		case "esc", "q":
			return m, func() tea.Msg { return FormCancelMsg{} }
		}
	}
	return m, nil
}

func (m PipelineSelectModel) View() string {
	var b strings.Builder

	b.WriteString(BreadcrumbStyle.Render("Menu"))
	b.WriteString(BreadcrumbStyle.Render(" > "))
	b.WriteString(BreadcrumbActiveStyle.Render("Add Pipeline"))
	b.WriteString("\n\n")

	b.WriteString(SubtitleStyle.Render("➕ Choose Pipeline Type"))
	b.WriteString("\n\n")

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

	b.WriteString(HelpStyle.Render("↑/↓ navigate • enter select • esc back"))

	return b.String()
}
