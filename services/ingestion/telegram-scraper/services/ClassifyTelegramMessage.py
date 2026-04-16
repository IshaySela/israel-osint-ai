from dataclasses import dataclass
from typing import Annotated
from openai import AsyncOpenAI
from .Configuration import TelegramScraperConfig
from enum import Enum
from pydantic import BaseModel, Field

class EventTypes(str,Enum):
    """The text indicates the launch or interception of rockets, missiles, or mortar fire."""
    rocket_fire = "rocket_fire"
    """The text indicates a usage of firearms fire in a security settings (not criminal)"""
    shooting = "shooting"
    """Hostile acts involving physical assault, stabbings, vehicle rammings, or complex tactical incursions not covered by specific projectile or firearm labels."""
    attack = "attack"
    """Any place that was hit by a missile or rocket."""
    missile_hit = "missile_hit"
    """The content does not meet the criteria for any defined labels, or is describing a criminal event"""
    not_relevant = "not_relevant"

def build_events_description() -> str:
    descriptions = {}
    
    for event in EventTypes:
        if event.__doc__ is None:
            raise RuntimeError(f"Event {event} is missing a description docstring.")
        descriptions[event] = event.__doc__.strip()
    
    return "\n".join([f"{key.value}: {value}" for key, value in descriptions.items()])


class EventClassifierResponse(BaseModel):
    event_type: Annotated[EventTypes,
                           Field(description=build_events_description())]


config = TelegramScraperConfig.get()

client = AsyncOpenAI(
    api_key=config.openai_api_key
)


developerPrompt = f"""
You are a specialized Natural Language Processing classifier optimized for analysis of Hebrew security alerts.

Goal: Perform a multiclass classification task on provided Hebrew text strings.
Map each input to exactly one of the defined labels.
"""

async def classify_telegram_msg(message: str) -> EventTypes:
    result = await client.responses.parse(
        input=[
            { "role": "system", "content": developerPrompt },
            {"role": "user", "content": message }
        ],
        model="gpt-5-nano-2025-08-07",
        text_format=EventClassifierResponse
    )
    
    if result.output_parsed is None:
        return EventTypes.not_relevant
    
    return result.output_parsed.event_type