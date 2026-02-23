"""SuperFlowSQL CLI - command-line interface for managing data stacks."""

from __future__ import annotations

import os
import sys
from pathlib import Path

import click
from rich.console import Console
from rich.panel import Panel
from rich.table import Table
from rich.text import Text

from superflowsql import __version__
from superflowsql.config import find_project_dir, load_config, update_config
from superflowsql.docker import (
    get_service_urls,
    get_stack_status,
    is_docker_running,
    restart_stack,
    start_stack,
    stop_stack,
)
from superflowsql.pipeline import PipelineConfig, PipelineType, add_pipeline, list_pipelines
from superflowsql.project import ProjectConfig, init_project

console = Console()

LOGO = r"""
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
 ╚═╝     ╚══════╝ ╚═════╝  ╚══╝╚══╝ ╚══════╝ ╚══▀▀═╝ ╚══════╝
"""


def _find_project() -> Path:
    """Find the current project directory or exit with an error."""
    project = find_project_dir()
    if project is None:
        console.print(
            "[red]✗ No SuperFlowSQL project found.[/red]"
            "\n  Run [cyan]superflowsql init[/cyan] to create one."
        )
        sys.exit(1)
    return project


@click.group()
@click.version_option(version=__version__, prog_name="superflowsql")
def main():
    """SuperFlowSQL - Data orchestration CLI.

    Bootstrap and manage your Airflow + PostgreSQL + Superset + PgAdmin stack.
    """
    pass


# ──────────────────────────────────────────────
# superflowsql init
# ──────────────────────────────────────────────
@main.command()
@click.option("--name", prompt="Project name", default="my-superflowsql-project", help="Project name.")
@click.option("--pg-user", default="airflow", help="PostgreSQL user.")
@click.option("--pg-password", default="airflow", help="PostgreSQL password.")
@click.option("--pg-db", default="airflow", help="PostgreSQL database name.")
@click.option("--pg-port", default="5432", help="PostgreSQL port.")
@click.option("--airflow-user", default="admin", help="Airflow admin user.")
@click.option("--airflow-password", default="admin", help="Airflow admin password.")
@click.option("--airflow-port", default="8080", help="Airflow webserver port.")
@click.option("--pgadmin-email", default="admin@admin.com", help="PgAdmin admin email.")
@click.option("--pgadmin-password", default="admin", help="PgAdmin admin password.")
@click.option("--pgadmin-port", default="5050", help="PgAdmin port.")
@click.option("--superset-user", default="admin", help="Superset admin user.")
@click.option("--superset-password", default="admin", help="Superset admin password.")
@click.option("--superset-port", default="8088", help="Superset port.")
def init(
    name,
    pg_user,
    pg_password,
    pg_db,
    pg_port,
    airflow_user,
    airflow_password,
    airflow_port,
    pgadmin_email,
    pgadmin_password,
    pgadmin_port,
    superset_user,
    superset_password,
    superset_port,
):
    """Scaffold a new SuperFlowSQL project with all configuration files."""
    console.print(Text(LOGO, style="bold blue"))

    config = ProjectConfig(
        project_name=name,
        postgres_user=pg_user,
        postgres_password=pg_password,
        postgres_db=pg_db,
        postgres_port=pg_port,
        airflow_user=airflow_user,
        airflow_password=airflow_password,
        airflow_port=airflow_port,
        pgadmin_email=pgadmin_email,
        pgadmin_password=pgadmin_password,
        pgadmin_port=pgadmin_port,
        superset_user=superset_user,
        superset_password=superset_password,
        superset_port=superset_port,
    )

    try:
        project_dir = init_project(config, os.getcwd())
        console.print(f"\n[green]✓ Project created at:[/green] {project_dir}")
        console.print("\n[bold]Next steps:[/bold]")
        console.print(f"  [cyan]cd {name}[/cyan]")
        console.print("  [cyan]superflowsql up[/cyan]")
    except FileExistsError as e:
        console.print(f"\n[red]✗ {e}[/red]")
        sys.exit(1)
    except Exception as e:
        console.print(f"\n[red]✗ Failed to create project: {e}[/red]")
        sys.exit(1)


