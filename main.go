package main

import (
	"log"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	// Telegram token
	token := os.Getenv("DICEBOT_TOKEN")

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			response := ""
			dices := parseDiceTrows(update.Message.Text)
			for _, d := range dices {
				line := ""
				if d.Ok {
					line = strconv.Itoa(d.Times) + "d" + strconv.Itoa(d.Max) + " 🎲 "
					for i := 0; i < d.Times; i++ {
						line = line + strconv.Itoa(getRandom(d.Max)) + " "
					}
				} else {
					line = d.Msg
				}
				response = response + line + "\n"
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
			msg.ReplyToMessageID = update.Message.MessageID

			bot.Send(msg)

		}
	}
}

type DiceThrow struct {
	Times int
	Max   int
	Msg   string
	Ok    bool
}

func parseDiceTrows(message string) []DiceThrow {
	diceThrows := []DiceThrow{}

	r, err := regexp.Compile("(\\d+)d(\\d+)")
	if err != nil {
		panic("Error on main regex")
	}
	message = strings.ToLower(message)
	matches := r.FindAllStringSubmatch(message, -1)

	if len(matches) == 0 {
		diceThrows = append(diceThrows, DiceThrow{
			Msg: "'" + message + "' is not valid.",
			Ok:  false,
		})
		return diceThrows
	}

	for _, m := range matches {
		tms, err := strconv.Atoi(m[1])
		if err != nil {
			log.Println("err 1", err)
			diceThrows = append(diceThrows, DiceThrow{
				Msg: strings.Join(m, "") + " is not valid.",
			})
			continue
		}
		mx, err := strconv.Atoi(m[2])
		if err != nil {
			log.Println("err 2", err)
			diceThrows = append(diceThrows, DiceThrow{
				Msg: strings.Join(m, "") + " is not valid.",
				Ok:  false,
			})
			continue
		}

		diceThrows = append(diceThrows, DiceThrow{
			Times: tms,
			Max:   mx,
			Ok:    true,
		})
	}

	return diceThrows
}

func getRandom(max int) int {
	return rand.Intn(max) + 1
}
