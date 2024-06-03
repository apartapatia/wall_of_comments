FROM golang:1.22-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o server .

FROM alpine:3.19.1

RUN apk --no-cache add ca-certificates bash

WORKDIR /work/

COPY --from=builder /app/server .

COPY .env .

EXPOSE 8090

CMD ["./server"]
