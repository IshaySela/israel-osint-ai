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
    def __init__(self, rabbit_host: str, rabbit_queue: str, max_retries: int = 5, retry_delay: int = 5, timeout: int = 10) -> None:
        self.rabbit_host = rabbit_host
        self.rabbit_queue = rabbit_queue
        self.max_retries = max_retries
        self.retry_delay = retry_delay
        self.connection = None
        self.timeout = timeout

    async def connect(self):
        self.connection = await aio_pika.connect_robust(host=self.rabbit_host)
        await self.connection.connect(self.timeout)
        
        async with self.connection:
            channel = self.connection.channel()  
            await channel.declare_queue(name=self.rabbit_queue)
    async def publish_event_async(self,event_data: Dict[str, Any]) -> None:
        if self.connection is None:
            raise RuntimeError("Must call .connect before trying to publish")
            
        async with self.connection:
            channel = self.connection.channel()
            msg = aio_pika.Message(json.dumps(event_data).encode())
            await channel.default_exchange.publish(msg, routing_key=self.rabbit_queue)
            

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
                        await channel.declare_exchange(exchange, aio_pika.ExchangeType.FANOUT, durable=True)
                    
                    # For empty queue_name, use non-durable, auto-delete, exclusive
                    is_exclusive = not queue_name
                    queue = await channel.declare_queue(
                        queue_name, 
                        durable=not is_exclusive,
                        auto_delete=is_exclusive,
                        exclusive=is_exclusive
                    )
                    
                    if exchange:
                        await queue.bind(exchange, routing_key="#")
                    
                    logger.info(f"Waiting for messages in {queue.name}...")
                    async with queue.iterator() as queue_iter:
                        async for message in queue_iter:
                            async with message.process():
                                try:
                                    event_data = json.loads(message.body.decode())
                                    await callback(event_data)
                                except json.JSONDecodeError:
                                    logger.error(f"Failed to decode message body: {message.body!r}")
                                except Exception as e:
                                    logger.error(f"Error in async callback: {e}")
                                    # message.process() will handle nack if exception is raised
                                    raise
            except Exception as e:
                retries += 1
                logger.warning(f"Async connection to RabbitMQ failed ({e}), retrying in {self.retry_delay}s ({retries}/{self.max_retries})...")
                await asyncio.sleep(self.retry_delay)
        
        if retries >= self.max_retries:
            logger.error(f"Failed to connect to RabbitMQ (async) after {self.max_retries} retries")
            raise RuntimeError(f"Failed to connect to RabbitMQ (async) after {self.max_retries} retries")
