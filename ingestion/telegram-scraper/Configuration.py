
from dataclasses import dataclass
import os

@dataclass
class TelegramScraperConfig:
    def __init__(self, api_id: str, api_hash: str):
        self.api_id = api_id
        self.api_hash = api_hash
    
    @staticmethod
    def load_config_from_env() -> 'TelegramScraperConfig':
        """Loads the Telegram API configuration from environment variables.

        Raises:
            ValueError: If one of the environment variables is not set

        Returns:
            TelegramScraperConfig: The loaded config from env
        """
        api_id = os.environ.get('TELEGRAM_API_ID')
        api_hash = os.environ.get('TELEGRAM_API_HASH')
        
        if api_id is None or api_hash is None:
            raise ValueError('TELEGRAM_API_ID and TELEGRAM_API_HASH must be set')
        
        return TelegramScraperConfig(api_id=api_id, api_hash=api_hash)
