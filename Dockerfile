FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY shared/ ./shared/
COPY server/ ./server/

RUN ls -la
RUN ls -la shared/
RUN ls -la server/

RUN go build -v -o hr-server ./server/
RUN ls -la

FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/hr-server ./hr-server

EXPOSE 8888

CMD ["./hr-server"]
