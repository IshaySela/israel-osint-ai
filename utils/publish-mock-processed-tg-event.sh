#!/bin/bash

mock="{
  \"dbId\": \"abc123xyz\",
  \"source\": \"telegram\",
  \"summary\": \"נפילת טיל בסמוך לצומת גלילות. כוחות ביטחון בדרך לאזור.\",
  \"timestamp_epoch\": 1744372800,
  \"channel_title\": \"mock\",
  \"channel_main_lang\": \"he\",
  \"locations\": [
    {
      \"name\": \"צומת גלילות\",
      \"lat\": \"32.1094\",
      \"lon\": \"34.8374\"
    }
  ]
}"

if command -v amqp-publish > /dev/null 2>&1; then
    amqp-publish \
    --url="amqp://guest:guest@localhost:5672" \
    --exchange="processed_events" --routing-key="#" \
    --body="$mock"
else
    echo "Please install amqp-tools (amqp-publish is missing) apt install amqp-tools"
fi