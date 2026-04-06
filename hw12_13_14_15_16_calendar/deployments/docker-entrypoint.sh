#!/bin/bash
set -e

echo "Substituting environment variables in config..."
envsubst < /run/config.yaml > /tmp/config.yaml

echo "Starting application..."
exec "$@"