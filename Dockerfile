FROM golang:alpine as builder

ENV GO111MODULE=on

WORKDIR /app

RUN apk add build-base

COPY go.mod .

RUN go mod download

COPY . .

RUN go build -o urls-api ./server/http

FROM alpine

RUN apk update && apk add sqlite

COPY --from=builder /app/urls-api .
COPY --from=builder /app/database/sqlite ./database/sqlite
COPY --from=builder /app/swagger.yaml .

CMD ["/urls-api"]

