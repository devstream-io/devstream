import logging
import os
import sys
from typing import Dict, List, Optional, Tuple

import oyaml as yaml

from .env_manager import EXTERNAL_ENVS
from .namespace import get_prioritized_namespace_path
from .path import INTERFACE_FILENAMES
from .schema import RuntimeParameter, WorkflowConfig
from .step import WorkflowStep
from .utils import get_sys_language_code

logger = logging.getLogger(__name__)


class Workflow:
    TRIGGER_PREFIX = "/"
    HELP_FLAG_PREFIX = "--help"

    def __init__(self, config: WorkflowConfig):
        self._config = config

        self._runtime_param = None

    @property
    def config(self) -> WorkflowConfig:
        return self._config

    @property
    def runtime_param(self):
        return self._runtime_param

    @staticmethod
    def parse_trigger(user_input: str) -> Tuple[Optional[str], Optional[str]]:
        """
        Check if the user input should trigger a workflow.
        Return a tuple of (workflow_name, the input without workflow trigger).

        User input is considered a workflow trigger if it starts with the Workflow.PREFIX.
        The workflow name is the first word after the prefix.
        """
        striped = user_input.strip()
        if not striped:
            return None, user_input
        if striped[0] != Workflow.TRIGGER_PREFIX:
            return None, user_input

        workflow_name = striped.split()[0][1:]

        # remove the trigger prefix and the workflow name
        actual_input = user_input.replace(f"{Workflow.TRIGGER_PREFIX}{workflow_name}", "", 1)
        return workflow_name, actual_input

    @staticmethod
    def load(workflow_name: str) -> Optional["Workflow"]:
        """
        Load a workflow from the interface.yml by name.
        A workflow name is the relative path of interface.yml
        to the /workflows dir joined by "."
        e.g
        - "unit_tests": means the interface file of the workflow is unit_tests/interface.yml
        - "commit.en": means the interface file is commit/en/interface.yml
        - "pr.review.zh": means the interface file is pr/review/zh/interface.yml

        TODO: 单个路径组件合法的正则表达式为：`^[A-Za-z0-9_-]{1,255}$`
        """
        path_parts = workflow_name.split(".")
        if len(path_parts) < 1:
            return None

        rel_path = os.path.join(*path_parts)

        found = False
        workflow_dir = ""
        prioritized_dirs = get_prioritized_namespace_path()
        for wf_dir in prioritized_dirs:
            for fn in INTERFACE_FILENAMES:
                yaml_file = os.path.join(wf_dir, rel_path, fn)
                if os.path.exists(yaml_file):
                    workflow_dir = wf_dir
                    found = True
                    break
            if found:
                break
        if not found:
            return None

        # Load and override yaml conf in top-down order
        config_dict = {}
        for i in range(len(path_parts)):
            cur_path = os.path.join(workflow_dir, *path_parts[: i + 1])
            for fn in INTERFACE_FILENAMES:
                cur_yaml = os.path.join(cur_path, fn)

                if os.path.exists(cur_yaml):
                    with open(cur_yaml, "r", encoding="utf-8") as file:
                        yaml_content = file.read()
                        cur_conf = yaml.safe_load(yaml_content)
                        cur_conf["dirpath"] = cur_path

                    # convert relative path to absolute path for dependencies file
                    if cur_conf.get("runtime", {}).get("requirements"):
                        rel_dep = cur_conf["runtime"]["requirements"]
                        abs_dep = os.path.join(cur_path, rel_dep)
                        cur_conf["runtime"]["requirements"] = abs_dep

                    config_dict.update(cur_conf)

        config = WorkflowConfig.model_validate(config_dict)

        return Workflow(config)

    def setup(
        self,
        user_input: Optional[str],
        model_name: Optional[str] = None,
        history_messages: Optional[List[Dict]] = None,
        parent_hash: Optional[str] = None,
    ):
        """
        Setup the workflow with the runtime parameters and env variables.
        """
        # NOTE: prepare an internal default python if needed
        default_python = sys.executable

        workflow_py = None
        if self.config.runtime.venv:
            env_name = self.config.runtime.venv
            if env_name in EXTERNAL_ENVS:
                # Use the external python set in the user settings
                workflow_py = EXTERNAL_ENVS[env_name].py_bin
                msg = [
                    "Using external Python from user settings:",
                    f"- env_name: {env_name}",
                    f"- python_bin: {workflow_py}",
                    "This Python environment's version and dependencies should be "
                    "ensured by the user to meet the requirements.",
                ]
                logger.debug("\n".join(msg))

            else:
                # version = self.config.runtime.python_version
                # req_file = self.config.runtime.requirements
                # manager = PyEnvManager()
                # workflow_py = manager.ensure(env_name, version, req_file)
                pass
        else:
            pass

        workflow_py = workflow_py or default_python
        runtime_param = {
            # from user interaction
            "user_input": user_input,
            "model_name": model_name,
            "history_messages": history_messages,
            "parent_hash": parent_hash,
            # from user setting or system
            "workflow_python": workflow_py,
        }

        self._runtime_param = RuntimeParameter.model_validate(runtime_param)

    def run(self) -> int:
        """
        Run the steps of the workflow.
        """
        steps = self.config.runtime.run

        for s in steps:
            step = WorkflowStep(cmd=s)
            result = step.run(self.config, self.runtime_param)
            return_code = result[0]
            if return_code != 0:
                # stop the workflow if any step fails
                return return_code

        return 0

    def get_help_doc(self, user_input: str) -> Optional[str]:
        """
        Get the help doc content of the workflow.
        """
        help_info = self.config.help
        help_file = None

        if isinstance(help_info, str):
            # return the only help doc
            help_file = help_info

        elif isinstance(help_info, list) and len(help_info) > 0:
            first = help_info[0]
            assert isinstance(first, dict)
            default_file = list(first.values())[0]

            # get language code from user input
            code = user_input.strip().removeprefix(Workflow.HELP_FLAG_PREFIX)
            code = code.removeprefix(".").strip()
            code = code.lower()
            code = code or get_sys_language_code()
            help_info_dict = {list(d.keys())[0].lower(): list(d.values())[0] for d in help_info}
            help_file = help_info_dict.get(code, default_file)

        if not help_file:
            return None  # or raise error? no help file configured

        help_path = os.path.join(self.config.dirpath, help_file)
        if os.path.exists(help_path):
            with open(help_path, "r", encoding="utf-8") as file:
                return file.read()

        return None  # or raise error? help file not found

    def should_show_help(self, user_input) -> bool:
        return user_input.strip().startswith(Workflow.HELP_FLAG_PREFIX)
