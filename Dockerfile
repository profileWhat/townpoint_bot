# Builder
FROM golang:1.23-alpine AS builder

RUN apk add --no-cache ca-certificates

WORKDIR /app

# GitHub private connections
ARG GITHUB_TOKEN

RUN apk update && apk add git

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# Build
COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -tags timetzdata -o /bin/app ./cmd/app

# Runner
FROM scratch AS runner

COPY --from=builder /app/config.toml /config.toml
#COPY --from=builder /app/migrations /migrations
COPY --from=builder /bin/app /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 8080
CMD ["/app"]