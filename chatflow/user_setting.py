from pathlib import Path

import oyaml as yaml

from .path import USER_SETTINGS_FILE
from .schema import UserSettings


def _load_user_settings() -> UserSettings:
    """
    Load the user settings from the settings.yml file.
    """
    settings_path = Path(USER_SETTINGS_FILE)
    if not settings_path.exists():
        return UserSettings()

    with open(settings_path, "r", encoding="utf-8") as file:
        content = yaml.safe_load(file.read())

    if content:
        return UserSettings.model_validate(content)

    return UserSettings()


USER_SETTINGS = _load_user_settings()
