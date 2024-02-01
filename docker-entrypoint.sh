#!/bin/sh
# entrypoint for ./Dockerfile 
set -e
SERVICE_LOG_DIR="${SERVICE_LOG_DIR:-/tmp/log}"
SERVICE_NAME="${SERVICE_NAME:-service}"
SERVICE_EXECUTABLE="${SERVICE_EXECUTABLE:-/app/server}"
# ----------------------------------------
if [ -d "$SERVICE_LOG_DIR" ]; then
  # write stdout output to a file while still outputting to stdout
  $SERVICE_EXECUTABLE 2>&1 | tee "${SERVICE_LOG_DIR}"/"${SERVICE_NAME}".json
else
  $SERVICE_EXECUTABLE
fi