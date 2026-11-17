# Multi-stage Dockerfile for kube-cost-lens
# Stage 1: Build binary statically
FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git make

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/bin/kube-cost-lens main.go

# Stage 2: Minimal runtime image
FROM alpine:3.20

RUN apk add --no-cache ca-certificates

WORKDIR /root/

COPY --from=builder /app/bin/kube-cost-lens /usr/local/bin/kube-cost-lens
COPY --from=builder /app/assets/sample-pricing.json /etc/kube-cost-lens/pricing.json

ENTRYPOINT ["kube-cost-lens"]
CMD ["--demo"]
