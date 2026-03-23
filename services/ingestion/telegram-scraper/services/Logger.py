import logging
import sys
from types import FrameType
from typing import Final, Optional, Union, TYPE_CHECKING
from loguru import logger

class InterceptHandler(logging.Handler):
    """
    Bridges standard library logging to Loguru.
    """
    def emit(self, record: logging.LogRecord) -> None:
        try:
            level: Union[str, int] = logger.level(record.levelname).name
        except ValueError:
            level = record.levelno

        frame: Optional[FrameType] = logging.currentframe()
        depth: int = 2
        while frame and frame.f_code.co_filename == logging.__file__:
            frame = frame.f_back
            depth += 1

        logger.opt(depth=depth, exception=record.exc_info).log(level, record.getMessage())

def setup_logging() -> None:
    """
    Configures Loguru to output strictly to stdout.
    
    Args:
        level: Minimum log level to record.
        use_json: If True, outputs structured JSON. If False, outputs colorized text.
    """
    logger.remove()

    LOG_FORMAT: Final[str] = (
        "<green>{time:YYYY-MM-DD HH:mm:ss}</green> | "
        "<level>{level: <8}</level> | "
        "<cyan>{name}</cyan>:<cyan>{function}</cyan>:<cyan>{line}</cyan> - "
        "<level>{message}</level>"
    )
    
    logger.add(
        sys.stderr,
        format=LOG_FORMAT,
        colorize=True
    )
    
    for library in ["telethon", "pika", "asyncio", "urllib3"]:
        logging.getLogger(library).setLevel(logging.WARNING)

    logging.basicConfig(handlers=[InterceptHandler()], level=0, force=True)