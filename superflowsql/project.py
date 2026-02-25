"""Project scaffolding - generates a full SuperFlowSQL project from templates."""

from __future__ import annotations

import os
from dataclasses import dataclass, field
from pathlib import Path

from jinja2 import Environment, PackageLoader, select_autoescape


@dataclass
class ProjectConfig:
    """Configuration for a new SuperFlowSQL project."""

    project_name: str = "my-superflowsql-project"
    postgres_user: str = "airflow"
    postgres_password: str = "airflow"
    postgres_db: str = "airflow"
    postgres_port: str = "5432"
    postgres_host: str = "postgres"
    airflow_user: str = "admin"
    airflow_password: str = "admin"
    airflow_email: str = "admin@example.com"
    airflow_port: str = "8080"
    airflow_fernet_key: str = "ZmDfcTF7_60GrrY167zsiPd67pEvs0aGOv2oasOM1Pg="
    pgadmin_email: str = "admin@admin.com"
    pgadmin_password: str = "admin"
    pgadmin_port: str = "5050"
    superset_user: str = "admin"
    superset_password: str = "admin"
    superset_email: str = "admin@superset.com"
    superset_port: str = "8088"
    superset_secret_key: str = "superflowsql_secret_key_change_me"

    def as_dict(self) -> dict:
        """Convert to a dict for template rendering."""
        return {k: v for k, v in self.__dict__.items()}


# Template file mapping: (template_name, output_relative_path)
_TEMPLATE_FILES: list[tuple[str, str]] = [
    ("env.j2", ".env"),
    ("docker-compose.yml.j2", "docker-compose.yml"),
    ("Dockerfile.airflow.j2", "Dockerfile.airflow"),
    ("requirements.txt.j2", "requirements.txt"),
    ("init.sql.j2", "postgres/init.sql"),
    ("pgadmin_servers.json.j2", "pgadmin_servers.json"),
    ("superset_config.py.j2", "superset/superset_config.py"),
    ("example_dag.py.j2", "dags/example_dag.py"),
]


def init_project(config: ProjectConfig, target_dir: str | Path = ".") -> Path:
    """Scaffold a new SuperFlowSQL project.

    Args:
        config: Project configuration with all service settings.
        target_dir: Parent directory in which to create the project folder.

    Returns:
        Path to the created project directory.

    Raises:
        FileExistsError: If the project directory already exists.
        OSError: If directory or file creation fails.
    """
    target_dir = Path(target_dir).resolve()
    project_dir = target_dir / config.project_name

    if project_dir.exists():
        raise FileExistsError(f"Project directory already exists: {project_dir}")

    # Create directory structure
    dirs = [
        project_dir,
        project_dir / "dags",
        project_dir / "postgres",
        project_dir / "superset",
    ]
    for d in dirs:
        d.mkdir(parents=True, exist_ok=True)

    # Render templates
    env = _get_jinja_env()
    context = config.as_dict()

    for template_name, output_path in _TEMPLATE_FILES:
        tmpl = env.get_template(template_name)
        rendered = tmpl.render(**context)
        out_file = project_dir / output_path
        out_file.parent.mkdir(parents=True, exist_ok=True)
        out_file.write_text(rendered, encoding="utf-8")

    return project_dir


def _get_jinja_env() -> Environment:
    """Create a Jinja2 environment loading from the templates package."""
    return Environment(
        loader=PackageLoader("superflowsql", "templates"),
        autoescape=select_autoescape([]),
        keep_trailing_newline=True,
        trim_blocks=True,
        lstrip_blocks=True,
    )
