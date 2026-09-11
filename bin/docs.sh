#!/bin/bash

set -e

cd "$(dirname "$0")/.."

PORT=6060
URL="http://localhost:${PORT}/pkg/monkey/"

if ! command -v godoc >/dev/null 2>&1; then
    echo "godoc not found, installing..."
    go install golang.org/x/tools/cmd/godoc@latest
fi

# Start the doc server in the background if it isn't already running.
if ! (exec 3<>"/dev/tcp/localhost/${PORT}") 2>/dev/null; then
    godoc -http=":${PORT}" >/dev/null 2>&1 &
    sleep 1
fi

if command -v xdg-open >/dev/null 2>&1; then
    xdg-open "$URL" >/dev/null 2>&1 &
elif command -v open >/dev/null 2>&1; then
    open "$URL"
else
    echo "Open this URL in your browser: $URL"
fi
