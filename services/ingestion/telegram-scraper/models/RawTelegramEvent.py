from shared.RawOsintEvent import RawOsintEvent
from pydantic import Field

class RawTelegramEvent(RawOsintEvent):
    text: str
    event_type: str
    chat_id: int
    message_id: int
    source: str = Field(default="telegram",init=False)