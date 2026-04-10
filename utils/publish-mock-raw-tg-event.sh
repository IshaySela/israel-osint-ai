#!/bin/bash

mock="{
  \"id\": \"tg_99999\",
  \"source\": \"telegram\",
  \"timestamp\": \"2026-04-11T12:00:00+00:00\",
  \"text\": \"דיווח: פיצוץ חזק נשמע בסמוך לצומת גלילות. כוחות ביטחון בדרך לאזור.\",
  \"event_type\": \"security_incident\",
  \"chat_id\": -1003756841569,
  \"channel_title\": \"test\",
  \"channel_main_lang\": \"he\",
  \"message_id\": 99999
}"

if command -v amqp-publish > /dev/null 2>&1; then
    amqp-publish \
    --url="amqp://guest:guest@localhost:5672" \
    --exchange="raw_events" --routing-key="#" \
    --body="$mock"
else
    echo "Please install amqp-tools (amqp-publish is missing) apt install amqp-tools"
fi