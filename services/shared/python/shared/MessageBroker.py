import pika
import json
import time
import asyncio
import aio_pika
from typing import Optional, Dict, Any, Callable, Awaitable
from loguru import logger
from pika.exceptions import AMQPChannelError, AMQPConnectionError
from pika.adapters.blocking_connection import BlockingChannel

class MessageBroker:
    connection: Optional[pika.BlockingConnection]
    channel: Optional[BlockingChannel]

    def __init__(self, rabbit_host: str, rabbit_queue: str, max_retries: int = 5, retry_delay: int = 5) -> None:
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

    async def listen_async(
        self,
        exchange: str,
        queue_name: str,
        callback: Callable[[Dict[str, Any]], Awaitable[None]]
    ) -> None:
        """
        Listens for messages from a specific exchange and queue asynchronously.
        
        Args:
            exchange: The name of the exchange to bind to.
            queue_name: The name of the queue to listen on.
            callback: An async function to call when a message is received.
        """
        logger.info(f"Starting async listener on exchange: '{exchange}', queue: '{queue_name}'")
        retries = 0
        while retries < self.max_retries:
            try:
                connection = await aio_pika.connect_robust(
                    f"amqp://{self.rabbit_host}/",
                )
                async with connection:
                    channel = await connection.channel()
                    
                    if exchange:
                        # Using TOPIC exchange type by default for flexibility
                        await channel.declare_exchange(exchange, aio_pika.ExchangeType.TOPIC, durable=True)
                    
                    queue = await channel.declare_queue(queue_name, durable=True)
                    
                    if exchange:
                        await queue.bind(exchange, routing_key="#")
                    
                    logger.info(f"Waiting for messages in {queue_name}...")
                    async with queue.iterator() as queue_iter:
                        async for message in queue_iter:
                            async with message.process():
                                try:
                                    event_data = json.loads(message.body.decode())
                                    await callback(event_data)
                                    await message.ack()
                                except json.JSONDecodeError:
                                    logger.error(f"Failed to decode message body: {message.body!r}")
                                    await message.nack(requeue=False)
                                except Exception as e:
                                    logger.error(f"Error in async callback: {e}")
                                    await message.nack(requeue=False)
            except Exception as e:
                retries += 1
                logger.warning(f"Async connection to RabbitMQ failed ({e}), retrying in {self.retry_delay}s ({retries}/{self.max_retries})...")
                await asyncio.sleep(self.retry_delay)
        
        if retries >= self.max_retries:
            logger.error(f"Failed to connect to RabbitMQ (async) after {self.max_retries} retries")
            raise RuntimeError(f"Failed to connect to RabbitMQ (async) after {self.max_retries} retries")
