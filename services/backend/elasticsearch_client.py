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
                "size": size
            }
            # The Elasticsearch client in v8+ has different return types for search
            response: Any = self.client.search(index=self.index, **query)
            events: List[Dict[str, Any]] = []
            
            hits = response.get('hits', {}).get('hits', [])
            for hit in hits:
                source: Dict[str, Any] = hit.get('_source', {})
                # Ensure chat_id is returned as a string for GraphQL compatibility
                if 'chat_id' in source:
                    source['chat_id'] = str(source['chat_id'])
                events.append(source)
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
