# Production
deploy:
	go build cmd/goat-cg/main.go
	nohup ./main &

# Local
lrun:
	ENV=local go run cmd/goat-cg/main.go

ltest:
	ENV=local go test

# Local (docker compose)
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

# Docker Container
run:
	go run cmd/goat-cg/main.go



