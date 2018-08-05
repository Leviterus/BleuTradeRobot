package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/yanzay/tbot"
)

var apikey, apisecret, signkey string

func main() {
	bot, err := tbot.NewServer(os.Getenv("TELEGRAM_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}
	bot.HandleFunc("/api {key}", keyHandler)
	bot.HandleFunc("/secret {key}", secretHandler)
	bot.HandleFunc("/hi", HiHandler)
	bot.HandleFunc("/balance", balanceHandler)
	bot.HandleFunc("/withdraw {currency} {quantity} {address}", withdrawHandler)
	bot.HandleFunc("/tradebuy {market} {rate} {quantity}", tradeBuyHandler)
	bot.HandleFunc("/tradesell {market} {rate} {quantity}", tradeSellHandler)

	bot.ListenAndServe()
}

func HiHandler(message *tbot.Message) {
	// Handler can reply with several messages
	message.Replyf("Hello, %s!", message.From.FirstName)
	time.Sleep(1 * time.Second)
	message.Reply("What's up?")
	time.Sleep(1 * time.Second)
	message.Reply("Type /help to see my available commands.")

}

func keyHandler(m *tbot.Message) {
	// m.Vars contains all variables, parsed during routing
	apikey = m.Vars["key"]

	fmt.Println(apikey)
	// Convert string variable to integer seconds value
	// seconds, err := strconv.Atoi(apikey)
	// if err != nil {
	// 	m.Reply("Invalid API Key")
	// 	return
	// }

}

func secretHandler(m *tbot.Message) {
	// m.Vars contains all variables, parsed during routing
	apisecret = m.Vars["key"]

	fmt.Println(apisecret)
	// Convert string variable to integer seconds value
	// seconds, err := strconv.Atoi(apikey)
	// if err != nil {
	// 	m.Reply("Invalid API Key")
	// 	return
	// }

}
