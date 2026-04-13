from pydantic import BaseModel

class RawOsintEvent(BaseModel):
    """Base class for all unprocessed osint events
    """
    app_ev_id: str
    raw_message: str
    timestamp: str
    source: str