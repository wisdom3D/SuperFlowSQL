"""Configuration management - read/write .env files for the data stack."""

from __future__ import annotations

from dataclasses import dataclass
from pathlib import Path

from dotenv import dotenv_values


@dataclass
class ConfigEntry:
    """A single configuration key-value pair."""

    key: str
    value: str
    group: str = "General"

    @staticmethod
    def get_group(key: str) -> str:
        """Infer the config group from the key prefix."""
        prefix = key.split("_")[0].upper()
        groups = {
            "POSTGRES": "PostgreSQL",
            "AIRFLOW": "Airflow",
            "PGADMIN": "PgAdmin",
            "SUPERSET": "Superset",
        }
        return groups.get(prefix, "General")


def load_config(project_dir: str | Path) -> list[ConfigEntry]:
    """Load configuration from the project's .env file.

    Args:
        project_dir: Path to the SuperFlowSQL project directory.

    Returns:
        List of ConfigEntry objects.

    Raises:
        FileNotFoundError: If the .env file doesn't exist.
    """
    env_path = Path(project_dir).resolve() / ".env"

    if not env_path.exists():
        raise FileNotFoundError(f"No .env file found at: {env_path}")

    values = dotenv_values(env_path)
    entries = []
    for key, value in values.items():
        if value is not None:
            entries.append(
                ConfigEntry(
                    key=key,
                    value=value,
                    group=ConfigEntry.get_group(key),
                )
            )

    return entries


def save_config(project_dir: str | Path, entries: list[ConfigEntry]) -> None:
    """Save configuration entries back to the .env file.

    Args:
        project_dir: Path to the SuperFlowSQL project directory.
        entries: List of ConfigEntry objects to write.
    """
    env_path = Path(project_dir).resolve() / ".env"

    # Group entries
    groups: dict[str, list[ConfigEntry]] = {}
    group_order: list[str] = []
    for entry in entries:
        if entry.group not in groups:
            groups[entry.group] = []
            group_order.append(entry.group)
        groups[entry.group].append(entry)

    lines = [
        "# ============================================",
        "# SuperFlowSQL - Environment Configuration",
        "# ============================================",
        "",
    ]
    for group_name in group_order:
        lines.append(f"# --- {group_name} ---")
        for entry in groups[group_name]:
            lines.append(f"{entry.key}={entry.value}")
        lines.append("")

    env_path.write_text("\n".join(lines), encoding="utf-8")


def update_config(project_dir: str | Path, key: str, value: str) -> None:
    """Update a single configuration value.

    Args:
        project_dir: Path to the SuperFlowSQL project directory.
        key: The configuration key to update.
        value: The new value.

    Raises:
        KeyError: If the key doesn't exist in the config.
    """
    entries = load_config(project_dir)
    found = False
    for entry in entries:
        if entry.key == key:
            entry.value = value
            found = True
            break

    if not found:
        raise KeyError(f"Configuration key not found: {key}")

    save_config(project_dir, entries)


def find_project_dir(start_dir: str | Path = ".") -> Path | None:
    """Find a SuperFlowSQL project directory by walking up from start_dir.

    Looks for a directory containing both docker-compose.yml and .env.

    Args:
        start_dir: Starting directory for the upward search.

    Returns:
        Path to the project directory, or None if not found.
    """
    current = Path(start_dir).resolve()

    while True:
        if (current / "docker-compose.yml").exists() and (current / ".env").exists():
            return current
        parent = current.parent
        if parent == current:
            break
        current = parent

    return None
