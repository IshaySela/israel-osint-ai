import asyncio
from typing import Literal, get_args, TypeGuard
from openai import AsyncOpenAI
from .Configuration import TelegramScraperConfig

config = TelegramScraperConfig.get()

EventTypes = Literal['rocket_fire', 'shooting','attack', 'not_relevant']


def _is_valid_event(event: str) -> TypeGuard[EventTypes]:
    valid_events = get_args(EventTypes)
    return event in valid_events

client = AsyncOpenAI(
    api_key=config.openai_api_key
)

developerPrompt = """
You are a specialized Natural Language Processing classifier optimized for analysis of Hebrew security alerts.

Goal: Perform a multiclass classification task on provided Hebrew text strings.
Map each input to exactly one of the defined labels.

Label Definitions:
1. rocket_fire: The text indicates the launch, interception, or impact of rockets, missiles, or mortar fire.
2. shooting: The text indicates a kinetic engagement involving firearms or small arms fire.
3. not_relevant: Not a rocket_fire or shooting event.
4. attack: erate hostile acts involving physical assault, stabbings, vehicle rammings, or complex tactical incursions not covered by the above.

Output Constraints:
- Return only the label string.
- Do not include delimiters, preamble, or prose.
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