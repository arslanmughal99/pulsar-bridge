FROM golang:alpine AS builder

WORKDIR /build

COPY . .

RUN go mod download

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

RUN apk add git

RUN go build -ldflags="-s -w" -o app .

RUN chmod +x -R /build

FROM alpine

COPY --from=builder /build/app /app

CMD ["sh", "-c", "/app"]