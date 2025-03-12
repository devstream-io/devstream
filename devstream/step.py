import json
import os
import shlex
import subprocess
from enum import Enum
from typing import Dict, List, Tuple

from .path import WORKFLOWS_BASE
from .schema import RuntimeParameter, WorkflowConfig


class BuiltInVars(str, Enum):
    """
    Built-in variables within the workflow step command.
    """

    wf_python = "$__python__"
    wf_dirpath = "$__idir__"
    user_input = "$__input__"


class BuiltInEnvs(str, Enum):
    """
    Built-in environment variables for the step subprocess.
    """

    llm_model = "LLM_MODEL"
    parent_hash = "PARENT_HASH"
    context_contents = "CONTEXT_CONTENTS"


class WorkflowStep:
    def __init__(self, cmd: str):
        """
        Initialize a workflow step with the given configuration.
        """
        self._kwargs = {"run": cmd}

    @property
    def command_raw(self) -> str:
        """
        The raw command string from the config.
        """
        return self._kwargs.get("run", "")

    def _setup_env(self, wf_config: WorkflowConfig, rt_param: RuntimeParameter) -> Dict[str, str]:
        """
        Setup the environment variables for the subprocess.
        """
        env = os.environ.copy()

        # set PYTHONPATH for the subprocess
        new_paths = [WORKFLOWS_BASE]

        paths = [os.path.normpath(p) for p in new_paths]
        paths = [p.replace("\\", "\\\\") for p in paths]
        joined = os.pathsep.join(paths)

        env["PYTHONPATH"] = joined
        env[BuiltInEnvs.llm_model] = rt_param.model_name or ""
        env[BuiltInEnvs.parent_hash] = rt_param.parent_hash or ""
        env[BuiltInEnvs.context_contents] = ""
        if rt_param.history_messages:
            # convert dict to json string
            env[BuiltInEnvs.context_contents] = json.dumps(rt_param.history_messages)

        return env

    def _validate_and_interpolate(
        self, wf_config: WorkflowConfig, rt_param: RuntimeParameter
    ) -> List[str]:
        """
        Validate the step configuration and interpolate variables in the command.

        Return the command parts as a list of strings.
        """
        command_raw = self.command_raw
        parts = shlex.split(command_raw)

        args = []
        for p in parts:
            arg = p

            if p.startswith(BuiltInVars.wf_python):
                if not rt_param.workflow_python:
                    raise ValueError(
                        f"The command uses {BuiltInVars.wf_python}, "
                        "but the python path is not set yet."
                    )
                arg = arg.replace(BuiltInVars.wf_python, rt_param.workflow_python)

            if p.startswith(BuiltInVars.wf_dirpath):
                path_parts = p.split("/")
                # replace "$__idir__" with the root path in path_parts
                arg = os.path.join(wf_config.dirpath, *path_parts[1:])

            if BuiltInVars.user_input in p:
                arg = arg.replace(BuiltInVars.user_input, rt_param.user_input)

            args.append(arg)

        return args

    def run(self, wf_config: WorkflowConfig, rt_param: RuntimeParameter) -> Tuple[int, str, str]:
        """
        Run the step in a subprocess.

        Returns the return code, stdout, and stderr.
        """
        # setup the environment variables
        env = self._setup_env(wf_config, rt_param)

        command_args = self._validate_and_interpolate(wf_config, rt_param)

        # NOTE: handle stdout & stderr if needed
        # def _pipe_reader(pipe, data, out_file):
        #     """
        #     Read from the pipe, then write and save the data.
        #     """
        #     while pipe:
        #         pipe_data = pipe.read(1)
        #         if pipe_data == "":
        #             break
        #         data["data"] += pipe_data
        #         print(pipe_data, end="", file=out_file, flush=True)
        with subprocess.Popen(
            command_args,
            # stdout=subprocess.PIPE,
            # stderr=subprocess.PIPE,
            env=env,
            text=True,
        ) as proc:
            # stdout_data, stderr_data = {"data": ""}, {"data": ""}
            # stdout_thread = threading.Thread(
            #     target=_pipe_reader, args=(proc.stdout, stdout_data, sys.stdout)
            # )
            # stderr_thread = threading.Thread(
            #     target=_pipe_reader, args=(proc.stderr, stderr_data, sys.stderr)
            # )
            # stdout_thread.start()
            # stderr_thread.start()
            # stdout_thread.join()
            # stderr_thread.join()

            proc.wait()
            return_code = proc.returncode

            # return return_code, stdout_data["data"], stderr_data["data"]
            return return_code, "", ""
