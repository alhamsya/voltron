# --- build stage ---
FROM golang:1.22-alpine AS builder

WORKDIR /app
RUN apk add --no-cache ca-certificates tzdata

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# change ./cmd/server to your main package path if different
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

# --- runtime stage ---
FROM alpine:3.20
WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /app/server /app/server

ENV TZ=Asia/Jakarta

ENTRYPOINT ["/app/server", "run", "rest"]