import os
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

    def __new__(cls: Type['Config']) -> 'Config':
        with cls._lock:
            if cls._instance is None:
                cls._instance = super(Config, cls).__new__(cls)
                load_dotenv()
                cls._instance.elasticsearch_urls = os.getenv("ELASTICSEARCH_URLS", "http://localhost:9200").split(",")
                cls._instance.elasticsearch_index = os.getenv("ELASTICSEARCH_INDEX", "osint_events")
                cls._instance.port = int(os.getenv("PORT", "5000"))
                cls._instance.host = os.getenv("HOST", "127.0.0.1")
                cls._instance.debug = os.getenv("DEBUG", "False").lower() in ("true", "1", "t")
        return cls._instance

def get_config() -> Config:
    return Config()
