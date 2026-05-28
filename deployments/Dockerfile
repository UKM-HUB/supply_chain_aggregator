# Generic multi-stage Dockerfile for any service in this monorepo.
# Build with: docker build --build-arg SERVICE_NAME=auth-service -t auth-service .
# Must be built from the repo root as the build context.

ARG SERVICE_NAME

# ── Stage 1: Build ──────────────────────────────────────────────────────────
FROM golang:1.23-alpine AS builder

ARG SERVICE_NAME
ENV SERVICE_NAME=${SERVICE_NAME}

WORKDIR /app

COPY services/${SERVICE_NAME}/go.mod services/${SERVICE_NAME}/go.sum ./
RUN go mod download

COPY services/${SERVICE_NAME} .

RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -o /bin/service ./cmd/${SERVICE_NAME}/...

# ── Stage 2: Runtime ─────────────────────────────────────────────────────────
FROM alpine:3.20

RUN apk --no-cache add ca-certificates tzdata

COPY --from=builder /bin/service /bin/service

EXPOSE 8080

CMD ["/bin/service"]
