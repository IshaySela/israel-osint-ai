from datetime import datetime, timezone, timedelta
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
                "sort": [{"timestamp_epoch": {"order": "desc"}}],
                "size": size
            }
            response: Any = self.client.search(index=self.index, **query)
            events: List[Dict[str, Any]] = []

            hits = response.get('hits', {}).get('hits', [])
            for hit in hits:
                source: Dict[str, Any] = hit.get('_source', {})
                events.append({
                    "raw_message": source.get("raw_message", ""),
                    "summary": source.get("summary", ""),
                    "timestamp_epoch": source.get("timestamp_epoch", 0),
                    "locations": source.get("locations", [])
                })
            logger.info(f"Retrived {len(events)} events from the database")
            return events
        except Exception as e:
            logger.error(f"Error fetching from Elasticsearch: {e}")
            return []

    def get_events_in_range(self, from_minutes_ago: int, to_minutes_ago: int, size: int = 200) -> List[Dict[str, Any]]:
        now = datetime.now(timezone.utc)
        range_from = int((now - timedelta(minutes=from_minutes_ago)).timestamp())
        range_to   = int((now - timedelta(minutes=to_minutes_ago)).timestamp())
        try:
            query: Dict[str, Any] = {
                "query": {"range": {"timestamp_epoch": {"gte": range_from, "lte": range_to}}},
                "sort": [{"timestamp_epoch": {"order": "desc"}}],
                "size": size
            }

            response: Any = self.client.search(index=self.index, **query)
            events: List[Dict[str, Any]] = []
            hits = response.get('hits', {}).get('hits', [])
            for hit in hits:
                source: Dict[str, Any] = hit.get('_source', {})
                events.append({
                    "raw_message": source.get("raw_message", ""),
                    "summary": source.get("summary", ""),
                    "timestamp_epoch": source.get("timestamp_epoch", 0),
                    "locations": source.get("locations", [])
                })
            logger.info(f"Retrieved {len(events)} events in range [{range_from} → {range_to}]")
            return events
        except Exception as e:
            logger.error(f"Error fetching events by range from Elasticsearch: {e}")
            return []

_es_instance: Optional[ESClient] = None

def get_es_client() -> ESClient:
    global _es_instance
    if _es_instance is None:
        _es_instance = ESClient()
    return _es_instance
