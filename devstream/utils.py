import locale


def get_sys_language_code():
    """
    Return the two-letter language code following ISO 639-1.
    """
    lang, _ = locale.getdefaultlocale()
    return lang.split("_")[0].lower()
