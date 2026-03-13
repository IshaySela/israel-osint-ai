
from dataclasses import dataclass
import os
from dotenv import load_dotenv


@dataclass
class TelegramScraperConfig:
    def __init__(self, api_id: str, api_hash: str, openai_api_key: str):
        self.api_id = api_id
        self.api_hash = api_hash
        self.openai_api_key = openai_api_key
    __configSingleton: 'None | TelegramScraperConfig' = None
    
    @staticmethod
    def load_from_env() -> 'TelegramScraperConfig':
        """Loads the Telegram API configuration from environment variables.

        Raises:
            ValueError: If one of the environment variables is not set

        Returns:
            TelegramScraperConfig: The loaded config from env
        """
        
        if TelegramScraperConfig.__configSingleton is not None:
            return TelegramScraperConfig.__configSingleton
        
        load_dotenv()
        openai_api_key = os.environ.get('OPENAI_API_KEY')
        telegram_api_id = os.environ.get('TELEGRAM_API_ID')
        telegram_api_hash = os.environ.get('TELEGRAM_API_HASH')
        
        if telegram_api_id is None or telegram_api_hash is None or openai_api_key is None:
            raise ValueError('TELEGRAM_API_ID and TELEGRAM_API_HASH must be set')
        
        return TelegramScraperConfig(api_id=telegram_api_id, api_hash=telegram_api_hash, openai_api_key=openai_api_key)
