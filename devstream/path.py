import os

from .env_vars import DEVSTREAM_BASE

# -------------------------------
# devchat basic paths
# -------------------------------
USE_DIR = os.path.expanduser("~")
CHAT_DIR = os.path.join(USE_DIR, ".chat")
CHAT_CONFIG_FILENAME = "config.yml"
DEVCHAT_WORKFLOWS_BASE_NAME = "scripts"


# -------------------------------
# workflow paths
# -------------------------------
WORKFLOWS_BASE = DEVSTREAM_BASE or os.path.join(CHAT_DIR, DEVCHAT_WORKFLOWS_BASE_NAME)

SYS_WORKFLOWS = os.path.join(WORKFLOWS_BASE, "sys")
COMMUNITY_WORKFLOWS = os.path.join(WORKFLOWS_BASE, "comm")

INTERFACE_FILENAMES = [
    # the order matters
    "interface.yml",
    "interface.yaml",
    "command.yml",
    "command.yaml",
]


# -------------------------------
# workflow related cache data
# -------------------------------
CACHE_DIR = os.path.join(WORKFLOWS_BASE, "cache")
ENV_CACHE_DIR = os.path.join(CACHE_DIR, "env_cache")
os.makedirs(ENV_CACHE_DIR, exist_ok=True)


# -------------------------------
# custom/usr paths
# -------------------------------
USER_BASE = os.path.join(WORKFLOWS_BASE, "usr")
USER_CONFIG_FILE = os.path.join(USER_BASE, "config.yml")
USER_SETTINGS_FILE = os.path.join(WORKFLOWS_BASE, "settings.yml")


# ----  ---------------------------
#  Python environments paths
# -------------------------------
MAMBA_ROOT = os.path.join(CHAT_DIR, "mamba")
MAMBA_PY_ENVS = os.path.join(MAMBA_ROOT, "envs")
