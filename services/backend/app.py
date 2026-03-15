import os
from flask import Flask, request, jsonify, Response
from flask_cors import CORS
from ariadne import load_schema_from_path, make_executable_schema, graphql_sync, QueryType, EnumType
from ariadne.explorer import ExplorerGraphiQL
from typing import Any, Dict, List, Tuple, Union, Optional
from elasticsearch_client import get_es_client, ESClient
from config import get_config, Config
from loguru import logger
import sys

# Initialize Flask app
app: Flask = Flask(__name__)
CORS(app)

# Load GraphQL schema with absolute path
BASE_DIR: str = os.path.dirname(os.path.abspath(__file__))
schema_path: str = os.path.join(BASE_DIR, "schema.graphql")
type_defs: str = load_schema_from_path(schema_path)
query: QueryType = QueryType()

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

if __name__ == "__main__":
    cfg: Config = get_config()
    logger.info(f"Starting BFF on {cfg.host}:{cfg.port} (debug={cfg.debug}), elasticsearch={cfg.elasticsearch_urls}")
    app.run(host=cfg.host, port=cfg.port, debug=cfg.debug)
