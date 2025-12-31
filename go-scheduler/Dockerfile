FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod ./
RUN go mod download

COPY . .
RUN go build -o scheduler

FROM alpine:3.20

WORKDIR /app
COPY --from=builder /app/scheduler /app/scheduler

# Directory for JSON files
RUN mkdir -p /data/incoming /data/processed

ENV DATA_DIR=/data

CMD ["/app/scheduler"]
