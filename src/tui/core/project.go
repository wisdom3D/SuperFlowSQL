package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// ProjectConfig holds the configuration for a new project
type ProjectConfig struct {
	ProjectName       string
	PostgresUser      string
	PostgresPassword  string
	PostgresDB        string
	PostgresPort      string
	PostgresHost      string
	AirflowUser       string
	AirflowPassword   string
	AirflowEmail      string
	AirflowPort       string
	AirflowFernetKey  string
	PgAdminEmail      string
	PgAdminPassword   string
	PgAdminPort       string
	SupersetUser      string
	SupersetPassword  string
	SupersetEmail     string
	SupersetPort      string
	SupersetSecretKey string
}

// DefaultProjectConfig returns a config with sensible defaults
func DefaultProjectConfig() ProjectConfig {
	return ProjectConfig{
		ProjectName:       "my-superflowsql-project",
		PostgresUser:      "airflow",
		PostgresPassword:  "airflow",
		PostgresDB:        "airflow",
		PostgresPort:      "5432",
		PostgresHost:      "postgres",
		AirflowUser:       "admin",
		AirflowPassword:   "admin",
		AirflowEmail:      "admin@example.com",
		AirflowPort:       "8080",
		AirflowFernetKey:  "ZmDfcTF7_60GrrY167zsiPd67pEvs0aGOv2oasOM1Pg=",
		PgAdminEmail:      "admin@admin.com",
		PgAdminPassword:   "admin",
		PgAdminPort:       "5050",
		SupersetUser:      "admin",
		SupersetPassword:  "admin",
		SupersetEmail:     "admin@superset.com",
		SupersetPort:      "8088",
		SupersetSecretKey: "superflowsql_secret_key_change_me",
	}
}

// InitProject scaffolds a new SuperFlowSQL project
func InitProject(config ProjectConfig, targetDir string) error {
	projectDir := filepath.Join(targetDir, config.ProjectName)

	// Create directory structure
	dirs := []string{
		projectDir,
		filepath.Join(projectDir, "dags"),
		filepath.Join(projectDir, "postgres"),
		filepath.Join(projectDir, "superset"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Generate files from templates
	files := map[string]string{
		filepath.Join(projectDir, ".env"):                           envTemplate,
		filepath.Join(projectDir, "docker-compose.yml"):             dockerComposeTemplate,
		filepath.Join(projectDir, "Dockerfile.airflow"):             dockerfileAirflowTemplate,
		filepath.Join(projectDir, "requirements.txt"):               requirementsTemplate,
		filepath.Join(projectDir, "postgres", "init.sql"):           initSQLTemplate,
		filepath.Join(projectDir, "pgadmin_servers.json"):           pgAdminServersTemplate,
		filepath.Join(projectDir, "superset", "superset_config.py"): supersetConfigTemplate,
		filepath.Join(projectDir, "dags", "example_dag.py"):         exampleDAGTemplate,
	}

	for filePath, tmplContent := range files {
		if err := renderTemplate(filePath, tmplContent, config); err != nil {
			return fmt.Errorf("failed to generate %s: %w", filepath.Base(filePath), err)
		}
	}

	return nil
}

func renderTemplate(filePath, tmplContent string, data interface{}) error {
	funcMap := template.FuncMap{
		"lower": strings.ToLower,
	}
	tmpl, err := template.New(filepath.Base(filePath)).Funcs(funcMap).Parse(tmplContent)
	if err != nil {
		return fmt.Errorf("template parse error: %w", err)
	}

	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("file create error: %w", err)
	}
	defer f.Close()

	return tmpl.Execute(f, data)
}
