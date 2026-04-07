import json
import os
from pathlib import Path
from threading import Lock
from typing import List, Optional, Type
from dotenv import load_dotenv

class Config:
    _instance: Optional['Config'] = None
    _lock: Lock = Lock()

    elasticsearch_urls: List[str]
    elasticsearch_index: str
    port: int
    host: str
    debug: bool
    rabbitmq_host: str
    processed_events_exchange: str

    def __new__(cls: Type['Config']) -> 'Config':
        with cls._lock:
            if cls._instance is None:
                cls._instance = super(Config, cls).__new__(cls)
                load_dotenv()
                shared = _load_shared_config()
                cls._instance.elasticsearch_urls = os.getenv("ELASTICSEARCH_URLS", "http://localhost:9200").split(",")
                cls._instance.elasticsearch_index = os.getenv("ELASTICSEARCH_INDEX", shared.get("elasticsearch", {}).get("index", "osint_events"))
                cls._instance.port = int(os.getenv("PORT", "5000"))
                cls._instance.host = os.getenv("HOST", "127.0.0.1")
                cls._instance.debug = os.getenv("DEBUG", "False").lower() in ("true", "1", "t")
                cls._instance.rabbitmq_host = os.getenv("RABBITMQ_HOST", "localhost")
                cls._instance.processed_events_exchange = os.getenv("PROCESSED_EVENTS_EXCHANGE", shared.get("messaging", {}).get("processed_events_exchange", "processed_events"))
        return cls._instance

def _load_shared_config() -> dict:
    config_path = Path("/shared/config/config.json")
    if config_path.exists():
        with open(config_path) as f:
            return json.load(f)
    return {}

def get_config() -> Config:
    return Config()
