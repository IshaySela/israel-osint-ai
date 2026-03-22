import asyncio
from typing import List
from services.Configuration import TelegramScraperConfig
from pyrogram import filters
from pyrogram.client import Client
from pyrogram.methods.utilities.idle import idle
from pyrogram.types import Message
from services.ClassifyTelegramMessage import classify_telegram_msg
from services.MessageBroker import MessageBroker
from services.Logger import setup_logging
from loguru import logger

setup_logging()

CONFIG = TelegramScraperConfig.get()
MONITORED_CHANNEL_IDS: List[str | int] = [c.channelId for c in CONFIG.channels]

broker = MessageBroker(CONFIG.rabbit_host, CONFIG.rabbit_queue)

async def main():
    client = Client("israel-osint-ai-telegram", CONFIG.api_id, CONFIG.api_hash)
  
    @client.on_message(filters.channel & filters.chat(MONITORED_CHANNEL_IDS))
    async def process_messages(client: Client, message: Message) -> None:
        """Process messages from monitored channels."""
        text = f"{message.caption or ''} {message.text or ''}".strip()
        if not text:
            return

        logger.info(f"Recived Message from channel {message.chat.title} {message.chat.id}")

        event_type = await classify_telegram_msg(text)
        logger.info(f"Classified Message from channel{message.chat.id}: {text[:30]}... | Type: {event_type}")
        
        if event_type != 'not_relevant':
            event_data = {
                'text': text,
                'event_type': event_type,
                'chat_id': message.chat.id,
                'message_id': message.id,
                'date': str(message.date)
            }
            
            broker.publish_event(event_data)
            logger.info(f"Published event: {event_type} {text[:30]}")
            
            
    await client.start()
    
    logger.info("Initializing channel cache...")
    async for dialog in client.get_dialogs(limit=100): # type: ignore
        if dialog.chat.id in MONITORED_CHANNEL_IDS:
            logger.info(f"Init channel: {dialog.chat.id} {dialog.chat.title}") 
    
    logger.info("Telegram Scraper is now online and listening.")
    
    await idle() 
    
    await client.stop()


if __name__ == "__main__":
    channelsInfo = [{c.channelName, c.channelId} for c in CONFIG.channels]
    logger.info(f"Starting Telegram Scraper, listening on channels: {channelsInfo}")
    try:
        asyncio.run(main())
    except KeyboardInterrupt:
        logger.info("Scraper stopped by user")
