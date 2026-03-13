from services.Configuration import TelegramScraperConfig
import asyncio
from pyrogram import filters
from pyrogram.client import Client
from pyrogram.types import Message
from services.ClassifyTelegramMessage import classify_telegram_msg
from services.MessageBroker import MessageBroker


CONFIG = TelegramScraperConfig.get()
LIVE_TEST_CHANNEL = -1001613161072
TZOFAR_TEST_CHANNEL = -1001436772127
MY_TEST_CHANNEL = -1003756841569
client = Client("israel-osint-ai-telegram", CONFIG.api_id, CONFIG.api_hash)

broker = MessageBroker(CONFIG.rabbit_host, CONFIG.rabbit_queue)

@client.on_message(filters.channel & filters.chat([LIVE_TEST_CHANNEL, TZOFAR_TEST_CHANNEL, MY_TEST_CHANNEL]))
async def listen_for_messages(client: Client, message: Message):
    text = f"{message.caption or ''} {message.text or ''}"
    event_type = await classify_telegram_msg(text)
    print('Received msg', {'text': text, 'event_type': event_type})
    
    if event_type != 'not_relevant':
        event_data = {
            'text': text,
            'event_type': event_type,
            'chat_id': message.chat.id,
            'message_id': message.id,
            'date': str(message.date)
        }
        broker.publish_event(event_data)
        print(f"Published event: {event_type}")
    
client.run()
