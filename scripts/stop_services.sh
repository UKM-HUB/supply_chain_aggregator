#!/usr/bin/env bash
# Stop all locally running services started by start_services.sh

set -euo pipefail

PID_FILE="/tmp/sca_pids"

GREEN='\033[0;32m'; YELLOW='\033[1;33m'; NC='\033[0m'

if [[ ! -f "$PID_FILE" ]]; then
  echo -e "${YELLOW}No PID file found at $PID_FILE. Services may not be running.${NC}"
  exit 0
fi

echo "Stopping services..."

while IFS= read -r pid; do
  if kill -0 "$pid" 2>/dev/null; then
    kill "$pid" 2>/dev/null && echo -e "${GREEN}  ✓ stopped PID $pid${NC}" || true
  fi
done < "$PID_FILE"

rm -f "$PID_FILE"

echo -e "${GREEN}Done.${NC}"
