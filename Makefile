.SILENT:

CONFIG_PATH = C:\Users\sasha\Documents\go\UrlShortener\config\config.yaml
APP_SECRET = test_secret

DB_URL_DOCKER = postgres://postgres:postgres@shortener-db:5432/shortener?sslmode=disable
DB_URL_LOCAL = postgres://postgres:1234@localhost:5432/shortener?sslmode=disable


export CONFIG_PATH
export APP_SECRET

set-env:
	@echo "CONFIG_PATH is $$CONFIG_PATH"
	@echo "APP_SECRET is $$APP_SECRET"

run-local: set-env
	go run ./cmd/ .

run-docker: 
	docker-compose up --build

create-migration:
	migrate create -ext sql -dir ./migrations -seq $(NAME)

migrate-up-local:
	migrate -path ./migrations -database "$(DB_URL_LOCAL)" up

migrate-down-local:
	migrate -path ./migrations -database "$(DB_URL_LOCAL)" down 1

migrate-up-docker:
	migrate -path ./migrations -database "$(DB_URL_DOCKER)" up

migrate-down-docker:
	migrate -path ./migrations -database "$(DB_URL_DOCKER)" down 1

