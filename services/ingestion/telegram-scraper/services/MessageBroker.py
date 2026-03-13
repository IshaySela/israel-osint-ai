import pika
import json
import time
from typing import Optional, Dict, Any
from loguru import logger
from services.Configuration import TelegramScraperConfig
from pika.exceptions import AMQPChannelError, AMQPConnectionError
from pika.adapters.blocking_connection import BlockingChannel

class MessageBroker:
    connection: Optional[pika.BlockingConnection]
    channel: Optional[BlockingChannel]

    def __init__(self,rabbit_host: str, rabbit_queue: str, max_retries: int = 5, retry_delay: int = 5) -> None:
        self.rabbit_host = rabbit_host
        self.rabbit_queue = rabbit_queue
        self.max_retries = max_retries
        self.retry_delay = retry_delay
        self.channel = None
        self._connect()

    def _connect(self) -> None:
        retries = 0
        connected = False
        while retries < self.max_retries and not connected:
            retries += 1
            try:
                self.connection = pika.BlockingConnection(
                    pika.ConnectionParameters(host=self.rabbit_host)
                )
                self._setup_channel()
                logger.info(f"Connected to RabbitMQ at {self.rabbit_host}")
                connected = True
            except AMQPConnectionError:
                logger.warning(f"Failed to connect to RabbitMQ at {self.rabbit_host}, retrying in {self.retry_delay}s...")
                time.sleep(self.retry_delay)
        if not connected:
            raise RuntimeError(f"Failed to connect to RabbitMQ after {self.max_retries} retries")        
        
    def _setup_channel(self):
        if self.connection:
            self.channel = self.connection.channel()
        if self.channel:
            self.channel.queue_declare(queue=self.rabbit_queue)
    def publish_event(self, event_data: Dict[str, Any]) -> None:
        try:
            if self.channel:
                self.channel.basic_publish(
                    exchange='',
                    routing_key=self.rabbit_queue,
                    body=json.dumps(event_data)
                )
        except (AMQPConnectionError, AMQPChannelError):
            logger.error("RabbitMQ connection lost, reconnecting...")
            self._connect()
            if self.channel:
                self.channel.basic_publish(
                    exchange='',
                    routing_key=self.rabbit_queue,
                    body=json.dumps(event_data)
                )
