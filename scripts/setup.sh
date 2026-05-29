#!/usr/bin/env bash
# setup.sh — Download semua dependency dan verifikasi setiap service bisa di-build.
#
# Jalankan SEKALI setelah clone / extract zip, sebelum `make run` atau e2e_flow.sh.
#
# Usage:
#   cd supply_chain_aggregator
#   bash scripts/setup.sh

set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

GREEN='\033[0;32m'; YELLOW='\033[1;33m'; RED='\033[0;31m'; BOLD='\033[1m'; NC='\033[0m'
ok()   { echo -e "${GREEN}  ✓ $1${NC}"; }
info() { echo -e "${YELLOW}  → $1${NC}"; }
fail() { echo -e "${RED}  ✗ $1${NC}"; exit 1; }

require() {
  command -v "$1" &>/dev/null || fail "Tool tidak ditemukan: $1. Silakan install terlebih dahulu."
}

echo -e "${BOLD}Supply Chain Aggregator — Setup${NC}"
echo -e "Root: $ROOT\n"

require go
require git

GO_VER=$(go version | awk '{print $3}' | tr -d 'go')
info "Go version: $GO_VER"

# ── 1. Tidy shared pkg dulu (services depend on pkg) ─────────────────────────
info "Step 1/2 — go mod tidy: pkg"
(cd pkg && go mod tidy) && ok "pkg selesai"

# ── 2. Tidy setiap service ────────────────────────────────────────────────────
SERVICES=(
  api-gateway
  auth-service
  sme-service
  nearby-service
  transaction-service
  payment-service
  communication-service
  report-service
)

info "Step 2/2 — go mod tidy: semua service"
for svc in "${SERVICES[@]}"; do
  dir="services/$svc"
  if [[ -f "$dir/go.mod" ]]; then
    (cd "$dir" && go mod tidy) && ok "$svc"
  else
    echo -e "${YELLOW}  skip $svc (go.mod tidak ditemukan)${NC}"
  fi
done

echo ""
echo -e "${GREEN}${BOLD}✓ Setup selesai!${NC}"
echo ""
echo "Langkah selanjutnya:"
echo "  1. make docker-infra-up          — start Redis, Elasticsearch, Postgres, RabbitMQ"
echo "  2. make run SVC=auth-service     — jalankan service satu per satu, atau"
echo "  3. bash scripts/e2e_flow.sh      — jalankan E2E flow otomatis"
