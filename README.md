# Israel OSINT AI

This project aims to ingest, process, and visualize OSINT data from various sources (Telegram, RSS, Web Scraping) in a unified map-based interface.

## Infrastructure Setup

To start the infrastructure required for local development (Elasticsearch and RabbitMQ), follow these steps:

### Prerequisites
- [Docker](https://www.docker.com/get-started)
- [Docker Compose](https://docs.docker.com/compose/install/)

### Start the Services
Run the following command from the root of the project:
```bash
docker compose up -d
```

### Accessing the Services
- **Elasticsearch**: `http://localhost:9200`
- **Kibana**: `http://localhost:5601`
- **RabbitMQ Management UI**: `http://localhost:15672` (Username: `guest`, Password: `guest`)
- **RabbitMQ AMQP Broker**: `localhost:5672`

### Stop the Services
To stop and remove the containers:
```bash
docker compose down
```

## Project Architecture
Refer to [ARCHITECTURE.md](ARCHITECTURE.md) for more details.
