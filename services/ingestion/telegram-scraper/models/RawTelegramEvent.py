from shared.RawOsintEvent import RawOsintEvent
from pydantic import Field

class RawTelegramEvent(RawOsintEvent):
    text: str
    event_type: str
    chat_id: int
    channel_title: str
    message_id: int
    source: str = Field(default="telegram",init=False)