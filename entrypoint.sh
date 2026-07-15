#!/bin/sh
set -e

PORT=8081 /app/bin/go-library &
PORT=8082 /app/bin/go-bookshop &
caddy run --config /app/Caddyfile --adapter caddyfile &

wait -n
exit $?
