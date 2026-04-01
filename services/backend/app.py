import os
from flask import Flask, request, jsonify, Response
from flask_cors import CORS
from ariadne import load_schema_from_path, make_executable_schema, graphql_sync, QueryType, EnumType
from ariadne.explorer import ExplorerGraphiQL
from typing import Any, Dict, List, Tuple, Union, Optional
from elasticsearch_client import get_es_client, ESClient
from config import get_config, Config
from loguru import logger
from flask_sse import sse
from shared.MessageBroker import MessageBroker
import asyncio
import threading
import json
from models.ProcessedMessageEvent import ProcessedEventMessage

# Initialize Flask app
app: Flask = Flask(__name__)
app.register_blueprint(sse, url_prefix='/events-stream')
CORS(app)

# Load GraphQL schema with absolute path
BASE_DIR: str = os.path.dirname(os.path.abspath(__file__))
schema_path: str = os.path.join(BASE_DIR, "schema.graphql")
type_defs: str = load_schema_from_path(schema_path)
query: QueryType = QueryType()
cfg: Config = get_config()

broker = MessageBroker(rabbit_host=cfg.rabbitmq_host, rabbit_queue="")

async def publish_events_to_clients(msg: Dict[str, Any]) -> None:
    ev = ProcessedEventMessage.model_validate(msg)
    sse.publish(ev,type="new_event")
    

@query.field("latestEvents")
def resolve_latest_events(*_: Any) -> List[Dict[str, Any]]:
    es: ESClient = get_es_client()
    return es.get_latest_events(size=50)

schema: Any = make_executable_schema(type_defs, query)
explorer: ExplorerGraphiQL = ExplorerGraphiQL()


@app.route("/graphql", methods=["GET"])
def graphql_playground() -> Union[str, Tuple[str, int]]:
    return explorer.html(None), 200

@app.route("/graphql", methods=["POST"])
def graphql_server() -> Tuple[Response, int]:
    data: Optional[Dict[str, Any]] = request.get_json()
    success, result = graphql_sync(
        schema,
        data,
        context_value=request,
        debug=app.debug
    )
    status_code: int = 200 if success else 400
    return jsonify(result), status_code

async def main():
    logger.info(f"Starting BFF on {cfg.host}:{cfg.port} (debug={cfg.debug}), elasticsearch={cfg.elasticsearch_urls}")
    broker_task = asyncio.create_task(broker.listen_async(cfg.processed_events_exchange,"", publish_events_to_clients))
    
    flask_thread = threading.Thread(
        target=app.run, 
        kwargs={"host": cfg.host, "port": cfg.port, "debug": False}
    )
    
    flask_thread.start()
    
    try:
        await asyncio.gather(broker_task)
    except asyncio.CancelledError:
        logger.info("Shutting down...")
    finally:
        broker_task.cancel()

if __name__ == "__main__":
    asyncio.run(main())

