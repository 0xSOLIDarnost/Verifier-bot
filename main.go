package main

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var yesNoKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Yes"),
		tgbotapi.NewKeyboardButton("No")),
)

var mainKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Verify personal wallet"),
		tgbotapi.NewKeyboardButton("Verify collective wallet")),
)

//to operate the bot, put a text file containing key for your bot acquired from telegram "botfather" to the same directory with this file
var tgApiKey, err = os.ReadFile(".secret")
var bot, error1 = tgbotapi.NewBotAPI(string(tgApiKey))

//type containing all the info about user input
type user struct {
	tgid                int64
	status            int64
	//exportTokenName   string
	//exportTokenSymbol string
	//exportTokenSupply uint64
	//exportTokenType   uint64
	//tokenTypeString   string
}

//main database, key (int64) is telegram user id
var userDatabase = make(map[int64]user)

var msgTemplates = make (map[string] string)


func main() {


	msgTemplates["hello"] = "Hey, this bot is attaching personal wallets to telegram user & collective wallets to chat id"
	msgTemplates["case0"] = "Go to link and attach your tg_id to your metamask wallet"

	bot, err = tgbotapi.NewBotAPI(string(tgApiKey))
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	//whenever bot gets a new message, check for user id in the database happens, if it's a new user, the entry in the database is created.
	for update := range updates {

		if update.Message != nil {
			if _, ok := userDatabase[update.Message.From.ID]; !ok {

				userDatabase[update.Message.From.ID] = user{update.Message.Chat.ID, 0,}
				msg := tgbotapi.NewMessage(userDatabase[update.Message.From.ID].tgid, msgTemplates["hello"])
				msg.ReplyMarkup = mainKeyboard
				bot.Send(msg)
			} else {

				switch userDatabase[update.Message.From.ID].status {

				//first check for user status, (for a new user status 0 is set automatically), then user reply for the first bot message is logged to a database as name AND user status is updated
				case 0:
					if updateDb, ok := userDatabase[update.Message.From.ID]; ok {
						//updateDb.exportTokenName = update.Message.Text
						updateDb.status = 1
						userDatabase[update.Message.From.ID] = updateDb

						msg := tgbotapi.NewMessage(userDatabase[update.Message.From.ID].tgid,msgTemplates["case0"] )
						bot.Send(msg)
						
					}

				//logic is that 1 incoming message fro the user equals one status check in database, so each status check ends with the message asking the next question
				case 1:
					if updateDb, ok := userDatabase[update.Message.From.ID]; ok {
						//updateDb.exportTokenSymbol = update.Message.Text
						updateDb.status = 2
						userDatabase[update.Message.From.ID] = updateDb
					}
					msg := tgbotapi.NewMessage(userDatabase[update.Message.From.ID].tgid, msgTemplates["hello"])
					bot.Send(msg)


				}
			}
		}
	}
}