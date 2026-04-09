from dataclasses import dataclass
from datetime import date


@dataclass
class RawOsintEvent:
    date: date
    source: str