
FROM golang:1.18-alpine

RUN mkdir /app
ADD . /app
WORKDIR /app

RUN go mod download
RUN go build -o main .

# Telegram token
ENV DICEBOT_TOKEN=$DICEBOT_TOKEN

CMD ["/app/main"]
