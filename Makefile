.PHONY: build run

APP_NAME=app
BINARY=dist/$(APP_NAME)

build:
	@mkdir -p dist
	go build -o $(BINARY) src/main.go

run: build
	./$(BINARY)
