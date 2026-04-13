from pydantic import BaseModel, Field
from shared.RawOsintEvent import RawOsintEvent


class TelegramEventData(BaseModel):
    event_type: str
    chat_id: int
    channel_title: str
    channel_main_lang: str
    msg_id: int


class RawTelegramEvent(RawOsintEvent):
    source: str = Field(default="telegram", init=False)
    data: TelegramEventData