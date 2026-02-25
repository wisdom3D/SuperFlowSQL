"""Docker Compose operations - start, stop, status for the data stack."""

from __future__ import annotations

import shutil
import subprocess
from dataclasses import dataclass
from pathlib import Path


@dataclass
class ServiceStatus:
    """Status of a single Docker service."""

    name: str
    status: str  # running, stopped, not_found
    health: str  # healthy, unhealthy, starting, -
    ports: str

    @property
    def is_running(self) -> bool:
        return self.status == "running"

    @property
    def icon(self) -> str:
        if self.status == "running":
            return "●"
        elif self.status == "stopped":
            return "○"
        return "?"


def is_docker_running() -> bool:
    """Check if the Docker daemon is accessible."""
    try:
        subprocess.run(
            ["docker", "info"],
            capture_output=True,
            timeout=10,
        )
        return True
    except (subprocess.SubprocessError, FileNotFoundError):
        return False


def _docker_compose_cmd(project_dir: str | Path, *args: str) -> subprocess.CompletedProcess:
    """Run a docker compose command in the given project directory.

    Tries 'docker compose' (V2) first, falls back to 'docker-compose' (V1).
    """
    project_dir = str(Path(project_dir).resolve())

    # Try Docker Compose V2
    try:
        result = subprocess.run(
            ["docker", "compose", *args],
            cwd=project_dir,
            capture_output=True,
            text=True,
            timeout=300,
        )
        if result.returncode == 0:
            return result
    except FileNotFoundError:
        pass

    # Fallback to V1
    if shutil.which("docker-compose"):
        return subprocess.run(
            ["docker-compose", *args],
            cwd=project_dir,
            capture_output=True,
            text=True,
            timeout=300,
        )

    raise RuntimeError(
        "Docker Compose not found. Install Docker Desktop or docker-compose."
    )


def start_stack(project_dir: str | Path) -> str:
    """Start all services with docker compose up -d --build.

    Args:
        project_dir: Path to the SuperFlowSQL project directory.

    Returns:
        Combined stdout + stderr output from docker compose.
    """
    result = _docker_compose_cmd(project_dir, "up", "-d", "--build")
    output = (result.stdout + "\n" + result.stderr).strip()
    if result.returncode != 0:
        raise RuntimeError(f"Failed to start stack:\n{output}")
    return output


def stop_stack(project_dir: str | Path) -> str:
    """Stop all services with docker compose down.

    Args:
        project_dir: Path to the SuperFlowSQL project directory.

    Returns:
        Combined stdout + stderr output.
    """
    result = _docker_compose_cmd(project_dir, "down")
    output = (result.stdout + "\n" + result.stderr).strip()
    if result.returncode != 0:
        raise RuntimeError(f"Failed to stop stack:\n{output}")
    return output


def restart_stack(project_dir: str | Path) -> str:
    """Restart all services (down then up)."""
    stop_stack(project_dir)
    return start_stack(project_dir)


def get_stack_status(project_dir: str | Path) -> list[ServiceStatus]:
    """Get the status of all services in the stack.

    Returns:
        List of ServiceStatus for each running/stopped service.
    """
    try:
        result = _docker_compose_cmd(
            project_dir, "ps", "--format", "{{.Name}}|{{.Status}}|{{.Ports}}"
        )
    except RuntimeError:
        return _default_services()

    services: list[ServiceStatus] = []
    for line in result.stdout.strip().splitlines():
        line = line.strip()
        if not line:
            continue

        parts = line.split("|", 2)
        name = parts[0] if len(parts) >= 1 else "unknown"
        raw_status = parts[1].lower() if len(parts) >= 2 else ""
        ports = parts[2] if len(parts) >= 3 else ""

        if "up" in raw_status:
            status = "running"
            if "healthy" in raw_status:
                health = "healthy"
            elif "unhealthy" in raw_status:
                health = "unhealthy"
            else:
                health = "starting"
        else:
            status = "stopped"
            health = "-"

        services.append(ServiceStatus(name=name, status=status, health=health, ports=ports))

    return services if services else _default_services()


def _default_services() -> list[ServiceStatus]:
    """Return default expected services as not found."""
    return [
        ServiceStatus(name="postgres", status="not_found", health="-", ports=""),
        ServiceStatus(name="pgadmin", status="not_found", health="-", ports=""),
        ServiceStatus(name="airflow-webserver", status="not_found", health="-", ports=""),
        ServiceStatus(name="airflow-scheduler", status="not_found", health="-", ports=""),
        ServiceStatus(name="superset", status="not_found", health="-", ports=""),
    ]


def get_service_urls(config: dict | None = None) -> dict[str, str]:
    """Return default service access URLs.

    Args:
        config: Optional dict with port overrides (airflow_port, pgadmin_port, etc.)

    Returns:
        Dict mapping service names to their URLs.
    """
    cfg = config or {}
    return {
        "Airflow": f"http://localhost:{cfg.get('airflow_port', '8080')}",
        "PgAdmin": f"http://localhost:{cfg.get('pgadmin_port', '5050')}",
        "Superset": f"http://localhost:{cfg.get('superset_port', '8088')}",
        "PostgreSQL": f"localhost:{cfg.get('postgres_port', '5432')}",
    }
