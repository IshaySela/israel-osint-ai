#!/bin/sh

amqp-consume \
    --url="amqp://guest:guest@localhost:5672" \
    --exchange="raw_events" --routing-key="#" \
    cat