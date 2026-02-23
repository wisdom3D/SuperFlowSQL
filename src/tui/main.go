package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"superflowsql-tui/core"
	"superflowsql-tui/ui"
)

// appState tracks which view is currently active
type appState int

const (
	stateMenu appState = iota
	stateInitForm
	stateStartStack
	stateStopStack
	stateStackStatus
	statePipelineSelect
	statePipelineForm
	stateConfig
)

// model is the top-level Bubbletea model
type model struct {
	state            appState
	menuModel        ui.MenuModel
	formModel        ui.FormModel
	statusModel      ui.StatusModel
	actionModel      ui.ActionModel
	configModel      ui.ConfigViewModel
	pipelineSelect   ui.PipelineSelectModel
	selectedPipeline core.PipelineType
	width            int
	height           int
	projectDir       string
}

func initialModel() model {
	// Try to find existing project directory
	cwd, _ := os.Getwd()
	projectDir, _ := core.FindProjectDir(cwd)
	if projectDir == "" {
		projectDir = cwd
	}

	return model{
		state:      stateMenu,
		menuModel:  ui.NewMenuModel(),
		projectDir: projectDir,
	}
}

func (m model) Init() tea.Cmd {
	return tea.EnterAltScreen
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	switch m.state {
	case stateMenu:
		return m.updateMenu(msg)
	case stateInitForm:
		return m.updateInitForm(msg)
	case stateStartStack:
		return m.updateAction(msg)
	case stateStopStack:
		return m.updateAction(msg)
	case stateStackStatus:
		return m.updateStatus(msg)
	case statePipelineSelect:
		return m.updatePipelineSelect(msg)
	case statePipelineForm:
		return m.updatePipelineForm(msg)
	case stateConfig:
		return m.updateConfig(msg)
	}

	return m, nil
}

// updateMenu handles the main menu
func (m model) updateMenu(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case ui.MenuSelectionMsg:
		switch msg.View {
		case ui.ViewInitProject:
			m.state = stateInitForm
			defaults := core.DefaultProjectConfig()
			m.formModel = ui.NewFormModel("⚡ Init Project", []ui.FormField{
				{Label: "Project Name", Placeholder: defaults.ProjectName},
				{Label: "PostgreSQL User", Placeholder: defaults.PostgresUser},
				{Label: "PostgreSQL Password", Placeholder: defaults.PostgresPassword, IsPassword: true},
				{Label: "PostgreSQL Database", Placeholder: defaults.PostgresDB},
				{Label: "PostgreSQL Port", Placeholder: defaults.PostgresPort},
				{Label: "Airflow User", Placeholder: defaults.AirflowUser},
				{Label: "Airflow Password", Placeholder: defaults.AirflowPassword, IsPassword: true},
				{Label: "Airflow Port", Placeholder: defaults.AirflowPort},
				{Label: "PgAdmin Email", Placeholder: defaults.PgAdminEmail},
				{Label: "PgAdmin Password", Placeholder: defaults.PgAdminPassword, IsPassword: true},
				{Label: "PgAdmin Port", Placeholder: defaults.PgAdminPort},
				{Label: "Superset User", Placeholder: defaults.SupersetUser},
				{Label: "Superset Password", Placeholder: defaults.SupersetPassword, IsPassword: true},
				{Label: "Superset Port", Placeholder: defaults.SupersetPort},
			})
			return m, m.formModel.Init()

		case ui.ViewStartStack:
			m.state = stateStartStack
			m.actionModel = ui.NewActionModel("▶  Start Stack", core.StackStart, m.projectDir)
			return m, m.actionModel.Init()

		case ui.ViewStopStack:
			m.state = stateStopStack
			m.actionModel = ui.NewActionModel("⏹  Stop Stack", core.StackStop, m.projectDir)
			return m, m.actionModel.Init()

		case ui.ViewStackStatus:
			m.state = stateStackStatus
			m.statusModel = ui.NewStatusModel(m.projectDir)
			return m, m.statusModel.Init()

		case ui.ViewAddPipeline:
			m.state = statePipelineSelect
			m.pipelineSelect = ui.NewPipelineSelectModel()
			return m, nil

		case ui.ViewConfig:
			m.state = stateConfig
			m.configModel = ui.NewConfigViewModel(m.projectDir)
			return m, nil
		}

	default:
		var cmd tea.Cmd
		m.menuModel, cmd = m.menuModel.Update(msg)
		return m, cmd
	}

	return m, nil
}

