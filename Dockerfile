FROM golang:1.22-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o server .

FROM alpine:latest

RUN apk --no-cache add ca-certificates bash

WORKDIR /work/

COPY --from=builder /app/server .

COPY .env .

CMD ["./server"]
