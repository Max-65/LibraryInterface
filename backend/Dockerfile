FROM golang:1.26.3-alpine AS builder

WORKDIR /src

RUN apk add --no-cache file ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /out/server ./cmd/server
RUN file /out/server

FROM alpine:3.22

WORKDIR /app

RUN apk add --no-cache file ca-certificates

COPY --from=builder /out/server /app/server
RUN chmod +x /app/server
RUN file /app/server

EXPOSE 8080
ENTRYPOINT ["/app/server"]