import os
import unittest
from typing import Any, Dict, List
from flask import Flask
from flask.testing import FlaskClient
from app import app
from config import get_config, Config

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

    def test_graphql_endpoint_structure(self) -> None:
        """Test that the /graphql endpoint is reachable and returns the expected structure."""
        query: str = """
        query {
            latestEvents {
                raw_message
                summary
                timestamp
                locations {
                    name
                    lat
                    lon
                }
            }
        }
        """
        response: Any = self.client.post("/graphql", json={"query": query})
        self.assertEqual(response.status_code, 200)
        
        data: Dict[str, Any] = response.get_json()
        self.assertIn("data", data)
        self.assertIn("latestEvents", data["data"])
        
        events: List[Dict[str, Any]] = data["data"]["latestEvents"]
        self.assertIsInstance(events, list)
        
        if len(events) > 0:
            event: Dict[str, Any] = events[0]
            self.assertIn("raw_message", event)
            self.assertIn("summary", event)
            self.assertIn("timestamp", event)
            self.assertIn("locations", event)
            self.assertIsInstance(event["locations"], list)
            if len(event["locations"]) > 0:
                loc = event["locations"][0]
                self.assertIn("name", loc)
                self.assertIn("lat", loc)
                self.assertIn("lon", loc)

if __name__ == "__main__":
    unittest.main()
