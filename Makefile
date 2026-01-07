.PHONY: build run

APP_NAME=app
BINARY=dist/$(APP_NAME)

build:
	@mkdir -p dist
	CGO_ENABLED=0 go build -o $(BINARY) src/main.go

run: build
	./$(BINARY)
