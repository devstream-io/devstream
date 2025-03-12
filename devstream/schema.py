import re
from typing import Dict, List, Optional, Union

from pydantic import BaseModel, ConfigDict, Field, computed_field, field_validator


class Parameters(BaseModel):
    input_mode: str = Field(default="optional", alias="__input__")  # "required" | "optional"


class RuntimeConf(BaseModel):
    id: str = "shell-python"
    run: List[str] = Field(default_factory=list)
    venv: Optional[str] = None
    # python: Optional[str] = None
    python_version: Optional[str] = Field(
        default=None, alias="python"
    )  # the specific python version.
    requirements: Optional[str] = None  # absolute path to the requirements file

    @field_validator("python_version")
    @classmethod
    def validate_python_version(cls, value: Optional[str]) -> Optional[str]:
        if value is None:
            return value

        pattern = r"^\d+\.\d+(\.\d+)?$"
        if not re.match(pattern, value):
            raise ValueError(f"Invalid version format: {value}. Should use the specific version.")
        return value


class WorkflowConfig(BaseModel):
    description: str
    runtime: RuntimeConf
    dirpath: str  # the path of the workflow dir
    parameters: Parameters = Field(default_factory=Parameters)
    help: Optional[Union[str, List[Dict[str, str]]]] = None

    @computed_field
    @property
    def input_required(self) -> bool:
        return self.parameters.input_mode.lower() == "required"

    model_config = ConfigDict(extra="ignore")


class UserSettings(BaseModel):
    environments: Dict[str, str] = Field(default_factory=dict)

    model_config = ConfigDict(extra="ignore")


class RuntimeParameter(BaseModel):
    workflow_python: str
    user_input: Optional[str] = None
    model_name: Optional[str] = None
    history_messages: Optional[Dict] = None
    parent_hash: Optional[str] = None

    model_config = ConfigDict(extra="ignore")


class ExternalPyConf(BaseModel):
    env_name: str  # the env_name of workflow python to act as
    py_bin: str  # the python executable path

    model_config = ConfigDict(extra="ignore")
