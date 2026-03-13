import asyncio
from typing import Literal, get_args, TypeGuard
from openai import AsyncOpenAI
from Configuration import TelegramScraperConfig

config = TelegramScraperConfig.get()

type EventTypes = Literal['rocket fire'] | Literal['shooting'] |  Literal['not_relevant']

def _is_valid_event(event: str) -> TypeGuard[EventTypes]:
    return event in get_args(EventTypes)

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
        model="gpt5-nano"
    )
    
    event = result.output_text
    
    if not _is_valid_event(event):
        return 'not_relevant'
    
    return event