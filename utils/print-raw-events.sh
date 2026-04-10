#!/bin/bash

if command -v amqp-consume > /dev/null 2>&1; then
    amqp-consume \
    --url="amqp://guest:guest@localhost:5672" \
    --exchange="raw_events" --routing-key="#" \
    cat
else
    echo "Please install amqp-tools"
fi