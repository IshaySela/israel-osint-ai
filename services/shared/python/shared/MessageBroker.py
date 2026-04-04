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
            

    async def listen_async(self,
                           queu_name: str,
                           routing_key: str,
                           callback: Callable[[Dict[str, Any]],Awaitable[None]],
                           exchange_name: str) -> None:
        """Binds a queue to an exchange and consumes messages indefinitely.

        Declares the queue (auto-delete) and optionally binds it to a named
        exchange before entering a message loop. Each message body is
        JSON-decoded and forwarded to `callback`. Malformed JSON is logged and
        skipped; all other exceptions are logged and re-raised.

        Args:
            queu_name: Name of the queue to declare and consume from.
            routing_key: Routing key used when binding the queue to the exchange.
            callback: Async callable invoked with the parsed message dict.
            exchange_name: Exchange to bind to. Use '/' to skip exchange
                declaration and use the channel's default exchange.

        Raises:
            RuntimeError: If called before `connect_async`.
        """
        if self.connection is None:
            raise RuntimeError("Must call .connect before trying to publish")
        
        channel = await self.connection.channel()
        
        async with channel:
            queue = await channel.declare_queue(name=queu_name, auto_delete=True)
                
            exchange = await channel.declare_exchange(name=exchange_name,type="fanout")
            
            await queue.bind(exchange, routing_key=routing_key)
            
            async with queue.iterator() as queue_iter:
                async for msg in queue_iter:
                    try:
                        parsed = json.loads(msg.body)
                        await callback(parsed)
                    except json.JSONDecodeError as e:
                        logger.error(f"Failed to decode json message {e!r}")
                    except Exception as e:
                        logger.error(f"Unknown exception occured while parsing message from exchange {exchange_name} queue {queu_name} key: {routing_key} error: {e}")
                        raise e