import asyncio
from services.Configuration import TelegramScraperConfig
from pyrogram import filters
from pyrogram.client import Client
from pyrogram.types import Message
from services.ClassifyTelegramMessage import classify_telegram_msg
from services.MessageBroker import MessageBroker
from services.Logger import setup_logging
from loguru import logger

setup_logging()

CONFIG = TelegramScraperConfig.get()
MONITORED_CHANNEL_IDS = {c.channelId for c in CONFIG.channels}
client = Client("israel-osint-ai-telegram", CONFIG.api_id, CONFIG.api_hash)

# broker = MessageBroker(CONFIG.rabbit_host, CONFIG.rabbit_queue)

@client.on_message(filters.channel)
async def debug_messages(client: Client, message: Message) -> None:
    is_monitored = message.chat.id in MONITORED_CHANNEL_IDS
    logger.info(f"Received message from channel ID: {message.chat.id} (Monitored: {is_monitored})")
    
    if not is_monitored:
        return

    text = f"{message.caption or ''} {message.text or ''}"
    event_type = await classify_telegram_msg(text)
    logger.info(f"Received msg: {text[:100]}... | Type: {event_type}")
    
    if event_type != 'not_relevant':
        event_data = {
            'text': text,
            'event_type': event_type,
            'chat_id': message.chat.id,
            'message_id': message.id,
            'date': str(message.date)
        }
        # broker = MessageBroker(CONFIG.rabbit_host, CONFIG.rabbit_queue)
        # broker.publish_event(event_data)
        logger.info(f"Published event: {event_type}")

async def main() -> None:
    channel_names = [c.channelName for c in CONFIG.channels]
    logger.info(f"Starting Telegram Scraper, listening on channels: {channel_names}")
    async with client:
        await asyncio.Event().wait()

if __name__ == "__main__":
    try:
        asyncio.run(main())
    except KeyboardInterrupt:
        logger.info("Scraper stopped by user")
