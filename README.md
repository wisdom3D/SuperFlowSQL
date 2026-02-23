<div align="center">
  <img src="assets/banner.jpg">
</div>
<div align="center">
  <h3>The CLI that bootstraps and manages your Airflow + PostgreSQL + Superset + PgAdmin stack in minutes.</h3>
</div>
<div align="center">
  <a href="https://pypi.python.org/pypi/superflowsql"><img src="https://img.shields.io/pypi/v/superflowsql.svg" alt="PyPI"></a>
  <a href="https://github.com/superflowsql/superflowsql"><img src="https://img.shields.io/pypi/pyversions/superflowsql.svg" alt="versions"></a>
  <a href="https://github.com/superflowsql/superflowsql/blob/main/LICENSE"><img src="https://img.shields.io/github/license/superflowsql/superflowsql.svg" alt="license"></a>
</div>

---

**Alpha** — This project is under active development. APIs may change.

SuperFlowSQL is a developer-first CLI that scaffolds and manages a complete modern data stack — **Apache Airflow**, **PostgreSQL**, **Apache Superset**, and **PgAdmin** — from a single command.

Instead of cloning boilerplate Docker Compose repos and wiring services together manually, SuperFlowSQL gives you opinionated project scaffolding, pipeline templates, and stack lifecycle management in one tool.

What SuperFlowSQL **does**:

* Scaffolds a fully configured data stack project in seconds (`superflowsql init`)
* Starts/stops/restarts all services with a single command (`superflowsql up`, `down`, `restart`)
* Shows live service status — running, healthy, ports (`superflowsql status`)
* Generates ready-to-use Airflow DAGs from pipeline templates (`superflowsql add-pipeline`)
* Manages project configuration via `.env` (`superflowsql config`)
* Ships an interactive Go-based TUI with a VSCode-inspired interface (`superflowsql tui`)

What SuperFlowSQL **is not**:

* A replacement for Airflow, Superset, or any individual tool
* A hosted SaaS platform
* An enterprise data governance layer

---

SuperFlowSQL is for:

* **Solo devs and indie hackers** who want a local analytics stack without the setup pain
* **Startups** standing up their first data pipelines
* **Consultants** who repeatedly scaffold similar stacks for clients
* **Learners** who want to explore Airflow + Superset without fighting Docker configs

## Installation

```bash
pip install superflowsql
```

**Prerequisites:**

* Python 3.9+
* Docker Desktop (or Docker Engine + Docker Compose)

## Quick Start

### 1. Create a project

```bash
superflowsql init --name my-data-project
```

This scaffolds a complete project with:

```
my-data-project/
├── .env                          # All service credentials and ports
├── docker-compose.yml            # Full stack: Postgres, Airflow, Superset, PgAdmin
├── Dockerfile.airflow            # Custom Airflow image with Python dependencies
├── requirements.txt              # Python packages for your DAGs
├── pgadmin_servers.json          # Pre-configured PgAdmin server connection
├── dags/
│   └── example_dag.py            # Example pipeline to get started
├── postgres/
│   └── init.sql                  # Initial database schema
└── superset/
    └── superset_config.py        # Superset configuration
```

### 2. Start the stack

```bash
cd my-data-project
superflowsql up
```

Your entire data platform is now running:

| Service    | URL                        | Default credentials  |
|------------|----------------------------|----------------------|
| Airflow    | http://localhost:8080       | admin / admin        |
| PgAdmin    | http://localhost:5050       | admin@admin.com / admin |
| Superset   | http://localhost:8088       | admin / admin        |
| PostgreSQL | localhost:5432              | airflow / airflow    |

### 3. Add a pipeline

```bash
superflowsql add-pipeline --name sales_ingestion --type pandas --table sales_data
```

Four pipeline templates are available:

| Type      | Description                                                  |
|-----------|--------------------------------------------------------------|
| `generic` | Basic start → process → finish DAG                           |
| `pandas`  | Generate/transform data with Pandas, load into PostgreSQL    |
| `api`     | Fetch from a REST API, load into PostgreSQL                  |
| `csv`     | Read CSV files, load into PostgreSQL                         |

## Usage

### All commands

```bash
superflowsql init             # Scaffold a new project
superflowsql up               # Start all Docker services
superflowsql down             # Stop and remove containers
superflowsql restart          # Restart the stack
superflowsql status           # View service health and ports
superflowsql add-pipeline     # Create a new Airflow DAG from templates
superflowsql pipelines        # List existing pipelines
superflowsql config           # View all configuration
superflowsql config KEY       # View a single config value
superflowsql config KEY VALUE # Update a config value
superflowsql tui              # Launch the interactive terminal UI
```

### Project scaffolding with options

