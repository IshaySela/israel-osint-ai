import json
import asyncio
import aio_pika
from aio_pika.exceptions import AMQPConnectionError, AMQPChannelError
from typing import Dict, Any, Callable, Awaitable
from loguru import logger
from aio_pika.robust_connection import RobustConnection

class MessageBroker:
    """The class MessageBroker abstract the impl details, optimzations etc. for directly working with
    the message broker and provides a simple declartive API specific to this project.
    """
    def __init__(self, rabbit_host: str, rabbit_queue: str) -> None:
        self.rabbit_host = rabbit_host
        self.rabbit_queue = rabbit_queue
        self.connection = None
        self.is_connected = False

    async def connect_async(self):
        """Establishes a robust async connection to RabbitMQ."""
        self.connection = await aio_pika.connect_robust(host=self.rabbit_host)
        self.is_connected = True

    async def publish_event_async(self, event_data: Dict[str, Any]) -> None:
        """Publishes a JSON-serialized event to the default queue.

        Args:
            event_data: Dictionary to serialize and publish as a message.

        Raises:
            RuntimeError: If called before `connect_async`.
        """
        if self.connection is None:
            raise RuntimeError("Must call .connect before trying to publish")

        channel = await self.connection.channel()
        async with channel:
            msg = aio_pika.Message(json.dumps(event_data).encode())
            await channel.default_exchange.publish(msg, routing_key=self.rabbit_queue)
            

    async def _listen_async(self,
                            exchange: aio_pika.abc.AbstractExchange,
                            callback: Callable[[Dict[str, Any]], Awaitable[None]],
                            queue_name: str = "") -> None:
        if self.connection is None:
            raise RuntimeError("Must call .connect_async before trying to listen")

        channel = await self.connection.channel()

        async with channel:
            queue = await channel.declare_queue(name=queue_name, auto_delete=True)
            await queue.bind(exchange, routing_key="#")

            async with queue.iterator() as queue_iter:
                async for msg in queue_iter:
                    try:
                        parsed = json.loads(msg.body)
                        await callback(parsed)
                    except json.JSONDecodeError as e:
                        logger.error(f"Failed to decode json message {e!r}")
                    except Exception as e:
                        logger.error(f"Unknown exception occured while parsing message: {e}")
                        raise

    async def listen_raw_events_async(self,
                                     callback: Callable[[Dict[str, Any]], Awaitable[None]]) -> None:
        if self.connection is None:
            raise RuntimeError("Must call .connect_async before trying to listen")

        channel = await self.connection.channel()
        async with channel:
            await self._listen_async(channel.default_exchange, callback, queue_name=self.rabbit_queue)

    async def listen_processed_events_async(self,
                                            exchange_name: str,
                                            callback: Callable[[Dict[str, Any]], Awaitable[None]]) -> None:
        if self.connection is None:
            raise RuntimeError("Must call .connect_async before trying to listen")

        channel = await self.connection.channel()
        exchange = await channel.declare_exchange(name=exchange_name, type="fanout", durable=True)
        await self._listen_async(exchange, callback)