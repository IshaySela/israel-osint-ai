from services.Configuration import TelegramScraperConfig
CONFIG = TelegramScraperConfig.get()

from telethon import TelegramClient, events

# Initialize client
client = TelegramClient('session_name', int(CONFIG.api_id), CONFIG.api_hash)

@client.on(events.NewMessage())
async def handler(event):
    print('New message received:', event.message.message)

async def main():
    await client.start()
    # Keep the client running to receive updates
    await client.run_until_disconnected()

client.loop.run_until_complete(main())
