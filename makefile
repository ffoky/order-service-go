.PHONY: up build-app rebuild-app up-without-app up-app logs-app stop down clean

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

logs-app:
	docker compose logs app

stop:
	docker compose stop

down:
	docker compose down

clean:
	docker compose down -v