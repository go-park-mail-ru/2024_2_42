FROM golang:1.23-alpine3.19 AS builder

WORKDIR /usr/local/src

RUN apk --no-cache add bash make git gcc musl-dev

COPY ["go.mod", "go.sum", "./"]

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -mod=readonly -o ./bin/main ./cmd/main.go

FROM alpine:3.19 AS runner

WORKDIR /docker-42/

RUN mkdir -p /docker-42/logs

COPY --from=builder /usr/local/src/bin/main .
COPY .env .env

EXPOSE 8080

ENTRYPOINT [ "./main" ]