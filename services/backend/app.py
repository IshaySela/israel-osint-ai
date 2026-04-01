import os
import asyncio
import json
from fastapi import FastAPI, Request
from fastapi.middleware.cors import CORSMiddleware
from fastapi.responses import StreamingResponse
from ariadne.asgi import GraphQL
from ariadne import load_schema_from_path, make_executable_schema, QueryType
from typing import Any, Dict, List, Optional, AsyncGenerator
from elasticsearch_client import get_es_client, ESClient
from config import get_config, Config
from loguru import logger
from shared.MessageBroker import MessageBroker
from sse_starlette.sse import EventSourceResponse

app = FastAPI()

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

cfg: Config = get_config()

# Load GraphQL schema
BASE_DIR: str = os.path.dirname(os.path.abspath(__file__))
schema_path: str = os.path.join(BASE_DIR, "schema.graphql")
type_defs: str = load_schema_from_path(schema_path)
query: QueryType = QueryType()

@query.field("latestEvents")
def resolve_latest_events(*_: Any) -> List[Dict[str, Any]]:
    es: ESClient = get_es_client()
    return es.get_latest_events(size=50)

schema: Any = make_executable_schema(type_defs, query)
graphql_app = GraphQL(schema, debug=cfg.debug)

# Message Broker for SSE
broker = MessageBroker(rabbit_host=cfg.rabbitmq_host, rabbit_queue="")

@app.get("/events-stream")
async def events_stream(request: Request) -> EventSourceResponse:
    async def event_generator() -> AsyncGenerator[Dict[str, Any], None]:
        queue = asyncio.Queue()

        async def callback(event_data: Dict[str, Any]) -> None:
            await queue.put(event_data)

        # Start listening for msgs in the background
        listen_task = asyncio.create_task(
            broker.listen_async(
                exchange=cfg.processed_events_exchange,
                queue_name="",
                callback=callback
            )
        )

        try:
            while True:
                if await request.is_disconnected():
                    break
                
                try:
                    # Wait for a message with a timeout to check for disconnection
                    event_data = await asyncio.wait_for(queue.get(), timeout=1.0)
                    yield {
                        "data": json.dumps(event_data),
                        "event": "message"
                    }
                except asyncio.TimeoutError:
                    continue
        finally:
            listen_task.cancel()
            try:
                await listen_task
            except asyncio.CancelledError:
                pass

    return EventSourceResponse(event_generator())

@app.post("/graphql")
async def graphql_post(request: Request):
    return await graphql_app.handle_request(request)

@app.get("/graphql")
async def graphql_get(request: Request):
    return await graphql_app.handle_request(request)

if __name__ == "__main__":
    import uvicorn
    logger.info(f"Starting FastAPI BFF on {cfg.host}:{cfg.port} (debug={cfg.debug}), elasticsearch={cfg.elasticsearch_urls}")
    uvicorn.run(app, host=cfg.host, port=cfg.port)
