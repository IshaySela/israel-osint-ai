from shared.RawOsintEvent import RawOsintEvent
from dataclasses import dataclass

@dataclass
class RawTelegramEvent(RawOsintEvent):
    text: str
    event_type: str
    chat_id: str
    message_id: str