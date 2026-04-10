import os
from threading import Lock
from typing import List, Optional, Type
from dotenv import load_dotenv

from shared.config import SharedConfig, Topology


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
                topo = Topology.load()
                shared = SharedConfig.load()

                # Infrastructure — topology.json
                cls._instance.elasticsearch_urls = [topo.elasticsearch.url]
                cls._instance.rabbitmq_host = topo.rabbitmq.host

                # Secrets / service-specific — env var only
                cls._instance.port = int(os.getenv("PORT", "5000"))
                cls._instance.host = os.getenv("HOST", "127.0.0.1")
                cls._instance.debug = os.getenv("DEBUG", "False").lower() in ("true", "1", "t")

                # Topology — config.json only
                cls._instance.elasticsearch_index = shared.elasticsearch.index
                cls._instance.processed_events_exchange = shared.messaging.processed_events_exchange
        return cls._instance


def get_config() -> Config:
    return Config()