.PHONY: run build clean dev

GO ?= go
APP_NAME = sheet-server

run:
	$(GO) run ./cmd/server --config config.yaml

build:
	$(GO) build -o bin/$(APP_NAME) ./cmd/server

clean:
	rm -rf bin/

dev:
	$(GO) run ./cmd/server --config config.yaml

tidy:
	$(GO) mod tidy

air:
	air -- --config config.yaml
