# Supply Chain Aggregator — root Makefile
#
# Suggested build order (mirrors Step 15):
#   1. make tidy-all          — download all module dependencies
#   2. make build-all         — compile every service
#   3. make docker-infra-up   — start PostgreSQL + RabbitMQ
#   4. make migrate-up        — run all SQL migrations
#   5. make run SVC=auth-service      (then sme, nearby, transaction, payment, communication, report)
#   6. make e2e               — run end-to-end flow
#   7. make docker-up         — build + start full stack in containers

.DEFAULT_GOAL := help

SERVICES := \
	auth-service \
	sme-service \
	nearby-service \
	transaction-service \
	payment-service \
	communication-service \
	report-service

SERVICES_DIR  := services
PROTO_DIR     := proto
MIGRATIONS_DIR := migrations

# Database defaults (override via env or .env)
POSTGRES_DSN  ?= postgres://sca:sca_pass@localhost:5432/supply_chain?sslmode=disable
RABBITMQ_URL  ?= amqp://guest:guest@localhost:5672/

# ── Help ──────────────────────────────────────────────────────────────────────
.PHONY: help
help:
	@echo ""
	@echo "Supply Chain Aggregator — available targets"
	@echo ""
	@echo "  Service targets (set SVC=<service-name>):"
	@echo "    build          Build a single service binary"
	@echo "    run            Run a single service locally"
	@echo "    tidy           Run go mod tidy for a single service"
	@echo ""
	@echo "  Bulk targets:"
	@echo "    build-all      Build all services"
	@echo "    tidy-all       Run go mod tidy for all services"
	@echo ""
	@echo "  Proto:"
	@echo "    proto          Generate Go code from all .proto files"
	@echo ""
	@echo "  Database:"
	@echo "    migrate-up     Apply all migrations"
	@echo "    migrate-down   Roll back the last migration"
	@echo "    migrate-reset  Roll back all migrations"
	@echo ""
	@echo "  Docker:"
	@echo "    docker-infra-up    Start PostgreSQL + RabbitMQ"
	@echo "    docker-infra-down  Stop infrastructure containers"
	@echo "    docker-up          Build + start full stack"
	@echo "    docker-down        Stop full stack"
	@echo "    docker-logs        Tail logs from all containers"
	@echo ""
	@echo "  Testing:"
	@echo "    e2e            Run the end-to-end MVP flow script"
	@echo ""
	@echo "  Example usage:"
	@echo "    make build SVC=auth-service"
	@echo "    make run   SVC=transaction-service"
	@echo "    make tidy  SVC=payment-service"
	@echo ""

# ── Single-service targets ────────────────────────────────────────────────────
.PHONY: build
build:
	@test -n "$(SVC)" || (echo "Usage: make build SVC=<service-name>"; exit 1)
	@echo "Building $(SVC)..."
	cd $(SERVICES_DIR)/$(SVC) && go build ./cmd/$(SVC)/...
	@echo "Done: $(SVC)"

.PHONY: run
run:
	@test -n "$(SVC)" || (echo "Usage: make run SVC=<service-name>"; exit 1)
	@echo "Running $(SVC)..."
	cd $(SERVICES_DIR)/$(SVC) && go run ./cmd/$(SVC)/...

.PHONY: tidy
tidy:
	@test -n "$(SVC)" || (echo "Usage: make tidy SVC=<service-name>"; exit 1)
	@echo "Tidying $(SVC)..."
	cd $(SERVICES_DIR)/$(SVC) && go mod tidy

# ── Bulk targets ──────────────────────────────────────────────────────────────
.PHONY: build-all
build-all:
	@echo "Building all services..."
	@for svc in $(SERVICES); do \
		echo "  → $$svc"; \
		(cd $(SERVICES_DIR)/$$svc && go build ./cmd/$$svc/...) || exit 1; \
	done
	@echo "All services built."

.PHONY: tidy-all
tidy-all:
	@echo "Running go mod tidy for all services..."
	@for svc in $(SERVICES); do \
		echo "  → $$svc"; \
		(cd $(SERVICES_DIR)/$$svc && go mod tidy) || exit 1; \
	done
	@echo "Done."

# ── Proto generation ──────────────────────────────────────────────────────────
.PHONY: proto
proto:
	@echo "Generating protobuf code..."
	@command -v protoc >/dev/null 2>&1 || \
		(echo "protoc not found. Install: https://grpc.io/docs/protoc-installation/"; exit 1)
	@command -v protoc-gen-go >/dev/null 2>&1 || \
		(echo "protoc-gen-go not found. Run: go install google.golang.org/protobuf/cmd/protoc-gen-go@latest"; exit 1)
	@command -v protoc-gen-go-grpc >/dev/null 2>&1 || \
		(echo "protoc-gen-go-grpc not found. Run: go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest"; exit 1)
	protoc \
		--go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		$(PROTO_DIR)/auth/auth.proto \
		$(PROTO_DIR)/user/user.proto \
		$(PROTO_DIR)/sme/sme.proto \
		$(PROTO_DIR)/nearby/nearby.proto \
		$(PROTO_DIR)/transaction/transaction.proto \
		$(PROTO_DIR)/payment/payment.proto
	@echo "Protobuf code generated."

# ── Database migrations ───────────────────────────────────────────────────────
.PHONY: migrate-up
migrate-up:
	@command -v migrate >/dev/null 2>&1 || \
		(echo "golang-migrate not found. Run: go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest"; exit 1)
	migrate -path $(MIGRATIONS_DIR) -database "$(POSTGRES_DSN)" up
	@echo "Migrations applied."

.PHONY: migrate-down
migrate-down:
	@command -v migrate >/dev/null 2>&1 || \
		(echo "golang-migrate not found."; exit 1)
	migrate -path $(MIGRATIONS_DIR) -database "$(POSTGRES_DSN)" down 1
	@echo "Last migration rolled back."

.PHONY: migrate-reset
migrate-reset:
	@command -v migrate >/dev/null 2>&1 || \
		(echo "golang-migrate not found."; exit 1)
	migrate -path $(MIGRATIONS_DIR) -database "$(POSTGRES_DSN)" down
	@echo "All migrations rolled back."

# ── Docker ────────────────────────────────────────────────────────────────────
.PHONY: docker-infra-up
docker-infra-up:
	docker compose -f deployments/docker-compose.infra.yml up -d
	@echo "Infrastructure started. PostgreSQL: localhost:5432  RabbitMQ: localhost:5672 (UI: localhost:15672)"

.PHONY: docker-infra-down
docker-infra-down:
	docker compose -f deployments/docker-compose.infra.yml down

.PHONY: docker-up
docker-up:
	docker compose -f deployments/docker-compose.yml up -d --build
	@echo "Full stack started."

.PHONY: docker-down
docker-down:
	docker compose -f deployments/docker-compose.yml down

.PHONY: docker-logs
docker-logs:
	docker compose -f deployments/docker-compose.yml logs -f

# ── E2E ───────────────────────────────────────────────────────────────────────
.PHONY: e2e
e2e:
	bash scripts/e2e_flow.sh
