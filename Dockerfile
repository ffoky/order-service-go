FROM golang:1.24.0-alpine AS deps
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download

FROM golang:1.24.0-alpine AS build
WORKDIR /build
COPY --from=deps /go/pkg/mod /go/pkg/mod
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN apk add --no-cache make
RUN go build -o app cmd/app/main.go

FROM alpine AS runner

WORKDIR /app

RUN apk add --no-cache curl

COPY --from=build /build/app ./app
RUN chmod +x /app/app
COPY --from=build /build/cmd/app/config/config.yml ./config.yml
COPY --from=build /build/.env .env

CMD ["./app", "--config=/app/config.yml" ]
