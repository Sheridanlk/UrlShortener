.SILENT:

CONFIG_PATH = C:\Users\sasha\Documents\go\UrlShortener\config\cfg.yaml
export CONFIG_PATH

set-env:
	@echo "CONFIG_PATH is $$CONFIG_PATH"

run: set-env
	go run ./cmd/ .