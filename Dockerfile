FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o gopher-watch ./cmd/gopher-watch

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/gopher-watch .
COPY --from=builder /app/configs ./configs
EXPOSE 8080
CMD ["./gopher-watch"]