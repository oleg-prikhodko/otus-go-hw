#!/bin/bash
set -e

echo "Substituting environment variables in config..."
envsubst < /etc/calendar/config.yaml > /etc/calendar/config.yaml.tmp
mv /etc/calendar/config.yaml.tmp /etc/calendar/config.yaml

echo "Starting application..."
exec "$@"