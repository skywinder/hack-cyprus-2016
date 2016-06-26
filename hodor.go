package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	//"github.com/mrd0ll4r/tbotapi"
	//"github.com/mrd0ll4r/tbotapi/examples/boilerplate"
)

func main() {
	apiToken := os.Getenv("HODOR_TOKEN")
	if len(apiToken) == 0 {
		fmt.Println("MAIN", "You need to set HODOR_TOKEN environment variable")
		os.Exit(1)
	}

	bot, err := tgbotapi.NewBotAPI(apiToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	var a chan string = make(chan string)
	var b chan int64 = make(chan int64)

	go tesselListener(a)
	go telegramListener(bot, a, b)
	go telegramSender(bot, a, b)

	var input string
	fmt.Scanln(&input)

	// fmt.Printf(sayHodor(rand.Intn(5)))
}
func sayHodor(times int) string {
	var out = "HODOR"
	for i := 0; i < times; i++ {
		out += " HODOR"
	}
	return out
}

func tesselListener(a chan string) {
	ln, err := net.Listen("tcp", ":8088")
	if err != nil {
		fmt.Println(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
		}
		// handleTesselRequest(net.Conn)
		a <- "ping"
		conn.Close()
	}
}

// func handleTesselRequest(conn net.Conn) {
//   // Make a buffer to hold incoming data.
//   buf := make([]byte, 1024)
//   // Read the incoming connection into the buffer.
//   reqLen, err := conn.Read(buf)
//   if err != nil {
//     fmt.Println("Error reading:", err.Error())
//   }
//   // Send a response back to person contacting us.
//   conn.Write([]byte("Message received."))
//   // Close the connection when you're done with it.
//   conn.Close()
// }

func telegramListener(bot *tgbotapi.BotAPI, a chan string, b chan int64) {

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		fmt.Println(err)
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, sayHodor(rand.Intn(5)))
		msg.ReplyToMessageID = update.Message.MessageID

		bot.Send(msg)
		b <- update.Message.Chat.ID
	}
}
func telegramSender(bot *tgbotapi.BotAPI, a chan string, b chan int64) {
	var m = make(map[int64]int64)
	for {
		select {
		case CommandToSend := <-a:
			fmt.Println("CommandToSend=", CommandToSend)
			for key, val := range m {
				msg := tgbotapi.NewMessage(key, sayHodor(rand.Intn(5)))
				//msg := tgbotapi.NewMessage(key, fmt.Sprintf("Key= %s", key))

				bot.Send(msg)

				fmt.Println("Send msg to Chat Id=", key, "val=", val)
			}
		case NewChatId := <-b:
			fmt.Println("NewChatID", NewChatId)
			m[NewChatId] = NewChatId
		}
	}
}
