APP_NAME = oat-tracker

all: update build

.PHONY: build
build:
	go build -o ./bin/$(APP_NAME).exe -trimpath ./cmd/...

.PHONY: update
update:
	go get -u ./...
	go mod tidy

.PHONY: run
run:
	go run ./cmd/...