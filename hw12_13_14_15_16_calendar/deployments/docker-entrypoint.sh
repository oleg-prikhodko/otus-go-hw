#!/bin/bash
set -e

CONFIG_TARGET="/tmp/config.yaml"

echo "Substituting environment variables in config..."

export SCHEDULER_INTERVAL=${SCHEDULER_INTERVAL:-1m}

if [ -f "/run/config.yaml" ]; then
    echo "Using config from /run/config.yaml"
    envsubst < /run/config.yaml > "$CONFIG_TARGET"
elif [ -f "/etc/calendar/config.yaml" ]; then
    echo "Using config from /etc/calendar/config.yaml"
    envsubst < /etc/calendar/config.yaml > "$CONFIG_TARGET"
else
    echo "Warning: no config file found"
fi

echo "Starting application..."
exec "$@"