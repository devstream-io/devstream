"""
Namespace management for workflows
"""

import logging
import os
from pathlib import Path
from typing import Dict, List, Set, Tuple

import oyaml as yaml
import yaml as pyyaml
from pydantic import BaseModel, Field, ValidationError

from .path import (
    COMMUNITY_WORKFLOWS,
    INTERFACE_FILENAMES,
    SYS_WORKFLOWS,
    USER_BASE,
    USER_CONFIG_FILE,
)

logger = logging.getLogger(__name__)


class CustomConfig(BaseModel):
    namespaces: List[str] = []  # active namespaces ordered by priority


class WorkflowMeta(BaseModel):
    name: str = Field(..., description="workflow name")
    namespace: str = Field(..., description="workflow namespace")
    active: bool = Field(..., description="active flag")
    command_conf: Dict = Field(description="command configuration", default_factory=dict)

    def __str__(self):
        return f"{'*' if self.active else ' '} {self.name} ({self.namespace})"


def _load_custom_config() -> CustomConfig:
    """
    Load the custom config file.
    """
    config = CustomConfig()

    if not os.path.exists(USER_CONFIG_FILE):
        return config

    with open(USER_CONFIG_FILE, "r", encoding="utf-8") as file:
        content = file.read()
        yaml_content = yaml.safe_load(content)
        try:
            if yaml_content:
                config = CustomConfig.model_validate(yaml_content)
        except ValidationError as err:
            logger.warning("Invalid custom config file: %s", err)

    return config


def get_prioritized_namespace_path() -> List[str]:
    """
    Get the prioritized namespaces.

    priority: usr > sys > community
    """
    config = _load_custom_config()

    namespaces = config.namespaces

    namespace_paths = [os.path.join(USER_BASE, ns) for ns in namespaces]

    namespace_paths.append(SYS_WORKFLOWS)
    namespace_paths.append(COMMUNITY_WORKFLOWS)

    return namespace_paths


def iter_namespace(ns_path: str, existing_names: Set[str]) -> Tuple[List[WorkflowMeta], Set[str]]:
    """
    Get all workflows under the namespace path.

    Args:
        ns_path: the namespace path
        existing_names: the existing workflow names to check if the workflow is the first priority

    Returns:
        List[WorkflowMeta]: the workflows
        Set[str]: the updated existing workflow names
    """
    root = Path(ns_path)
    interest_files = set(INTERFACE_FILENAMES)
    result = []
    unique_names = set(existing_names)
    for file in root.rglob("*"):
        try:
            if file.is_file() and file.name in interest_files:
                rel_path = file.relative_to(root)
                parts = rel_path.parts
                workflow_name = ".".join(parts[:-1])
                is_first = workflow_name not in unique_names

                # load the config content from file
                with open(file, "r", encoding="utf-8") as file_handle:
                    yaml_content = file_handle.read()
                    command_conf = yaml.safe_load(yaml_content)
                    # pop the "steps" field
                    command_conf.pop("steps", None)

                workflow = WorkflowMeta(
                    name=workflow_name,
                    namespace=root.name,
                    active=is_first,
                    command_conf=command_conf,
                )
                unique_names.add(workflow_name)
                result.append(workflow)
        except pyyaml.scanner.ScannerError as err:
            logger.error(f"Failed to load {rel_path}: {err}")
        except Exception as err:
            logger.error(f"Unknown error when loading {rel_path}: {err}")

    return result, unique_names
