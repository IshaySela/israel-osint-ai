from Configuration import TelegramScraperConfig
import asyncio
from pyrogram import filters
from pyrogram.client import Client
from pyrogram.types import Message

CONFIG = TelegramScraperConfig.load_from_env()
TEST_CHANNEL = -1001613161072

client = Client("israel-osint-ai-telegram", CONFIG.api_id, CONFIG.api_hash)

@client.on_message(filters.chat(TEST_CHANNEL))
async def listen_for_messages(client: Client, message: Message):
    print('Recived msg', message.text)
    
client.run()