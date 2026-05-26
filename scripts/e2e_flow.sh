#!/usr/bin/env bash
# End-to-end MVP flow demonstration for supply_chain_aggregator.
#
# Walks through the complete business flow:
#   SME registers → SME creates profile → Corporation logs in →
#   Corporation searches nearby SMEs → Corporation creates transaction →
#   Payment virtual account created → Xendit webhook (PAID) →
#   RabbitMQ publishes payment.paid → WhatsApp notification sent →
#   Report data available
#
# Usage:
#   ./scripts/e2e_flow.sh           # starts services, runs flow, cleans up
#   ./scripts/e2e_flow.sh --no-start  # skip service startup (already running)

set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
NO_START="${1:-}"

# ── colours ──────────────────────────────────────────────────────────────────
BLUE='\033[0;34m'; GREEN='\033[0;32m'; YELLOW='\033[1;33m'
RED='\033[0;31m';  BOLD='\033[1m';     NC='\033[0m'

step()  { echo -e "\n${BLUE}${BOLD}━━━ $1 ━━━${NC}"; }
ok()    { echo -e "${GREEN}  ✓ $1${NC}"; }
info()  { echo -e "${YELLOW}  → $1${NC}"; }
fail()  { echo -e "${RED}  ✗ $1${NC}"; exit 1; }
json()  { echo "$1" | python3 -m json.tool 2>/dev/null || echo "$1"; }

require() {
  command -v "$1" &>/dev/null || fail "required tool not found: $1"
}

# ── helpers ───────────────────────────────────────────────────────────────────
wait_for() {
  local name=$1 url=$2
  info "Waiting for $name at $url ..."
  for i in $(seq 1 20); do
    if curl -sf "$url" &>/dev/null; then
      ok "$name is up"
      return
    fi
    sleep 1
  done
  fail "$name did not become healthy at $url"
}

extract() {
  # extract a JSON field using python3 (no jq dependency)
  python3 -c "import sys,json; d=json.load(sys.stdin); print(d$1)" 2>/dev/null
}

PIDS=()
start_service() {
  local name=$1 dir=$2 port=$3
  info "Starting $name on :$port ..."
  (cd "$REPO_ROOT/services/$dir" && go run "./cmd/$name/..." ) &
  PIDS+=($!)
}

cleanup() {
  echo -e "\n${YELLOW}Stopping services...${NC}"
  for pid in "${PIDS[@]:-}"; do
    kill "$pid" 2>/dev/null || true
  done
  # also kill any stray processes on the ports
  for port in 8081 8082 8083 8084 8085 8087; do
    lsof -ti ":$port" 2>/dev/null | xargs kill -9 2>/dev/null || true
  done
  echo -e "${GREEN}Done.${NC}"
}

# ── preflight ─────────────────────────────────────────────────────────────────
require curl
require python3
require go

echo -e "${BOLD}Supply Chain Aggregator — E2E MVP Flow${NC}"
echo -e "Repo: $REPO_ROOT\n"

# ── start services ────────────────────────────────────────────────────────────
if [[ "$NO_START" != "--no-start" ]]; then
  trap cleanup EXIT

  start_service "auth-service"        "auth-service"        8081
  start_service "sme-service"         "sme-service"         8082
  start_service "nearby-service"      "nearby-service"      8083
  start_service "transaction-service" "transaction-service" 8084
  start_service "payment-service"     "payment-service"     8085
  start_service "report-service"      "report-service"      8087

  sleep 3

  wait_for "auth-service"        "http://localhost:8081/health"
  wait_for "sme-service"         "http://localhost:8082/health"
  wait_for "nearby-service"      "http://localhost:8083/health"
  wait_for "transaction-service" "http://localhost:8084/health"
  wait_for "payment-service"     "http://localhost:8085/health"
  wait_for "report-service"      "http://localhost:8087/health"
fi

# ─────────────────────────────────────────────────────────────────────────────
# STEP 1: SME registers account
# ─────────────────────────────────────────────────────────────────────────────
step "Step 1 — SME registers account"

SME_REGISTER=$(curl -sf -X POST http://localhost:8081/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name":      "UMKM Maju Food",
    "email":     "maju@umkm.test",
    "phone":     "628111111111",
    "password":  "secret123",
    "role":      "SME",
    "latitude":  -6.2245,
    "longitude": 106.8099
  }')

SME_USER_ID=$(echo "$SME_REGISTER" | extract "['user']['id']")
ok "SME registered — id: $SME_USER_ID"
json "$SME_REGISTER"

# ─────────────────────────────────────────────────────────────────────────────
# STEP 2: SME creates business profile
# ─────────────────────────────────────────────────────────────────────────────
step "Step 2 — SME creates business profile with category and location"

SME_PROFILE=$(curl -sf -X POST http://localhost:8082/api/v1/umkm \
  -H "Content-Type: application/json" \
  -d "{
    \"owner_id\":    \"$SME_USER_ID\",
    \"name\":        \"UMKM Maju Food\",
    \"phone\":       \"628111111111\",
    \"address\":     \"Jakarta Selatan\",
    \"description\": \"Food and snack supplier\",
    \"category_ids\":[\"food\"],
    \"products\":    [\"frozen food\",\"snack\"],
    \"capacity\":    \"500 pcs per week\",
    \"latitude\":    -6.2245,
    \"longitude\":   106.8099,
    \"status\":      \"active\"
  }")

SME_PROFILE_ID=$(echo "$SME_PROFILE" | extract "['data']['id']")
ok "SME profile created — id: $SME_PROFILE_ID"
json "$SME_PROFILE"

