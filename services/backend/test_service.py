import os
import unittest
from typing import Any, Dict, List, cast
from flask import Flask
from flask.testing import FlaskClient
from app import app
from config import get_config, Config
from elasticsearch_client import get_es_client, ESClient

class TestBackendService(unittest.TestCase):
    def setUp(self) -> None:
        self.app: Flask = app
        self.app.testing = True
        self.client: FlaskClient = self.app.test_client()

    def test_config_singleton(self) -> None:
        """Test that the configuration is a singleton and loads values."""
        cfg1: Config = get_config()
        cfg2: Config = get_config()
        
        self.assertIs(cfg1, cfg2, "Config should be a singleton instance")
        self.assertIsInstance(cfg1.port, int)
        self.assertIsInstance(cfg1.elasticsearch_urls, list)

    def test_graphql_endpoint_structure(self) -> None:
        """Test that the /graphql endpoint is reachable and returns the expected structure."""
        query: str = """
        query {
            latestEvents {
                text
                event_type
                chat_id
                message_id
                date
            }
        }
        """
        response: Any = self.client.post("/graphql", json={"query": query})
        self.assertEqual(response.status_code, 200)
        
        data: Dict[str, Any] = response.get_json()
        self.assertIn("data", data)
        self.assertIn("latestEvents", data["data"])
        
        # latestEvents should be a list (even if empty)
        events: List[Any] = data["data"]["latestEvents"]
        self.assertIsInstance(events, list)
        
        if len(events) > 0:
            event: Dict[str, Any] = events[0]
            expected_keys: List[str] = ["text", "event_type", "chat_id", "message_id", "date"]
            for key in expected_keys:
                self.assertIn(key, event)

if __name__ == "__main__":
    unittest.main()
