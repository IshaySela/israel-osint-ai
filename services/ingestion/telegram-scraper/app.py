import asyncio
from typing import List
from services.Logger import setup_logging
from loguru import logger
from services.Configuration import TelegramScraperConfig
from telethon import TelegramClient, events
from telethon.types import Message, Chat
from services.MessageBroker import MessageBroker
from services.ClassifyTelegramMessage import classify_telegram_msg

setup_logging()

CONFIG = TelegramScraperConfig.get()
MONITORED_CHANNEL_IDS: List[int] = [c.channelId for c in CONFIG.channels]


client = TelegramClient('israel-osint-ai-telegram', int(CONFIG.api_id), CONFIG.api_hash)
broker = MessageBroker(CONFIG.rabbit_host, CONFIG.rabbit_queue)


@client.on(events.NewMessage(chats=MONITORED_CHANNEL_IDS))
async def handler(event: events.NewMessage.Event):
    msg: Message = event.message
    text = msg.message or ''
    
    chat: Chat = await event.get_chat() # type: ignore
    
    if text is None:
        logger.error(f"Skipping, recived empty msg from channel {chat.id} {chat.title}")
    
    logger.error(f"Recived msg from channel {chat.id} {chat.title} {text[:30]}")
    
    event_type = await classify_telegram_msg(text)
    logger.info(f"Classified Message from channel{chat.id}: {text[:30]}... | Type: {event_type}")
        
    if event_type != 'not_relevant':
        event_data = {
            'text': text,
            'event_type': event_type,
            'chat_id': chat.id,
            'message_id': msg.id,
            'date': str(msg.date)
        }
        
        broker.publish_event(event_data)
        logger.info(f"Published event: {event_type} {text[:30]}")
    

async def main():
    logger.info(f"Starting telegram scraper service for channels: {[c.channelName for c in CONFIG.channels]}")
    client.start()
    await client.connect()
    
    await client.run_until_disconnected() # type: ignore
   
if __name__ == "__main__":
    asyncio.run(main())