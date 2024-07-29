deploy:
	go build cmd/goat-cg/main.go
	nohup ./main &

up:
	docker compose up -d

down:
	docker compose down

start:
	docker compose start

stop:
	docker compose stop

in:
	docker exec -i -t goat-cg_app bash

db:
	docker exec -i -t goat-cg_db bash

build:
	docker compose build --no-cache

run:
	go run cmd/goat-cg/main.go

test:
	go test



