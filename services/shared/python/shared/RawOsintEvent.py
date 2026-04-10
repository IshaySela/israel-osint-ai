from pydantic import BaseModel

class RawOsintEvent(BaseModel):
    """Base class for all unprocessed osint events
    """
    id: str
    timestamp: str
    source: str