# ──────────────────────────────────────────────
# superflowsql up
# ──────────────────────────────────────────────
@main.command()
def up():
    """Start all Docker services (Postgres, Airflow, Superset, PgAdmin)."""
    project = _find_project()

    if not is_docker_running():
        console.print("[red]✗ Docker is not running. Please start Docker first.[/red]")
        sys.exit(1)

    console.print("[blue]▶ Starting stack...[/blue]")
    with console.status("[bold blue]Running docker compose up..."):
        try:
            output = start_stack(project)
            console.print(f"[green]✓ Stack started successfully.[/green]")

            # Show service URLs
            urls = get_service_urls()
            table = Table(title="Service URLs", show_header=True, header_style="bold cyan")
            table.add_column("Service", style="bold")
            table.add_column("URL", style="cyan")
            for name, url in urls.items():
                table.add_row(name, url)
            console.print(table)

            if output.strip():
                console.print(Panel(output[-500:], title="Output", border_style="dim"))
        except RuntimeError as e:
            console.print(f"[red]✗ {e}[/red]")
            sys.exit(1)


# ──────────────────────────────────────────────
# superflowsql down
# ──────────────────────────────────────────────
@main.command()
def down():
    """Stop and remove all running containers."""
    project = _find_project()

    console.print("[yellow]⏹ Stopping stack...[/yellow]")
    with console.status("[bold yellow]Running docker compose down..."):
        try:
            output = stop_stack(project)
            console.print("[green]✓ Stack stopped.[/green]")
            if output.strip():
                console.print(Panel(output[-500:], title="Output", border_style="dim"))
        except RuntimeError as e:
            console.print(f"[red]✗ {e}[/red]")
            sys.exit(1)


# ──────────────────────────────────────────────
# superflowsql restart
# ──────────────────────────────────────────────
@main.command()
def restart():
    """Restart all services (stop then start)."""
    project = _find_project()

    if not is_docker_running():
        console.print("[red]✗ Docker is not running.[/red]")
        sys.exit(1)

    console.print("[blue]🔄 Restarting stack...[/blue]")
    with console.status("[bold blue]Restarting..."):
        try:
            output = restart_stack(project)
            console.print("[green]✓ Stack restarted.[/green]")

            urls = get_service_urls()
            table = Table(title="Service URLs", show_header=True, header_style="bold cyan")
            table.add_column("Service", style="bold")
            table.add_column("URL", style="cyan")
            for name, url in urls.items():
                table.add_row(name, url)
            console.print(table)
        except RuntimeError as e:
            console.print(f"[red]✗ {e}[/red]")
            sys.exit(1)


# ──────────────────────────────────────────────
# superflowsql status
# ──────────────────────────────────────────────
@main.command()
def status():
    """View the health and status of all services."""
    project = _find_project()

    # Docker check
    docker_ok = is_docker_running()
    if docker_ok:
        console.print("[green]● Docker is running[/green]")
    else:
        console.print("[red]○ Docker is not running[/red]")
        sys.exit(1)

    console.print()

    # Service table
    services = get_stack_status(project)
    table = Table(title="Stack Status", show_header=True, header_style="bold cyan")
    table.add_column("Service", style="bold", min_width=25)
    table.add_column("Status", min_width=12)
    table.add_column("Health", min_width=10)
    table.add_column("Ports")

    for svc in services:
        status_str = f"{svc.icon} {svc.status}"
        status_style = "green" if svc.is_running else "red" if svc.status == "stopped" else "yellow"
        health_style = (
            "green" if svc.health == "healthy"
            else "red" if svc.health == "unhealthy"
            else "yellow" if svc.health == "starting"
            else "dim"
        )
        table.add_row(
            svc.name,
            f"[{status_style}]{status_str}[/{status_style}]",
            f"[{health_style}]{svc.health}[/{health_style}]",
            svc.ports,
        )

    console.print(table)


