.SILENT:

CONFIG_PATH = C:\Users\sasha\Documents\go\UrlShortener\config\cfg.yaml
APP_SECRET = test_secret
export CONFIG_PATH
export APP_SECRET

set-env:
	@echo "CONFIG_PATH is $$CONFIG_PATH"
	@echo "APP_SECRET is $$APP_SECRET"

run: set-env
	go run ./cmd/ .

create-migration:
	migrate create -ext sql -dir ./migrations -seq $(NAME)

migrate:
	migrate -path ./migrations -database 'postgres://postgres:1234@localhost:5432/urlShortener?sslmode=disable' up

migrate-down:
	migrate -path ./migrations -database 'postgres://postgres:postgres@localhost:5432/urlShortener?sslmode=disable' down

