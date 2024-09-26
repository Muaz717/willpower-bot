.PHONY:

build:
	go build -o ./.bin/bot cmd/bot/main.go

run: build
	./.bin/bot

build-image:
	docker-compose build willpower_bot

start-container: 
	docker-compose up willpower_bot