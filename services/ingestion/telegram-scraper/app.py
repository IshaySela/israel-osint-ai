from services.Configuration import TelegramScraperConfig
import asyncio
from pyrogram import filters
from pyrogram.client import Client
from pyrogram.types import Message
from services.ClassifyTelegramMessage import classify_telegram_msg


CONFIG = TelegramScraperConfig.get()
LIVE_TEST_CHANNEL = -1001613161072
TZOFAR_TEST_CHANNEL = -1001436772127
MY_TEST_CHANNEL = -1003756841569
client = Client("israel-osint-ai-telegram", CONFIG.api_id, CONFIG.api_hash)

@client.on_message(filters.channel & filters.chat([LIVE_TEST_CHANNEL, TZOFAR_TEST_CHANNEL, MY_TEST_CHANNEL]))
async def listen_for_messages(client: Client, message: Message):
    text = f"{message.caption or ''} {message.text or ''}"
    event_type = await classify_telegram_msg(text)
    print('Recived msg', { text, event_type })
    
client.run()