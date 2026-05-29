#!/usr/bin/env bash
# Start all services locally for development.
# Each service runs in the background; a PID file is written to /tmp/sca_pids.
# Stop everything with: ./scripts/stop_services.sh

set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
PID_FILE="/tmp/sca_pids"

GREEN='\033[0;32m'; YELLOW='\033[1;33m'; RED='\033[0;31m'; NC='\033[0m'

ok()   { echo -e "${GREEN}  ✓ $1${NC}"; }
info() { echo -e "${YELLOW}  → $1${NC}"; }
fail() { echo -e "${RED}  ✗ $1${NC}"; exit 1; }

if [[ -f "$PID_FILE" ]]; then
  echo -e "${YELLOW}Services may already be running (found $PID_FILE).${NC}"
  echo -e "${YELLOW}Run ./scripts/stop_services.sh first, then retry.${NC}"
  exit 1
fi

> "$PID_FILE"

start() {
  local name=$1 port=$2
  info "Starting $name on :$port ..."
  ( cd "$REPO_ROOT/services/$name" && \
    HTTP_PORT="$port" go run "./cmd/$name/..." \
    > "/tmp/sca_${name}.log" 2>&1 ) &
  echo "$!" >> "$PID_FILE"
}

wait_healthy() {
  local name=$1 url=$2
  for i in $(seq 1 30); do
    if curl -sf "$url" &>/dev/null; then
      ok "$name  :$(echo "$url" | grep -o '[0-9]*/'| tr -d /)"
      return
    fi
    sleep 1
  done
  fail "$name did not start (check /tmp/sca_${name}.log)"
}

echo ""
echo "Supply Chain Aggregator — starting all services"
echo "Logs: /tmp/sca_<service>.log"
echo ""

start "auth-service"        8081
start "sme-service"         8082
start "nearby-service"      8083
start "transaction-service" 8084
start "payment-service"     8085
start "communication-service" 8086
start "report-service"      8087

sleep 3

echo ""
echo "Waiting for health checks..."
wait_healthy "auth-service"          "http://localhost:8081/health"
wait_healthy "sme-service"           "http://localhost:8082/health"
wait_healthy "nearby-service"        "http://localhost:8083/health"
wait_healthy "transaction-service"   "http://localhost:8084/health"
wait_healthy "payment-service"       "http://localhost:8085/health"
wait_healthy "communication-service" "http://localhost:8086/health"
wait_healthy "report-service"        "http://localhost:8087/health"

echo ""
echo -e "${GREEN}All services running.${NC}"
echo ""
echo "  Service               Port   Log"
echo "  ──────────────────────────────────────────────────────"
echo "  auth-service          8081   /tmp/sca_auth-service.log"
echo "  sme-service           8082   /tmp/sca_sme-service.log"
echo "  nearby-service        8083   /tmp/sca_nearby-service.log"
echo "  transaction-service   8084   /tmp/sca_transaction-service.log"
echo "  payment-service       8085   /tmp/sca_payment-service.log"
echo "  communication-service 8086   /tmp/sca_communication-service.log"
echo "  report-service        8087   /tmp/sca_report-service.log"
echo ""
echo "  Stop with: ./scripts/stop_services.sh"
echo "  E2E test:  ./scripts/e2e_flow.sh --no-start"
echo ""
