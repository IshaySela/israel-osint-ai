import json
import asyncio
import aio_pika
from aio_pika.exceptions import AMQPConnectionError, AMQPChannelError
from typing import Dict, Any, Callable, Awaitable
from loguru import logger


class MessageBroker:
    """The class MessageBroker abstract the impl details, optimzations etc. for directly working with
    the message broker and provides a simple declartive API specific to this project.
    """
    def __init__(
        self,
        rabbit_host: str,
        rabbit_queue: str = "osint_events",
        raw_events_exchange: str = "raw_events",
        processed_events_exchange: str = "processed_events",
        dlx_exchange: str = "dead_letter",
        dlx_queue: str = "dead_letter_queue",
    ) -> None:
        self.rabbit_host = rabbit_host
        self.rabbit_queue = rabbit_queue
        self.connection: aio_pika.abc.AbstractRobustConnection | None = None
        self.is_connected = False
        self.raw_events_exchange = raw_events_exchange
        self.processed_events_exchange = processed_events_exchange
        self.dlx_exchange = dlx_exchange
        self.dlx_queue = dlx_queue

    async def connect_async(self) -> None:
        """Establishes a robust async connection to RabbitMQ and declares all exchanges and queues."""
        self.connection = await aio_pika.connect_robust(host=self.rabbit_host)
        self.is_connected = True

        async with self.connection.channel() as channel:
            # DLX exchange + queue
            dlx = await channel.declare_exchange(self.dlx_exchange, aio_pika.ExchangeType.FANOUT, durable=True)
            dlx_q = await channel.declare_queue(self.dlx_queue, durable=True)
            await dlx_q.bind(dlx, routing_key="#")

            # Raw events exchange + queue (with DLX)
            raw = await channel.declare_exchange(self.raw_events_exchange, aio_pika.ExchangeType.FANOUT, durable=True)
            raw_q = await channel.declare_queue(
                self.rabbit_queue,
                durable=True,
                arguments={"x-dead-letter-exchange": self.dlx_exchange},
            )
            await raw_q.bind(raw, routing_key="")

            # Processed events exchange
            await channel.declare_exchange(self.processed_events_exchange, aio_pika.ExchangeType.FANOUT, durable=True)

    async def publish_raw_event_async(self, event_data: Dict[str, Any]) -> None:
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
            exchange = await channel.get_exchange(self.raw_events_exchange)
            msg = aio_pika.Message(json.dumps(event_data).encode())
            await exchange.publish(msg, routing_key="")
            

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
            exchange = await channel.get_exchange(self.raw_events_exchange)
            await self._listen_async(exchange, callback, queue_name="")

    async def listen_processed_events_async(self,
                                            callback: Callable[[Dict[str, Any]], Awaitable[None]]) -> None:
        if self.connection is None:
            raise RuntimeError("Must call .connect_async before trying to listen")

        channel = await self.connection.channel()
        exchange = await channel.get_exchange(self.processed_events_exchange)
        await self._listen_async(exchange, callback)