// updateInitForm handles the init project form
func (m model) updateInitForm(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case ui.FormSubmitMsg:
		// Build project config from form values
		config := core.ProjectConfig{
			ProjectName:       msg.Values["Project Name"],
			PostgresUser:      msg.Values["PostgreSQL User"],
			PostgresPassword:  msg.Values["PostgreSQL Password"],
			PostgresDB:        msg.Values["PostgreSQL Database"],
			PostgresPort:      msg.Values["PostgreSQL Port"],
			PostgresHost:      "postgres",
			AirflowUser:       msg.Values["Airflow User"],
			AirflowPassword:   msg.Values["Airflow Password"],
			AirflowEmail:      "admin@example.com",
			AirflowPort:       msg.Values["Airflow Port"],
			AirflowFernetKey:  "ZmDfcTF7_60GrrY167zsiPd67pEvs0aGOv2oasOM1Pg=",
			PgAdminEmail:      msg.Values["PgAdmin Email"],
			PgAdminPassword:   msg.Values["PgAdmin Password"],
			PgAdminPort:       msg.Values["PgAdmin Port"],
			SupersetUser:      msg.Values["Superset User"],
			SupersetPassword:  msg.Values["Superset Password"],
			SupersetEmail:     "admin@superset.com",
			SupersetPort:      msg.Values["Superset Port"],
			SupersetSecretKey: "superflowsql_secret_key_change_me",
		}

		cwd, _ := os.Getwd()
		err := core.InitProject(config, cwd)
		if err != nil {
			// Show error but stay on form... for now, go back to menu
			m.state = stateMenu
			return m, nil
		}

		// Update project dir to the new project
		m.projectDir = fmt.Sprintf("%s/%s", cwd, config.ProjectName)
		m.state = stateMenu
		return m, nil

	case ui.FormCancelMsg:
		m.state = stateMenu
		return m, nil

	default:
		var cmd tea.Cmd
		m.formModel, cmd = m.formModel.Update(msg)
		return m, cmd
	}
}

// updateAction handles start/stop stack actions
func (m model) updateAction(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case ui.FormCancelMsg:
		m.state = stateMenu
		return m, nil
	default:
		var cmd tea.Cmd
		m.actionModel, cmd = m.actionModel.Update(msg)
		return m, cmd
	}
}

// updateStatus handles the stack status view
func (m model) updateStatus(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case ui.FormCancelMsg:
		m.state = stateMenu
		return m, nil
	default:
		var cmd tea.Cmd
		m.statusModel, cmd = m.statusModel.Update(msg)
		return m, cmd
	}
}

// updatePipelineSelect handles pipeline type selection
func (m model) updatePipelineSelect(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case ui.PipelineTypeSelectedMsg:
		m.selectedPipeline = msg.Type
		m.state = statePipelineForm
		defaults := core.DefaultPipelineConfig()
		m.formModel = ui.NewFormModel("➕ Add Pipeline", []ui.FormField{
			{Label: "Pipeline Name", Placeholder: defaults.Name},
			{Label: "Description", Placeholder: defaults.Description},
			{Label: "Table Name", Placeholder: defaults.TableName},
			{Label: "Schedule", Placeholder: "None"},
			{Label: "Tags", Placeholder: "\"superflowsql\""},
		})
		return m, m.formModel.Init()

	case ui.FormCancelMsg:
		m.state = stateMenu
		return m, nil

	default:
		var cmd tea.Cmd
		m.pipelineSelect, cmd = m.pipelineSelect.Update(msg)
		return m, cmd
	}
}

// updatePipelineForm handles the pipeline detail form
func (m model) updatePipelineForm(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case ui.FormSubmitMsg:
		config := core.PipelineConfig{
			Name:         msg.Values["Pipeline Name"],
			Description:  msg.Values["Description"],
			TableName:    msg.Values["Table Name"],
			Schedule:     msg.Values["Schedule"],
			Tags:         msg.Values["Tags"],
			PipelineType: m.selectedPipeline,
		}

		err := core.AddPipeline(config, m.projectDir)
		if err != nil {
			// For now go back to menu
			m.state = stateMenu
			return m, nil
		}

		m.state = stateMenu
		return m, nil

	case ui.FormCancelMsg:
		m.state = statePipelineSelect
		return m, nil

	default:
		var cmd tea.Cmd
		m.formModel, cmd = m.formModel.Update(msg)
		return m, cmd
	}
}

// updateConfig handles the configuration view
func (m model) updateConfig(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case ui.FormCancelMsg:
		m.state = stateMenu
		return m, nil
	default:
		var cmd tea.Cmd
		m.configModel, cmd = m.configModel.Update(msg)
		return m, cmd
	}
}

func (m model) View() string {
	switch m.state {
	case stateMenu:
		return m.menuModel.View()
	case stateInitForm:
		return m.formModel.View()
	case stateStartStack:
		return m.actionModel.View()
	case stateStopStack:
		return m.actionModel.View()
	case stateStackStatus:
		return m.statusModel.View()
	case statePipelineSelect:
		return m.pipelineSelect.View()
	case statePipelineForm:
		return m.formModel.View()
	case stateConfig:
		return m.configModel.View()
	}
	return ""
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running SuperFlowSQL TUI: %v\n", err)
		os.Exit(1)
	}
}
