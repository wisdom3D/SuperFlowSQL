"""Pipeline creation - generate Airflow DAG files from templates."""

from __future__ import annotations

import re
from dataclasses import dataclass
from enum import Enum
from pathlib import Path

from jinja2 import Environment, PackageLoader, select_autoescape


class PipelineType(str, Enum):
    """Available pipeline template types."""

    GENERIC = "generic"
    PANDAS = "pandas"
    API = "api"
    CSV = "csv"

    @property
    def description(self) -> str:
        descriptions = {
            "generic": "A basic pipeline with start → process → finish stages",
            "pandas": "Generate/transform data with Pandas and load into PostgreSQL",
            "api": "Fetch data from a REST API and load into PostgreSQL",
            "csv": "Read data from CSV files and load into PostgreSQL",
        }
        return descriptions[self.value]

    @property
    def template_name(self) -> str:
        return f"pipeline_{self.value}.py.j2"


@dataclass
class PipelineConfig:
    """Configuration for a new pipeline DAG."""

    name: str = "my_pipeline"
    description: str = "A new data pipeline"
    schedule: str = "None"
    pipeline_type: PipelineType = PipelineType.GENERIC
    table_name: str = "my_table"
    tags: str = '"superflowsql"'

    @property
    def safe_name(self) -> str:
        """Sanitize the pipeline name for use as a Python identifier / filename."""
        name = self.name.lower()
        name = re.sub(r"[^a-z0-9_]", "_", name)
        name = re.sub(r"_+", "_", name)
        return name.strip("_")


def add_pipeline(config: PipelineConfig, project_dir: str | Path) -> Path:
    """Create a new pipeline DAG file in the project's dags directory.

    Args:
        config: Pipeline configuration (name, type, schedule, etc.)
        project_dir: Path to the SuperFlowSQL project directory.

    Returns:
        Path to the created DAG file.

    Raises:
        FileExistsError: If a pipeline with the same name already exists.
        FileNotFoundError: If the project directory doesn't exist.
    """
    project_dir = Path(project_dir).resolve()
    dags_dir = project_dir / "dags"

    if not project_dir.exists():
        raise FileNotFoundError(f"Project directory not found: {project_dir}")

    dags_dir.mkdir(parents=True, exist_ok=True)

    filename = f"{config.safe_name}.py"
    filepath = dags_dir / filename

    if filepath.exists():
        raise FileExistsError(f"Pipeline '{config.name}' already exists: {filepath}")

    # Render template
    env = Environment(
        loader=PackageLoader("superflowsql", "templates"),
        autoescape=select_autoescape([]),
        keep_trailing_newline=True,
        trim_blocks=True,
        lstrip_blocks=True,
    )
    tmpl = env.get_template(config.pipeline_type.template_name)
    rendered = tmpl.render(
        name=config.safe_name,
        description=config.description,
        schedule=config.schedule,
        table_name=config.table_name,
        tags=config.tags,
    )
    filepath.write_text(rendered, encoding="utf-8")

    return filepath


def list_pipelines(project_dir: str | Path) -> list[str]:
    """List all existing pipeline DAG files in the project.

    Args:
        project_dir: Path to the SuperFlowSQL project directory.

    Returns:
        List of pipeline filenames.
    """
    dags_dir = Path(project_dir).resolve() / "dags"
    if not dags_dir.exists():
        return []
    return sorted(f.name for f in dags_dir.glob("*.py") if not f.name.startswith("__"))
