.PHONY: build run

build:
	@go build -o bin/mailer

run: build
	@bin/mailer
