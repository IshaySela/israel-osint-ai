from dataclasses import dataclass
from typing import Literal, get_args, TypeGuard, TypedDict
from openai import AsyncOpenAI
from .Configuration import TelegramScraperConfig

config = TelegramScraperConfig.get()

EventTypes = Literal['rocket_fire', 'shooting', 'attack', 'missle_hit','not_relevant']

@dataclass(frozen=True)
class EventsDescription(TypedDict):
    rocket_fire: str
    shooting: str
    attack: str
    missle_hit: str
    not_relevant: str
    
eventsDescription: EventsDescription = {
    "rocket_fire": "The text indicates the launch, interception, or impact of rockets, missiles, or mortar fire.",
    "shooting": "The text indicates a kinetic engagement involving firearms or small arms fire.",
    "missle_hit": "Any place that was hit by a missile or rocket.",
    "attack": "Hostile acts involving physical assault, stabbings, vehicle rammings, or complex tactical incursions not covered by specific projectile or firearm labels.",
    "not_relevant": "The content does not meet the criteria for any defined tactical event labels."
}

def _create_prompt_mappings(ed: EventsDescription) -> str:
    """
    Generates a formatted string mapping event types to their descriptions for use in prompts.

    Args:
        ed (EventsDescription): A dictionary containing event types as keys and their descriptions as values.

    Returns:
        str: A newline-separated string where each line follows the format '-key: description'.
    """

    lines = [f"-{x}: {ed[x]}" for x in ed.keys()]
    
    return "\n".join(lines)

def _is_valid_event(event: str) -> TypeGuard[EventTypes]:
    valid_events = eventsDescription.keys()
    return event in valid_events

client = AsyncOpenAI(
    api_key=config.openai_api_key
)

eventsDescriptionPrompt = _create_prompt_mappings(eventsDescription)

developerPrompt = f"""
You are a specialized Natural Language Processing classifier optimized for analysis of Hebrew security alerts.

Goal: Perform a multiclass classification task on provided Hebrew text strings.
Map each input to exactly one of the defined labels.

Label Definitions:
{eventsDescriptionPrompt}

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