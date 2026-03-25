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
Refer to [ARCHITECTURE.md](ARCHITECTURE.md) for more details.

## Roadmap
- SSE between the backend and the clients
- Image & Video of events from telegram scraper
- Displaying events in a polygon over a certain area
- More scrapers types (RSS, x.com etc.)
- Analyzing events over time

## Implemented Features
### Sprint 1: Core Implementation (Walking Skeleton) - Done.
Basic implementation of the services, ensure that data flows correctly and visualized on the map, startup via docker compose.

### Sprint 2: Performance & Robustness - Done
- Persistent geocode result caching
- Events filtering
- Implement the worker pool pattern in the processing service
- Bug fixes
