package main

import (
	"crypto/rand"
	"log"
	"math/big"
	"os"
	"regexp"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const HelpStr = "Welcome to Dice Bot!!\n\n" +
	"I can throw messages if you write <number>d<max> i.e. 1d20 like in Dungeons And Dragons.\n\n" +
	"/show Show common dices\n" +
	"/close Close buttons panel"

var diceKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("1d20"),
		tgbotapi.NewKeyboardButton("1d12"),
		tgbotapi.NewKeyboardButton("1d10"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("1d8"),
		tgbotapi.NewKeyboardButton("1d6"),
		tgbotapi.NewKeyboardButton("1d4"),
	),
)

func main() {
	// Telegram token
	token := os.Getenv("DICEBOT_TOKEN")
	if token == "" {
		panic("DICEBOT_TOKEN must be set.")
	}

	log.Println("Starting with token " + token[0:3] + "..." + token[len(token)-3:])

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}
	bot.Self.SupportsInlineQueries = true

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil && update.Message.Command() != "" {
			command := update.Message.Command()
			if command == "show" {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Showing buttons.")
				msg.ReplyMarkup = diceKeyboard
				bot.Send(msg)
			} else if command == "close" {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Closing buttons.")
				msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
				bot.Send(msg)
			} else {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, HelpStr)
				msg.ReplyMarkup = diceKeyboard
				bot.Send(msg)
			}
		} else if update.MyChatMember != nil {
			msg := tgbotapi.NewMessage(update.MyChatMember.Chat.ID, HelpStr)
			msg.ReplyMarkup = diceKeyboard
			bot.Send(msg)

		} else if update.Message.Text != "" {
			go processMessage(bot, update)
		}

	}
}

type DiceThrow struct {
	Times int
	Max   int
	Msg   string
	Ok    bool
}

func processMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
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

func parseDiceTrows(message string) []DiceThrow {
	diceThrows := []DiceThrow{}

	r, err := regexp.Compile("(\\d+)d(\\d+)")
	if err != nil {
		panic("Error on parseDiceTrows regex")
	}
	message = strings.ToLower(message)
	matches := r.FindAllStringSubmatch(message, -1)

	if len(matches) == 0 {
		diceThrows = append(diceThrows, DiceThrow{
			Msg: HelpStr,
			Ok:  false,
		})
		return diceThrows
	}

	for _, m := range matches {
		times, err := strconv.Atoi(m[1])
		if err != nil {
			log.Println("err 1", err)
			diceThrows = append(diceThrows, DiceThrow{
				Msg: m[0] + " is not valid.",
			})
			continue
		}
		if times < 1 {
			times = 1
		} else if times > 100 {
			times = 100
		}
		maxNumber, err := strconv.Atoi(m[2])
		if err != nil || maxNumber < 1 {
			log.Println("err 2", err)
			diceThrows = append(diceThrows, DiceThrow{
				Msg: m[0] + " is not valid.",
				Ok:  false,
			})
			continue
		}

		diceThrows = append(diceThrows, DiceThrow{
			Times: times,
			Max:   maxNumber,
			Ok:    true,
		})
	}

	return diceThrows
}

func getRandom(max int) int {
	i64 := int64(max)
	v, e := rand.Int(rand.Reader, big.NewInt(i64))
	if e != nil {
		panic("Error on random generator.")
	}
	return int(v.Int64() + 1)
}
