# Israel OSINT AI

This project aims to ingest, process, and visualize OSINT data from various sources (Telegram, RSS, Web Scraping) in a unified map-based interface.

## Why?
The main intent of the project is to create a fully working product using various technologies from end to end.

## How?
Multiple microservices that activley ingest OSINT sources. The micro-service utilize the OpenAI API in order to filter 
events by relevance.
The processed events are then stored in Elasticsearch database and served to the React app via dedicated service.
Refere to detailed architecture choices in ARCHITECTURE.md.

## Development Startup

To start the infrastructure required for local development follow these steps:

### Telegram API
Generate the app id, app hash via [telegram API development tools](https://my.telegram.org/apps) and setup the ```APP_ID``` and ```APP_HASH``` environment variables

### OpenAI API
Create OpenAI API key in the [OpenAI API platform](https://platform.openai.com/api-keys) and set the OPENAI_API_KEY environment variable

```bash
docker compose up -d
```
Access the frontend via 

### Stop the Services
To stop and remove the containers:
```bash
docker compose down
```

## Project Architecture
Refer to [ARCHITECTURE.md](ARCHITECTURE.md) for more details.
