
import json
from pathlib import Path
from dataclasses import dataclass
import os
from dotenv import load_dotenv
from typing import List


@dataclass
class TelegramScraperConfig:
    def __init__(self, api_id: str, api_hash: str, openai_api_key: str, rabbit_host: str, rabbit_queue: str, channels: List[int]):
        self.api_id = api_id
        self.api_hash = api_hash
        self.openai_api_key = openai_api_key
        self.rabbit_host = rabbit_host
        self.rabbit_queue = rabbit_queue
        self.channels = channels
    __configSingleton: 'None | TelegramScraperConfig' = None
    
    @staticmethod
    def get() -> 'TelegramScraperConfig':
        """Loads the Telegram API configuration from environment variables and channels from JSON.

        Raises:
            ValueError: If one of the environment variables is not set or channels.json is missing

        Returns:
            TelegramScraperConfig: The loaded config
        """
        
        if TelegramScraperConfig.__configSingleton is not None:
            return TelegramScraperConfig.__configSingleton
        
        load_dotenv()
        openai_api_key = os.environ.get('OPENAI_API_KEY')
        telegram_api_id = os.environ.get('TELEGRAM_API_ID')
        telegram_api_hash = os.environ.get('TELEGRAM_API_HASH')
        rabbit_host = os.environ.get('RABBIT_HOST', 'localhost')
        rabbit_queue = os.environ.get('RABBIT_QUEUE', 'events')
        
        if telegram_api_id is None or telegram_api_hash is None or openai_api_key is None:
            raise ValueError('TELEGRAM_API_ID, TELEGRAM_API_HASH and OPENAI_API_KEY must be set')

        channels_file = Path(__file__).parent.parent / "channels.json"
        if not channels_file.exists():
            raise ValueError(f"channels.json not found at {channels_file}")

        try:
            with open(channels_file, 'r') as f:
                channels_data = json.load(f)
                channels = [c['channelId'] for c in channels_data['channels']]
        except (json.JSONDecodeError, KeyError, TypeError) as e:
            raise ValueError(f"Error parsing channels.json: {e}")
        
        if not channels:
            raise ValueError('No channels found in channels.json')
        
        return TelegramScraperConfig(
            api_id=telegram_api_id, 
            api_hash=telegram_api_hash, 
            openai_api_key=openai_api_key,
            rabbit_host=rabbit_host,
            rabbit_queue=rabbit_queue,
            channels=channels
        )