# ─────────────────────────────────────────────────────────────────────────────
# STEP 3: Corporation logs in
# ─────────────────────────────────────────────────────────────────────────────
step "Step 3 — Corporation registers and logs in"

curl -sf -X POST http://localhost:8081/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name":     "PT Maju Bersama",
    "email":    "corp@factory.test",
    "phone":    "628222222222",
    "password": "corp123",
    "role":     "CORPORATION"
  }' > /dev/null

CORP_LOGIN=$(curl -sf -X POST http://localhost:8081/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"corp@factory.test","password":"corp123"}')

CORP_TOKEN=$(echo "$CORP_LOGIN"  | extract "['token']")
CORP_USER_ID=$(echo "$CORP_LOGIN" | extract "['user']['id']")
ok "Corporation logged in — user_id: $CORP_USER_ID"
info "JWT token: ${CORP_TOKEN:0:40}..."

# ─────────────────────────────────────────────────────────────────────────────
# STEP 4: Corporation searches nearby SMEs
# ─────────────────────────────────────────────────────────────────────────────
step "Step 4 — Corporation searches nearby SMEs (lat=-6.2, lng=106.8, category=food)"

NEARBY=$(curl -sf \
  "http://localhost:8083/api/v1/nearby/umkm?lat=-6.2&lng=106.8&category_id=food&radius_km=10&limit=5")

NEARBY_COUNT=$(echo "$NEARBY" | extract "['total']")
ok "Found $NEARBY_COUNT nearby SME(s) matching category=food within 10 km"
json "$NEARBY"

# ─────────────────────────────────────────────────────────────────────────────
# STEP 5: Corporation creates a transaction
# ─────────────────────────────────────────────────────────────────────────────
step "Step 5 — Corporation creates a transaction"

TXN=$(curl -sf -X POST http://localhost:8084/api/v1/transactions \
  -H "Content-Type: application/json" \
  -d "{
    \"user_id\":        \"$CORP_USER_ID\",
    \"amount\":         3000000,
    \"payment_method\": \"virtual_account\"
  }")

TXN_ID=$(echo "$TXN"      | extract "['data']['id']")
INVOICE=$(echo "$TXN"     | extract "['data']['invoice_number']")
TXN_STATUS=$(echo "$TXN"  | extract "['data']['status']")

ok "Transaction created — id: $TXN_ID  invoice: $INVOICE  status: $TXN_STATUS"
json "$TXN"

# ─────────────────────────────────────────────────────────────────────────────
# STEP 6: Payment virtual account is created
# ─────────────────────────────────────────────────────────────────────────────
step "Step 6 — Create Xendit virtual account for the transaction"

VA=$(curl -sf -X POST http://localhost:8085/api/v1/gateway/create-va \
  -H "Content-Type: application/json" \
  -d "{
    \"invoice_number\": \"$INVOICE\",
    \"amount\":         3000000,
    \"user_phone\":     \"628222222222\"
  }")

PAYMENT_URL=$(echo "$VA" | extract "['data']['payment_url']")
ok "Virtual account created — payment URL: $PAYMENT_URL"
json "$VA"

# ─────────────────────────────────────────────────────────────────────────────
# STEP 7: Xendit webhook — payment is PAID
# ─────────────────────────────────────────────────────────────────────────────
step "Step 7 — Simulate Xendit webhook (PAID)"

WEBHOOK=$(curl -sf -X POST http://localhost:8085/api/v1/webhooks/xendit \
  -H "Content-Type: application/json" \
  -d "{
    \"id\":          \"xendit-mock-001\",
    \"external_id\": \"$INVOICE\",
    \"status\":      \"PAID\",
    \"amount\":      3000000,
    \"paid_at\":     \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"
  }")

ok "Webhook processed"
json "$WEBHOOK"
info "→ payment.paid event published to RabbitMQ (fire-and-forget)"
info "→ communication-service consumes and sends WhatsApp to 628222222222"

# ─────────────────────────────────────────────────────────────────────────────
# STEP 8: Verify transaction status
# ─────────────────────────────────────────────────────────────────────────────
step "Step 8 — Verify transaction detail"

TXN_DETAIL=$(curl -sf "http://localhost:8084/api/v1/transactions/$TXN_ID")
FINAL_STATUS=$(echo "$TXN_DETAIL" | extract "['data']['status']")
ok "Transaction $TXN_ID — status: $FINAL_STATUS"
json "$TXN_DETAIL"

# ─────────────────────────────────────────────────────────────────────────────
# STEP 9: Report data is available
# ─────────────────────────────────────────────────────────────────────────────
step "Step 9 — Report data is available"

TODAY=$(date +%Y-%m-%d)
DAILY=$(curl -sf "http://localhost:8087/api/v1/reports/daily?date=$TODAY")
ok "Daily report for $TODAY"
json "$DAILY"

MONTHLY=$(curl -sf "http://localhost:8087/api/v1/reports/monthly?year=$(date +%Y)&month=$(date +%-m)")
ok "Monthly report for $(date +%Y-%m)"
json "$MONTHLY"

# ─────────────────────────────────────────────────────────────────────────────
echo -e "\n${GREEN}${BOLD}✓ E2E MVP flow completed successfully${NC}"
echo -e "${YELLOW}Note: services use in-memory storage — data is not shared between${NC}"
echo -e "${YELLOW}      services or persisted across restarts. Run migrations and${NC}"
echo -e "${YELLOW}      connect POSTGRES_URL / RABBITMQ_URL for the full stack.${NC}\n"
