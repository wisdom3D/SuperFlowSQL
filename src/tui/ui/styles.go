package ui

import "github.com/charmbracelet/lipgloss"

// VSCode-inspired color palette
var (
	// Background colors
	ColorBg        = lipgloss.Color("#1E1E1E") // Editor background
	ColorSidebarBg = lipgloss.Color("#252526") // Sidebar background
	ColorTitleBg   = lipgloss.Color("#007ACC") // Title bar / accent
	ColorStatusBg  = lipgloss.Color("#007ACC") // Status bar
	ColorInputBg   = lipgloss.Color("#3C3C3C") // Input background

	// Text colors
	ColorFg       = lipgloss.Color("#D4D4D4") // Default text
	ColorFgDim    = lipgloss.Color("#808080") // Dimmed text
	ColorFgBright = lipgloss.Color("#FFFFFF") // Bright text
	ColorAccent   = lipgloss.Color("#007ACC") // Accent/links
	ColorSuccess  = lipgloss.Color("#4EC9B0") // Success/green
	ColorWarning  = lipgloss.Color("#DCDCAA") // Warning/yellow
	ColorError    = lipgloss.Color("#F44747") // Error/red
	ColorInfo     = lipgloss.Color("#9CDCFE") // Info/light blue
	ColorKeyword  = lipgloss.Color("#C586C0") // Keyword/purple
	ColorString   = lipgloss.Color("#CE9178") // String/orange
	ColorNumber   = lipgloss.Color("#B5CEA8") // Number/green

	// Styles
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorFgBright).
			Background(ColorTitleBg).
			Padding(0, 2).
			MarginBottom(1)

	SubtitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorAccent).
			MarginBottom(1)

	MenuItemStyle = lipgloss.NewStyle().
			Foreground(ColorFg).
			PaddingLeft(2)

	MenuItemSelectedStyle = lipgloss.NewStyle().
				Foreground(ColorFgBright).
				Background(ColorInputBg).
				Bold(true).
				PaddingLeft(1)

	MenuItemDescStyle = lipgloss.NewStyle().
				Foreground(ColorFgDim).
				PaddingLeft(4)

	StatusBarStyle = lipgloss.NewStyle().
			Foreground(ColorFgBright).
			Background(ColorStatusBg).
			Padding(0, 1)

	HelpStyle = lipgloss.NewStyle().
			Foreground(ColorFgDim).
			MarginTop(1)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(ColorError).
			Bold(true)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(ColorSuccess).
			Bold(true)

	WarningStyle = lipgloss.NewStyle().
			Foreground(ColorWarning)

	InfoStyle = lipgloss.NewStyle().
			Foreground(ColorInfo)

	InputLabelStyle = lipgloss.NewStyle().
			Foreground(ColorKeyword).
			Bold(true)

	InputStyle = lipgloss.NewStyle().
			Foreground(ColorFg).
			Background(ColorInputBg).
			Padding(0, 1)

	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorAccent).
			Padding(1, 2).
			MarginBottom(1)

	LogoStyle = lipgloss.NewStyle().
			Foreground(ColorAccent).
			Bold(true)

	DimBorderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorFgDim).
			Padding(1, 2)

	SectionHeaderStyle = lipgloss.NewStyle().
				Foreground(ColorKeyword).
				Bold(true).
				MarginTop(1).
				MarginBottom(1)

	TableHeaderStyle = lipgloss.NewStyle().
				Foreground(ColorAccent).
				Bold(true).
				PaddingRight(2)

	TableCellStyle = lipgloss.NewStyle().
			Foreground(ColorFg).
			PaddingRight(2)

	BreadcrumbStyle = lipgloss.NewStyle().
			Foreground(ColorFgDim)

	BreadcrumbActiveStyle = lipgloss.NewStyle().
				Foreground(ColorAccent).
				Bold(true)
)

// Logo returns the SuperFlowSQL ASCII art logo
func Logo() string {
	logo := `
 ███████╗██╗   ██╗██████╗ ███████╗██████╗ 
 ██╔════╝██║   ██║██╔══██╗██╔════╝██╔══██╗
 ███████╗██║   ██║██████╔╝█████╗  ██████╔╝
 ╚════██║██║   ██║██╔═══╝ ██╔══╝  ██╔══██╗
 ███████║╚██████╔╝██║     ███████╗██║  ██║
 ╚══════╝ ╚═════╝ ╚═╝     ╚══════╝╚═╝  ╚═╝
 ███████╗██╗      ██████╗ ██╗    ██╗███████╗ ██████╗ ██╗     
 ██╔════╝██║     ██╔═══██╗██║    ██║██╔════╝██╔═══██╗██║     
 █████╗  ██║     ██║   ██║██║ █╗ ██║███████╗██║   ██║██║     
 ██╔══╝  ██║     ██║   ██║██║███╗██║╚════██║██║▄▄ ██║██║     
 ██║     ███████╗╚██████╔╝╚███╔███╔╝███████║╚██████╔╝███████╗
 ╚═╝     ╚══════╝ ╚═════╝  ╚══╝╚══╝ ╚══════╝ ╚══▀▀═╝ ╚══════╝`
	return LogoStyle.Render(logo)
}
