# build stage
FROM golang:1.15-alpine3.13 AS builder

RUN apk add --no-cache git make build-base
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build -tags musl -o /bin/application .

# application
FROM alpine:3.13
COPY --from=builder /bin/application /bin/application
CMD ["application"]