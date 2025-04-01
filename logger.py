"""
Modeule for good looking logger.
"""

from logging import Logger, getLogger, Formatter, INFO
from logging.handlers import RotatingFileHandler
from os import getcwd
from os.path import join


def get_my_logger(name: str, path: str, log_level: int | str = 0) -> Logger:
    """
    Returns file logger with time and level message.
    """

    logger = getLogger(name)
    filename = join(path, f"{name}.log")
    handler = RotatingFileHandler(
        filename,
        mode="a",
        maxBytes=1024 * 1024 * 20,
        backupCount=1,
        encoding="utf-8",
    )
    formatter = Formatter(
        fmt="[%(asctime)s] %(levelname)s %(message)s",
        datefmt="%Y.%m.%d %H:%M:%S",
    )
    handler.setFormatter(formatter)
    logger.addHandler(handler)
    logger.setLevel(log_level)
    return logger


LOGGER = get_my_logger("log", join(getcwd(), "logs"), INFO)
