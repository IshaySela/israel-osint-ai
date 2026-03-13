import asyncio
from typing import Literal, get_args, TypeGuard
from openai import AsyncOpenAI
from .Configuration import TelegramScraperConfig

config = TelegramScraperConfig.get()

EventTypes = Literal['rocket_fire', 'shooting', 'not_relevant']


def _is_valid_event(event: str) -> TypeGuard[EventTypes]:
    valid_events = get_args(EventTypes)
    return event in valid_events

client = AsyncOpenAI(
    api_key=config.openai_api_key
)

developerPrompt = """
You classifiy hebrew messages as an event / not event as your job. Your job is to classifiy text to the following categories:
rocket_fire: The text describes a rocket fire event
shooting: The text describes a shooting event
not_relevant: Not a rocket_fire or shooting.
Answer only the category and nothing else.
"""

async def classify_telegram_msg(message: str) -> EventTypes:
    result = await client.responses.create(
        input=message,
        instructions=developerPrompt,
        model="gpt-5-nano-2025-08-07"
    )
    
    event = result.output_text
    
    if not _is_valid_event(event):
        return 'not_relevant'
    
    return event