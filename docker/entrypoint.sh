#!/bin/sh

if [ -z "$PIPED_API" ]; then
    echo "PIPED_API not set"
    exit 1
fi

if [ -z "$HYP_API" ]; then
    echo "HYP_API not set"
    exit 1
fi

find /usr/share/nginx/html -type f -exec sed -i s/pipedapi.kavin.rocks/$PIPED_API/g {} \; -exec sed -i s/hyperpipeapi.onrender.com/$HYP_API/g {} \;

nginx -g "daemon off;"
