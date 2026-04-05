#!/bin/bash
set -e

echo "Substituting environment variables in config..."
cp /etc/calendar/config.yaml /tmp/config.yaml
envsubst < /tmp/config.yaml > /etc/calendar/config.yaml

echo "Starting application..."
exec "$@"