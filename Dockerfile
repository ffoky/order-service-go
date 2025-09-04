FROM golang:1.24.0 AS deps
WORKDIR /build

RUN apt-get update && apt-get install -y \
    librdkafka-dev \
    pkg-config \
    build-essential \
    && rm -rf /var/lib/apt/lists/*

COPY go.mod go.sum ./
RUN go mod download

FROM golang:1.24.0 AS build
WORKDIR /build

RUN apt-get update && apt-get install -y \
    librdkafka-dev \
    pkg-config \
    build-essential \
    make \
    && rm -rf /var/lib/apt/lists/*

COPY --from=deps /go/pkg/mod /go/pkg/mod
COPY go.mod go.sum ./
RUN go mod download
COPY . .

ENV CGO_ENABLED=1
RUN go build -o app cmd/app/main.go

FROM ubuntu:22.04 AS runner

WORKDIR /app

RUN apt-get update && apt-get install -y \
    curl \
    librdkafka1 \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

COPY --from=build /build/app ./app
RUN chmod +x /app/app
COPY --from=build /build/config/config.yml ./config/config.yml
COPY --from=build /build/.env .env
COPY --from=build /build/static ./static

CMD ["./app", "--config=./config/config.yml"]