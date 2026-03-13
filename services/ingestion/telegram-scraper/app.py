from services.Configuration import TelegramScraperConfig
import asyncio
from pyrogram import filters
from pyrogram.client import Client
from pyrogram.types import Message

CONFIG = TelegramScraperConfig.get()
TEST_CHANNEL = -1001613161072

client = Client("israel-osint-ai-telegram", CONFIG.api_id, CONFIG.api_hash)

@client.on_message(filters.channel & filters.chat([TEST_CHANNEL,-1003756841569]))
async def listen_for_messages(client: Client, message: Message):
    print('Recived msg', message.text)
    
client.run()