from pydantic import BaseModel, Field
from typing import Dict

class Geocode(BaseModel):
    # Assuming Geocode has lat/lng; adjust fields as needed
    lat: float
    lng: float

class ProcessedEventMessage(BaseModel):
    db_id: str = Field(..., alias="dbId")
    summary: str
    locations: Dict[str, Geocode]
    timestamp: str

    class Config:
        # This allows you to populate the object using 'dbId' 
        # but refer to it as 'db_id' in your Python code.
        populate_by_name = True