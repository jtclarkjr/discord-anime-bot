ARG GO_VERSION=1.24.1
FROM golang:${GO_VERSION}-bookworm AS builder

WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go build -v -o /run-app .

FROM debian:bookworm

RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

# Create non-root user for security
RUN groupadd -r anime && useradd -r -g anime anime

COPY --from=builder /run-app /usr/local/bin/
RUN chown anime:anime /usr/local/bin/run-app

USER anime

CMD ["run-app"]
