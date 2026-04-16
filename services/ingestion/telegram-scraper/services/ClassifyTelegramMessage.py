from dataclasses import dataclass
from typing import Literal, get_args, TypeGuard, TypedDict
from openai import AsyncOpenAI
from .Configuration import TelegramScraperConfig
from enum import Enum
from pydantic import BaseModel

class EventTypes(str,Enum):
    rocket_fire = "rocket_fire"
    """The text indicates the launch or interception of rockets, missiles, or mortar fire."""
    shooting = "shooting"
    """The text indicates a usage of firearms fire in a security settings (not criminal)"""
    attack = "attack"
    """Hostile acts involving physical assault, stabbings, vehicle rammings, or complex tactical incursions not covered by specific projectile or firearm labels."""
    missile_hit = "missile_hit"
    """Any place that was hit by a missile or rocket."""
    not_relevant = "not_relevant"
    """The content does not meet the criteria for any defined tactical event labels."""

class EventClassifierResponse(BaseModel):
    event_type: EventTypes

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