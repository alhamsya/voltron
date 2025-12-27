# --- build stage ---
FROM golang:1.22-alpine AS builder

WORKDIR /app
RUN apk add --no-cache ca-certificates tzdata

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Sesuaikan path main kamu (contoh: ./cmd/api)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/app ./cmd

# --- runtime stage ---
FROM alpine:3.20

WORKDIR /app
RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /bin/app /bin/app

# Jika kamu pakai timezone
ENV TZ=Asia/Jakarta
EXPOSE 8080

# Healthcheck optional (kalau kamu punya endpoint /health)
# HEALTHCHECK --interval=30s --timeout=3s CMD wget -qO- http://localhost:8080/health || exit 1

ENTRYPOINT ["/bin/app"]