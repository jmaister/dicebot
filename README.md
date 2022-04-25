
# DiceBot

Telegram bot to throw dices.

[t.me/some_dice_bot](http://t.me/some_dice_bot)


Send a message with the way Dungeons & Dragons describes the dice throws and you get the result.

    1d20
    2d20
    2d10 1d8


# Build docker

docker build -t dicebot .

# Run docker

make docker-run DICEBOT_TOKEN=<telegram_token>