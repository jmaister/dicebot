

TOKEN=$(DICEBOT_TOKEN)

build:
	go build -o main main.go

run:
	go run main.go

test:
	go test

docker-build:
	docker build -t dicebot .

docker-run:
	docker run -it -d --restart=always --name dicebot01 --env DICEBOT_TOKEN=$(TOKEN) dicebot

docker-rm:
	docker rm --force dicebot01

rebuild: test docker-build docker-rm docker-run