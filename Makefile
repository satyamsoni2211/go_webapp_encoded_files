.DEFAULT_GOAL: build

generate:
	@go run cmd/generator/main.go

build: generate
	go build -mod vendor .

.PHONY: generate build