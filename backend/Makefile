.PHONY: build run stop clean dev logs migrate migrate-up migrate-down migrate-status migrate-create

# Docker commands
build:
	docker-compose build

dev:
	docker-compose up -d

down:
	docker-compose down

clean:
	docker-compose down -v
	docker system prune -f

reset:
	make down
	make dev

ps:
	docker-compose ps

logs:
	docker-compose logs -f

gen:
	gqlgen generate

# Database commands
db:
	docker-compose exec mysql mysql -u root -p graphql_db

# Migration commands
migrate-up:
	go run cmd/migrate/main.go -action=up

migrate-down:
	go run cmd/migrate/main.go -action=down -steps=1

migrate-down-all:
	go run cmd/migrate/main.go -action=down -steps=100

migrate-status:
	go run cmd/migrate/main.go -action=status

migrate-create:
	@if [ -z "$(name)" ]; then \
		echo "Usage: make migrate-create name=migration_name"; \
		exit 1; \
	fi
	go run cmd/migrate/main.go -action=create -name="$(name)"
