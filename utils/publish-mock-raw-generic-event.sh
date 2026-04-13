#!/bin/bash

mock="{
  \"app_ev_id\": \"generic_99999\",
  \"source\": \"twitter\",
  \"timestamp\": \"2026-04-13T12:00:00+00:00\",
  \"raw_message\": \"דיווח: פיצוץ חזק נשמע בסמוך לצומת גלילות. כוחות ביטחון בדרך לאזור.\",
  \"data\": {}
}"

if command -v amqp-publish > /dev/null 2>&1; then
    amqp-publish \
    --url="amqp://guest:guest@localhost:5672" \
    --exchange="raw_events" --routing-key="#" \
    --body="$mock"
else
    echo "Please install amqp-tools (amqp-publish is missing) apt install amqp-tools"
fi
