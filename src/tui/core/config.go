package core

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// ConfigEntry represents a key-value pair from the .env file
type ConfigEntry struct {
	Key     string
	Value   string
	Comment string // Comment above the entry, if any
}

// LoadConfig reads the .env file from the project directory
func LoadConfig(projectDir string) ([]ConfigEntry, error) {
	envPath := filepath.Join(projectDir, ".env")

	f, err := os.Open(envPath)
	if err != nil {
		return nil, fmt.Errorf("could not open .env file: %w", err)
	}
	defer f.Close()

	var entries []ConfigEntry
	var currentComment string

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines
		if line == "" {
			currentComment = ""
			continue
		}

		// Accumulate comments
		if strings.HasPrefix(line, "#") {
			currentComment = line
			continue
		}

		// Parse key=value
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			entries = append(entries, ConfigEntry{
				Key:     strings.TrimSpace(parts[0]),
				Value:   strings.TrimSpace(parts[1]),
				Comment: currentComment,
			})
			currentComment = ""
		}
	}

	return entries, scanner.Err()
}

// SaveConfig writes config entries back to the .env file
func SaveConfig(projectDir string, entries []ConfigEntry) error {
	envPath := filepath.Join(projectDir, ".env")

	f, err := os.Create(envPath)
	if err != nil {
		return fmt.Errorf("could not write .env file: %w", err)
	}
	defer f.Close()

	// Group entries by their section (based on key prefix)
	writer := bufio.NewWriter(f)
	defer writer.Flush()

	fmt.Fprintln(writer, "# ============================================")
	fmt.Fprintln(writer, "# SuperFlowSQL - Environment Configuration")
	fmt.Fprintln(writer, "# ============================================")
	fmt.Fprintln(writer)

	// Group by prefix
	groups := map[string][]ConfigEntry{}
	groupOrder := []string{}

	for _, entry := range entries {
		prefix := GetGroupPrefix(entry.Key)
		if _, exists := groups[prefix]; !exists {
			groupOrder = append(groupOrder, prefix)
		}
		groups[prefix] = append(groups[prefix], entry)
	}

	for _, prefix := range groupOrder {
		fmt.Fprintf(writer, "# --- %s ---\n", prefix)
		for _, entry := range groups[prefix] {
			fmt.Fprintf(writer, "%s=%s\n", entry.Key, entry.Value)
		}
		fmt.Fprintln(writer)
	}

	return nil
}

// UpdateConfigEntry updates a single entry in the config
func UpdateConfigEntry(projectDir string, key string, value string) error {
	entries, err := LoadConfig(projectDir)
	if err != nil {
		return err
	}

	found := false
	for i, e := range entries {
		if e.Key == key {
			entries[i].Value = value
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("config key '%s' not found", key)
	}

	return SaveConfig(projectDir, entries)
}

// GetConfigKeys returns all available config keys
func GetConfigKeys(entries []ConfigEntry) []string {
	keys := make([]string, len(entries))
	for i, e := range entries {
		keys[i] = e.Key
	}
	sort.Strings(keys)
	return keys
}

func GetGroupPrefix(key string) string {
	parts := strings.SplitN(key, "_", 2)
	if len(parts) > 0 {
		switch strings.ToUpper(parts[0]) {
		case "POSTGRES":
			return "PostgreSQL"
		case "AIRFLOW":
			return "Airflow"
		case "PGADMIN":
			return "PgAdmin"
		case "SUPERSET":
			return "Superset"
		}
	}
	return "General"
}

// FindProjectDir looks for a SuperFlowSQL project directory
// It checks the current directory and parent directories for a docker-compose.yml
func FindProjectDir(startDir string) (string, error) {
	dir := startDir
	for {
		composePath := filepath.Join(dir, "docker-compose.yml")
		envPath := filepath.Join(dir, ".env")
		if _, err := os.Stat(composePath); err == nil {
			if _, err := os.Stat(envPath); err == nil {
				return dir, nil
			}
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return "", fmt.Errorf("no SuperFlowSQL project found (looked for docker-compose.yml + .env)")
}
