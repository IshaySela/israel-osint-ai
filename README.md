# Israel OSINT AI

This project aims to ingest, process, and visualize OSINT data from various sources (Telegram, RSS, Web Scraping) in a unified map-based interface.

![Example Photo](https://github.com/IshaySela/israel-osint-ai/blob/master/static/example.jpg?raw=true)

## Why
The main intent of the project is to create a fully working product using various technologies from end to end.

## How
Multiple microservices that activley ingest OSINT sources. The microservice utilize the OpenAI API in order to filter 
events by relevance, communicate via RabbitMQ and display the events on a map using React client.

## Development Startup

To start the infrastructure required for local development follow these steps:

1. Make sure docker is installed and running.
2. Generate Telegram app id and hash via [telegram API development tools](https://my.telegram.org/apps)
3. Generate OpenAI api key using [OpenAI API platform](https://platform.openai.com/api-keys)
4. Configure the ```APP_ID```, ```APP_HASH``` and ```OPENAI_API_KEY``` environment variables.

Run the following command:
```bash
docker compose up -d
```
Access the frontend via ```http://localhost:5173/```

### Stop the Services
To stop and remove the containers:
```bash
docker compose down
```

## Project Architecture
The project is comprised of 4 layers:

- **Ingestion Layer** - Ingest data form various sources and publish to RabbitMQ
- **Processing Layer** - Listen to raw events from the queue and extract info (geocode, AI sumamry, etc.) and index to ES
- **Backend for Frontend** - Provide API for the frontend to view the data
- **Frontend** - User interface for the end user.
 
Refer to [ARCHITECTURE.md](ARCHITECTURE.md) for more details.

## Implemented Features
### Sprint 1: Core Implementation (Walking Skeleton) - Done
Basic implementation of the services, ensure that data flows correctly and visualized on the map, startup via docker compose.
- Telegram scraper ingestion ([#1](https://github.com/IshaySela/israel-osint-ai/issues/1))
- Minimal normalization/processing service ([#3](https://github.com/IshaySela/israel-osint-ai/issues/3))
- Basic map display frontend ([#4](https://github.com/IshaySela/israel-osint-ai/issues/4))
- Core infrastructure services ([#2](https://github.com/IshaySela/israel-osint-ai/issues/2))

### Sprint 2: Performance & Robustness - Done
- Persistent geocode result caching with Redis ([#14](https://github.com/IshaySela/israel-osint-ai/issues/14))
- Events filtering & restrict geocoder to Israel ([#8](https://github.com/IshaySela/israel-osint-ai/issues/8))
- Worker pool pattern in the processing service ([#19](https://github.com/IshaySela/israel-osint-ai/issues/19))
- Fixed auto-ack data loss in message broker ([#13](https://github.com/IshaySela/israel-osint-ai/issues/13))
- Critical bug fixes in ingestion error handling ([#20](https://github.com/IshaySela/israel-osint-ai/issues/20))

### Sprint 3: Reliability & Integrations - Done
- Improved RabbitMQ configuration with exchanges & DLX ([#22](https://github.com/IshaySela/israel-osint-ai/issues/22))
- Processing notifications via dedicated exchange ([#21](https://github.com/IshaySela/israel-osint-ai/issues/21))
- Shared configuration across all services ([#23](https://github.com/IshaySela/israel-osint-ai/issues/23))
- SSE notifications between React client and backend ([#26](https://github.com/IshaySela/israel-osint-ai/issues/26))
- GraphQL query for events in the last 24 hours ([#30](https://github.com/IshaySela/israel-osint-ai/issues/30))
- Source metadata from telegram scraper ([#28](https://github.com/IshaySela/israel-osint-ai/issues/28))

## Roadmap
### Features
- Add RSS news feed ingestion service ([#37](https://github.com/IshaySela/israel-osint-ai/issues/37))
- Support images & videos from telegram scraper ([#31](https://github.com/IshaySela/israel-osint-ai/issues/31))
- Geospatial search — GraphQL query by point, radius & time frame ([#39](https://github.com/IshaySela/israel-osint-ai/issues/39))
- Support locally deployed LLM via Ollama for development ([#25](https://github.com/IshaySela/israel-osint-ai/issues/25))
- Structured logs streamed to Elasticsearch across all services ([#34](https://github.com/IshaySela/israel-osint-ai/issues/34))
- Event ID field for end-to-end event tracking ([#41](https://github.com/IshaySela/israel-osint-ai/issues/41))

### Bugs
- AI summary translates English events into Hebrew ([#38](https://github.com/IshaySela/israel-osint-ai/issues/38))
- Events with too-broad locations still indexed ([#32](https://github.com/IshaySela/israel-osint-ai/issues/32))
- GraphQL query fails on missing/malformed ES documents ([#35](https://github.com/IshaySela/israel-osint-ai/issues/35), [#36](https://github.com/IshaySela/israel-osint-ai/issues/36))

### Chores
- Replace deprecated `build.bin` with `build.entrypoint` in `.air.toml` ([#33](https://github.com/IshaySela/israel-osint-ai/issues/33))
- Add Pydantic model in backend `elasticsearch_client.py` ([#40](https://github.com/IshaySela/israel-osint-ai/issues/40))
