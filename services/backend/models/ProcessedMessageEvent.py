from pydantic import BaseModel, Field
from typing import Dict

class Geocode(BaseModel):
    lat: float
    lng: float

class ProcessedEventMessage(BaseModel):
    db_id: str = Field(..., alias="dbId")
    summary: str
    locations: Dict[str, Geocode]
    timestamp: str