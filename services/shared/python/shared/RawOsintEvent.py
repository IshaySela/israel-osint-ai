from dataclasses import dataclass
from datetime import date


@dataclass
class RawOsintEvent:
    """Base class for all unprocessed osint events
    """
    id: str
    date: date
    source: str