# ──────────────────────────────────────────────
# superflowsql add-pipeline
# ──────────────────────────────────────────────
@main.command("add-pipeline")
@click.option("--name", prompt="Pipeline name", default="my_pipeline", help="Pipeline DAG name.")
@click.option("--description", default="A new data pipeline", help="Pipeline description.")
@click.option(
    "--type",
    "pipeline_type",
    type=click.Choice(["generic", "pandas", "api", "csv"], case_sensitive=False),
    prompt="Pipeline type",
    default="generic",
    help="Template type for the pipeline.",
)
@click.option("--table", "table_name", default="my_table", help="Target PostgreSQL table name.")
@click.option("--schedule", default="None", help='Cron schedule (e.g. "0 0 * * *") or None.')
@click.option("--tags", default='"superflowsql"', help="Pipeline tags (comma-separated, quoted).")
def add_pipeline_cmd(name, description, pipeline_type, table_name, schedule, tags):
    """Create a new Airflow DAG from templates."""
    project = _find_project()

    config = PipelineConfig(
        name=name,
        description=description,
        schedule=schedule,
        pipeline_type=PipelineType(pipeline_type),
        table_name=table_name,
        tags=tags,
    )

    try:
        filepath = add_pipeline(config, project)
        console.print(f"\n[green]✓ Pipeline created:[/green] {filepath}")
        console.print(f"  Type: [cyan]{pipeline_type}[/cyan]")
        console.print(f"  Schedule: [cyan]{schedule}[/cyan]")
    except FileExistsError as e:
        console.print(f"\n[red]✗ {e}[/red]")
        sys.exit(1)
    except Exception as e:
        console.print(f"\n[red]✗ Failed to create pipeline: {e}[/red]")
        sys.exit(1)


# ──────────────────────────────────────────────
# superflowsql pipelines
# ──────────────────────────────────────────────
@main.command()
def pipelines():
    """List all pipeline DAG files in the project."""
    project = _find_project()
    pipeline_list = list_pipelines(project)

    if not pipeline_list:
        console.print(
            "[yellow]No pipelines found.[/yellow]"
            "\n  Run [cyan]superflowsql add-pipeline[/cyan] to create one."
        )
        return

    table = Table(title="Pipelines", show_header=True, header_style="bold cyan")
    table.add_column("#", style="dim")
    table.add_column("File", style="bold")

    for i, name in enumerate(pipeline_list, 1):
        table.add_row(str(i), name)

    console.print(table)


# ──────────────────────────────────────────────
# superflowsql config
# ──────────────────────────────────────────────
@main.command()
@click.argument("key", required=False)
@click.argument("value", required=False)
def config(key, value):
    """View or update project configuration (.env).

    \b
    Usage:
      superflowsql config              # Show all config
      superflowsql config POSTGRES_PORT # Show one value
      superflowsql config POSTGRES_PORT 5433  # Update a value
    """
    project = _find_project()

    try:
        entries = load_config(project)
    except FileNotFoundError as e:
        console.print(f"[red]✗ {e}[/red]")
        sys.exit(1)

    # Update mode
    if key and value:
        try:
            update_config(project, key, value)
            console.print(f"[green]✓ Updated {key} = {value}[/green]")
        except KeyError as e:
            console.print(f"[red]✗ {e}[/red]")
            sys.exit(1)
        return

    # Show single key
    if key:
        for entry in entries:
            if entry.key == key:
                console.print(f"[bold]{entry.key}[/bold] = [cyan]{entry.value}[/cyan]")
                return
        console.print(f"[red]✗ Key not found: {key}[/red]")
        sys.exit(1)

    # Show all config
    table = Table(title="Configuration", show_header=True, header_style="bold cyan")
    table.add_column("Group", style="dim")
    table.add_column("Key", style="bold")
    table.add_column("Value", style="cyan")

    current_group = ""
    for entry in entries:
        group_display = entry.group if entry.group != current_group else ""
        current_group = entry.group
        table.add_row(group_display, entry.key, entry.value)

    console.print(table)


# ──────────────────────────────────────────────
# superflowsql tui
# ──────────────────────────────────────────────
@main.command()
def tui():
    """Launch the interactive TUI (Go-based terminal interface)."""
    import platform
    import shutil
    import subprocess

    # Look for the Go binary
    binary_name = "superflowsql-tui.exe" if platform.system() == "Windows" else "superflowsql-tui"
    binary = shutil.which(binary_name)

    if not binary:
        # Check relative to this package
        pkg_dir = Path(__file__).parent.parent / "src" / "tui"
        candidates = [
            pkg_dir / binary_name,
            pkg_dir / "superflowsql.exe",
            pkg_dir / "superflowsql",
        ]
        for c in candidates:
            if c.exists():
                binary = str(c)
                break

    if not binary:
        console.print(
            "[red]✗ TUI binary not found.[/red]"
            "\n  Build it with: [cyan]cd src/tui && go build -o superflowsql-tui .[/cyan]"
        )
        sys.exit(1)

    subprocess.run([binary])


if __name__ == "__main__":
    main()
