import json
from dataclasses import dataclass
from pathlib import Path


@dataclass
class RabbitMQTopology:
    host: str
    port: int
    management_port: int
    user: str
    password: str

    @classmethod
    def from_dict(cls, d: dict) -> 'RabbitMQTopology':
        return cls(
            host=d["host"],
            port=d["port"],
            management_port=d["management_port"],
            user=d["user"],
            password=d["password"],
        )


@dataclass
class ElasticsearchTopology:
    host: str
    port: int

    @classmethod
    def from_dict(cls, d: dict) -> 'ElasticsearchTopology':
        return cls(host=d["host"], port=d["port"])

    @property
    def url(self) -> str:
        return f"http://{self.host}:{self.port}"


@dataclass
class Topology:
    rabbitmq: RabbitMQTopology
    elasticsearch: ElasticsearchTopology

    @classmethod
    def load(cls, topo_path: str = "/shared/config/topology.json") -> 'Topology':
        path = Path(topo_path)
        if not path.exists():
            raise FileNotFoundError(f"topology.json not found at {path}")
        with open(path) as f:
            d = json.load(f)
        return cls(
            rabbitmq=RabbitMQTopology.from_dict(d["rabbitmq"]),
            elasticsearch=ElasticsearchTopology.from_dict(d["elasticsearch"]),
        )


@dataclass
class MessagingConfig:
    queue: str
    raw_events_exchange: str
    processed_events_exchange: str
    dlx_exchange: str
    dlx_queue: str

    @classmethod
    def from_dict(cls, d: dict) -> 'MessagingConfig':
        return cls(
            queue=d["queue"],
            raw_events_exchange=d["raw_events_exchange"],
            processed_events_exchange=d["processed_events_exchange"],
            dlx_exchange=d["dlx_exchange"],
            dlx_queue=d["dlx_queue"],
        )


@dataclass
class ElasticsearchConfig:
    index: str
    geocode_index: str

    @classmethod
    def from_dict(cls, d: dict) -> 'ElasticsearchConfig':
        return cls(index=d["index"], geocode_index=d["geocode_index"])


@dataclass
class OpenAIConfig:
    model: str

    @classmethod
    def from_dict(cls, d: dict) -> 'OpenAIConfig':
        return cls(model=d["model"])


@dataclass
class SharedConfig:
    messaging: MessagingConfig
    elasticsearch: ElasticsearchConfig
    openai: OpenAIConfig

    @classmethod
    def load(cls, cfg_path: str = "/shared/config/config.json") -> 'SharedConfig':
        path = Path(cfg_path)
        if not path.exists():
            raise FileNotFoundError(f"config.json not found at {path}")
        with open(path) as f:
            d = json.load(f)
        return cls(
            messaging=MessagingConfig.from_dict(d["messaging"]),
            elasticsearch=ElasticsearchConfig.from_dict(d["elasticsearch"]),
            openai=OpenAIConfig.from_dict(d["openai"]),
        )