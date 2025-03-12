import hashlib
import logging
import os
import shutil
import subprocess
import sys
from typing import Dict, Optional, Tuple

import virtualenv

from .env_vars import MAMBA_BIN_PATH
from .path import CHAT_CONFIG_FILENAME, CHAT_DIR, ENV_CACHE_DIR, MAMBA_PY_ENVS, MAMBA_ROOT
from .schema import ExternalPyConf
from .user_setting import USER_SETTINGS

logger = logging.getLogger(__name__)

PYPI_TUNA = "https://pypi.tuna.tsinghua.edu.cn/simple"
DEFAULT_CONDA_FORGE_URL = "https://mirrors.tuna.tsinghua.edu.cn/anaconda/cloud/conda-forge/"


def _get_external_envs() -> Dict[str, ExternalPyConf]:
    """
    Get the external python environments info from the user settings.
    """
    external_pythons: Dict[str, ExternalPyConf] = {}
    for env_name, python_bin in USER_SETTINGS.environments.items():
        external_pythons[env_name] = ExternalPyConf(env_name=env_name, py_bin=python_bin)

    return external_pythons


EXTERNAL_ENVS = _get_external_envs()


class PyEnvManager:
    mamba_bin = MAMBA_BIN_PATH
    mamba_root = MAMBA_ROOT

    def __init__(self):
        pass

    @staticmethod
    def get_py_version(py: str) -> Optional[str]:
        """
        Get the version of the python executable.
        """
        py_version_cmd = [py, "--version"]
        with subprocess.Popen(py_version_cmd, stdout=subprocess.PIPE, stderr=None) as proc:
            proc.wait()

            if proc.returncode != 0:
                return None

            out = proc.stdout.read().decode("utf-8")
            return out.split()[1]

    @staticmethod
    def get_dep_hash(reqirements_file: str) -> str:
        """
        Get the hash of the requirements file content.

        Used to check if the requirements file has been changed.
        """
        with open(reqirements_file, "r", encoding="utf-8") as f:
            content = f.read()
            return hashlib.md5(content.encode("utf-8")).hexdigest()

    def ensure(
        self,
        env_name: str,
        py_version: Optional[str] = None,
        reqirements_file: Optional[str] = None,
    ) -> Optional[str]:
        """
        Ensure the python environment exists with the given name and version.
        And install the requirements if provided.

        return the python executable path.
        """
        py = self.get_py(env_name)

        should_remove_old = False
        should_install_deps = False

        if py:
            # check the version of the python executable
            current_version = self.get_py_version(py)

            if py_version and current_version != py_version:
                should_remove_old = True

            if reqirements_file and self.should_reinstall(env_name, reqirements_file):
                should_install_deps = True

            if not should_remove_old and not should_install_deps:
                return py

        # log_file = get_logging_file()

        if should_remove_old:
            self.remove(env_name, py_version)

        # create the environment if it doesn't exist or needs to be recreated
        if should_remove_old or not py:
            if py_version:
                create_ok, msg = self.create(env_name, py_version)
            else:
                create_ok, msg = self.create_with_virtualenv(env_name)

            if not create_ok:
                logger.error(f"Failed to create {env_name}: {msg}")
                sys.exit(0)

        # install or update the requirements
        if reqirements_file:
            filename = os.path.basename(reqirements_file)
            action = "Updating" if should_install_deps else "Installing"
            logger.debug(f"- {action} dependencies from {filename}...")
            install_ok, msg = self.install(env_name, reqirements_file)
            if not install_ok:
                logger.error(f"Failed to {action.lower()} dependencies: {msg}")
                sys.exit(0)
            else:
                # save the hash of the requirements file content
                dep_hash = self.get_dep_hash(reqirements_file)
                cache_file = os.path.join(ENV_CACHE_DIR, f"{env_name}")
                with open(cache_file, "w", encoding="utf-8") as f:
                    f.write(dep_hash)

        return self.get_py(env_name)

    def create_with_virtualenv(self, env_name: str) -> Tuple[bool, str]:
        """
        Create a new python environment using virtualenv with the current Python interpreter.
        """
        env_path = os.path.join(MAMBA_PY_ENVS, env_name)
        if os.path.exists(env_path):
            return True, ""

        try:
            # Use virtualenv.cli_run to create a virtual environment
            virtualenv.cli_run([env_path, "--python", sys.executable])

            # Create sitecustomize.py in the lib/site-packages directory
            site_packages_dir = os.path.join(env_path, "Lib", "site-packages")
            if os.path.exists(site_packages_dir):
                sitecustomize_path = os.path.join(site_packages_dir, "sitecustomize.py")
                with open(sitecustomize_path, "w") as f:
                    f.write("import sys\n")
                    f.write('sys.path = [path for path in sys.path if path.find("conda") == -1]')

            return True, ""
        except Exception as e:
            return False, str(e)

    def install(self, env_name: str, requirements_file: str) -> Tuple[bool, str]:
        """
        Install or update requirements in the python environment.

        Args:
            env_name: the name of the python environment
            requirements_file: the absolute path to the requirements file.

        Returns:
            A tuple (success, message), where success is a boolean indicating
            whether the installation was successful, and message is a string
            containing output or error information.
        """
        py = self.get_py(env_name)
        if not py:
            return False, "Python executable not found."

        if not os.path.exists(requirements_file):
            return False, "Dependencies file not found."

        # Base command
        cmd = [
            py,
            "-m",
            "pip",
            "install",
            "-r",
            requirements_file,
            "-i",
            PYPI_TUNA,
            "--no-warn-script-location",
        ]

        # Check if this is an update or a fresh install
        cache_file = os.path.join(ENV_CACHE_DIR, f"{env_name}")
        if os.path.exists(cache_file):
            # This is an update, add --upgrade flag
            cmd.append("--upgrade")

        env = os.environ.copy()
        env.pop("PYTHONPATH", None)
        with subprocess.Popen(
            cmd, stdout=subprocess.DEVNULL, stderr=subprocess.PIPE, env=env
        ) as proc:
            _, err = proc.communicate()

            if proc.returncode != 0:
                return False, f"Installation failed: {err.decode('utf-8')}"

            return True, "Installation successful"

    def should_reinstall(self, env_name: str, requirements_file: str) -> bool:
        """
        Check if the requirements file has been changed.
        """
        cache_file = os.path.join(ENV_CACHE_DIR, f"{env_name}")
        if not os.path.exists(cache_file):
            return True

        dep_hash = self.get_dep_hash(requirements_file)
        with open(cache_file, "r", encoding="utf-8") as f:
            cache_hash = f.read()

        return dep_hash != cache_hash

    def create(self, env_name: str, py_version: str) -> Tuple[bool, str]:
        """
        Create a new python environment using mamba.
        """
        is_exist = os.path.exists(os.path.join(MAMBA_PY_ENVS, env_name))
        if is_exist:
            return True, ""

        # Get conda-forge URL from config file
        conda_forge_url = self._get_conda_forge_url()

        # create the environment
        cmd = [
            self.mamba_bin,
            "create",
            "-n",
            env_name,
            "-c",
            conda_forge_url,
            "-r",
            self.mamba_root,
            f"python={py_version}",
            "-y",
        ]
        with subprocess.Popen(cmd, stdout=subprocess.PIPE, stderr=subprocess.PIPE) as proc:
            out, err = proc.communicate()
            msg = f"err: {err.decode()}\n-----\nout: {out.decode()}"

            if proc.returncode != 0:
                return False, msg
            return True, ""

    def remove(self, env_name: str, py_version: Optional[str] = None) -> bool:
        if py_version:
            return self.remove_by_mamba(env_name)
        return self.remove_by_del(env_name)

    def remove_by_del(self, env_name: str) -> bool:
        """
        Remove the python environment.
        """
        env_path = os.path.join(MAMBA_PY_ENVS, env_name)
        try:
            # Remove the environment directory
            if os.path.exists(env_path):
                shutil.rmtree(env_path)
                return True
        except Exception as e:
            logger.warning(f"Failed to remove environment {env_name}: {e}")
            return False

    def remove_by_mamba(self, env_name: str) -> bool:
        """
        Remove the python environment.
        """
        is_exist = os.path.exists(os.path.join(MAMBA_PY_ENVS, env_name))
        if not is_exist:
            return True

        # remove the environment
        cmd = [
            self.mamba_bin,
            "env",
            "remove",
            "-n",
            env_name,
            "-r",
            self.mamba_root,
            "-y",
        ]
        with subprocess.Popen(cmd, stdout=subprocess.DEVNULL, stderr=None) as proc:
            proc.wait()

            if proc.returncode != 0:
                return False

            return True

    def get_py(self, env_name: str) -> Optional[str]:
        """
        Get the python executable path of the given environment.
        """
        env_path = None
        if sys.platform == "win32":
            env_path = os.path.join(MAMBA_PY_ENVS, env_name, "python.exe")
            if not os.path.exists(env_path):
                env_path = os.path.join(MAMBA_PY_ENVS, env_name, "Scripts", "python.exe")
        else:
            env_path = os.path.join(MAMBA_PY_ENVS, env_name, "bin", "python")

        if env_path and os.path.exists(env_path):
            return env_path

        return None

    def _get_conda_forge_url(self) -> str:
        """
        Read the conda-forge URL from the config file.
        If the config file does not exist or does not contain the conda-forge URL,
        use the default value.
        """
        config_file = os.path.join(CHAT_DIR, CHAT_CONFIG_FILENAME)

        try:
            if not os.path.exists(config_file):
                return DEFAULT_CONDA_FORGE_URL

            import yaml

            with open(config_file, "r", encoding="utf-8") as f:
                config = yaml.safe_load(f)

            return config.get("conda-forge-url", DEFAULT_CONDA_FORGE_URL)
        except Exception as e:
            logger.error(f"An error occurred when loading conda-forge-url from config file: {e}")
            return DEFAULT_CONDA_FORGE_URL
