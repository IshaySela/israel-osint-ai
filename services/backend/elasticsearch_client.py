from typing import Any, Dict, List, Optional
from elasticsearch import Elasticsearch
from config import get_config, Config
from loguru import logger

class ESClient:
    client: Elasticsearch
    index: str

    def __init__(self) -> None:
        config: Config = get_config()
        self.client = Elasticsearch(config.elasticsearch_urls)
        self.index = config.elasticsearch_index

    def get_latest_events(self, size: int = 50) -> List[Dict[str, Any]]:
        try:
            query: Dict[str, Any] = {
                "query": {"match_all": {}},
                "sort": [{"timestamp.keyword": {"order": "desc"}}],
                "size": size
            }
            # Use body=query for v8 compatibility if needed, or just pass kwargs
            response: Any = self.client.search(index=self.index, **query)
            events: List[Dict[str, Any]] = []
            
            hits = response.get('hits', {}).get('hits', [])
            for hit in hits:
                source: Dict[str, Any] = hit.get('_source', {})
                
                # Transform locations dictionary to list of objects
                raw_locations: Dict[str, Dict[str, str]] = source.get('locations', {})
                formatted_locations: List[Dict[str, str]] = [
                    {"name": name, "lat": str(loc.get("lat", "")), "lon": str(loc.get("lon", ""))}
                    for name, loc in raw_locations.items()
                ]
                
                event: Dict[str, Any] = {
                    "raw_message": source.get("raw_message", ""),
                    "summary": source.get("summary", ""),
                    "timestamp": source.get("timestamp", ""),
                    "locations": formatted_locations
                }
                logger.debug(f"Formatted event: {event}")
                events.append(event)
            return events
        except Exception as e:
            logger.error(f"Error fetching from Elasticsearch: {e}")
            return []

_es_instance: Optional[ESClient] = None

def get_es_client() -> ESClient:
    global _es_instance
    if _es_instance is None:
        _es_instance = ESClient()
    return _es_instance
