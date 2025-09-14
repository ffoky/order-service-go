.PHONY: up build-app rebuild-app up-without-app up-app logs-app stop down clean tests mocks up-with-logs

up:
	docker compose up -d

build-app:
	docker compose build app

rebuild-app:
	docker compose up --build --force-recreate -d app

up-without-app:
	docker compose up -d db zookeeper kafka kafka-ui

up-app:
	docker compose up -d app

up-with-logs:
	docker compose up -d && docker compose logs -f app

stop:
	docker compose stop

down:
	docker compose down

clean:
	docker compose down -v

tests:
	go test -v ./internal/usecases/service/ -run TestOrderService

mocks:
	mockery --name=Order --dir=internal/infrastructure/repository --output=internal/mocks --outpkg=mocks --with-expecter --filename=order_mock.go
	mockery --name=Cache --dir=internal/usecases --output=internal/mocks --outpkg=mocks --with-expecter --filename=cache_mock.go