```bash
superflowsql init \
  --name my-project \
  --pg-user datauser \
  --pg-password secret123 \
  --pg-db warehouse \
  --pg-port 5432 \
  --airflow-user admin \
  --airflow-password admin \
  --airflow-port 8080 \
  --pgadmin-email admin@company.com \
  --pgadmin-password pgadmin123 \
  --pgadmin-port 5050 \
  --superset-user admin \
  --superset-password admin \
  --superset-port 8088
```

### Configuration management

```bash
# View all config
superflowsql config

# Check a single value
superflowsql config POSTGRES_PORT

# Update a value
superflowsql config POSTGRES_PORT 5433
```

### Pipeline creation

```bash
# Interactive (prompts for all options)
superflowsql add-pipeline

# Non-interactive
superflowsql add-pipeline \
  --name user_sync \
  --type api \
  --table users \
  --schedule "0 */6 * * *" \
  --tags '"superflowsql", "users"'
```

### Interactive TUI

SuperFlowSQL ships with a Go-based terminal UI built with [Bubbletea](https://github.com/charmbracelet/bubbletea), offering a VSCode-inspired interface with:

* Keyboard-navigable menu
* Multi-field forms with tab navigation
* Live service status with health indicators
* docker compose progress with spinners
* Inline configuration editor

```bash
superflowsql tui
```

## Python API

SuperFlowSQL can also be used programmatically:

```python
from superflowsql.project import ProjectConfig, init_project
from superflowsql.docker import start_stack, stop_stack, get_stack_status
from superflowsql.pipeline import PipelineConfig, PipelineType, add_pipeline
from superflowsql.config import load_config, update_config

# Scaffold a project
config = ProjectConfig(project_name="analytics", postgres_user="analyst")
project_dir = init_project(config, target_dir=".")

# Start the stack
start_stack(project_dir)

# Check status
for svc in get_stack_status(project_dir):
    print(f"{svc.name}: {svc.status} ({svc.health})")

# Add a pipeline
pipeline = PipelineConfig(
    name="daily_report",
    pipeline_type=PipelineType.PANDAS,
    table_name="daily_metrics",
    schedule='"0 0 * * *"',
)
add_pipeline(pipeline, project_dir)

# View config
entries = load_config(project_dir)
for entry in entries:
    print(f"{entry.group}: {entry.key} = {entry.value}")

# Update config
update_config(project_dir, "POSTGRES_PORT", "5433")

# Stop the stack
stop_stack(project_dir)
```

## Architecture

SuperFlowSQL is a monorepo with three components:

```
superflowsql-proj/
├── superflowsql/            # Python package (PyPI)
│   ├── cli.py               # Click CLI entry point
│   ├── project.py           # Project scaffolding
│   ├── docker.py            # Docker Compose operations
│   ├── pipeline.py          # Pipeline DAG generation
│   ├── config.py            # .env configuration management
│   └── templates/           # Jinja2 templates for all generated files
├── src/
│   ├── tui/                 # Go TUI (Bubbletea)
│   │   ├── main.go          # App model and view routing
│   │   ├── ui/              # UI components (menu, forms, status, styles)
│   │   └── core/            # Business logic (project, docker, pipeline, config)
│   ├── docker-compose.yml   # Reference stack definition
│   └── dags/                # Reference DAG examples
├── superflowsqllanding/     # Next.js landing page
└── pyproject.toml            # Package configuration
```

### Component responsibilities

| Component         | Language | Purpose                                                  |
|-------------------|----------|----------------------------------------------------------|
| `superflowsql/`   | Python   | PyPI package — CLI commands, project scaffolding, Docker management, pipeline generation |
| `src/tui/`         | Go       | Interactive terminal UI — VSCode-inspired Bubbletea interface |
| `superflowsqllanding/` | TypeScript | Landing page — Next.js marketing site              |

### Stack services

| Service                | Image                     | Default Port |
|------------------------|---------------------------|-------------|
| PostgreSQL 15          | `postgres:15`             | 5432        |
| Apache Airflow 2.9     | Custom (Dockerfile)       | 8080        |
| Airflow Scheduler      | Custom (Dockerfile)       | —           |
| Apache Superset 3.0    | `apache/superset:3.0.2`   | 8088        |
| PgAdmin 4              | `dpage/pgadmin4:latest`   | 5050        |

## Development

```bash
# Clone the repo
git clone https://github.com/superflowsql/superflowsql.git
cd superflowsql

# Install Python package in editable mode
pip install -e ".[dev]"

# Build the Go TUI
cd src/tui
go build -o superflowsql-tui .

# Run the CLI
superflowsql --help

# Run the TUI
superflowsql tui
```
