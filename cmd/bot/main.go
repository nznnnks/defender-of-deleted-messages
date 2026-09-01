package main

import (
	"context"
	"log"
	"messagesdefender/database"
	"os"

	"github.com/joho/godotenv"

	tg "github.com/mr-linch/go-tg"
	"github.com/mr-linch/go-tg/tgb"
)

func main() {
	db, err := database.InitDB()
	checkErr(err)
	defer db.Close()

	err = godotenv.Load()
	checkErr(err)

	token := os.Getenv("BOT_TOKEN")

	client := tg.New(token)

	handler := tgb.HandlerFunc(func(
		ctx context.Context,
		update *tgb.Update,
	) error {

		u := update.Update

		if u.BusinessMessage != nil {
			msg := u.BusinessMessage
			log.Println("NEW MESSAGE")
			log.Printf("Message ID:%#v\n", msg.ID)
			log.Printf("Text:%#v\n", msg.Text)
		}

		if u.DeletedBusinessMessages != nil {
			log.Println("MESSAGE DELETED")
			log.Printf("%#v\n", u.DeletedBusinessMessages)
		}

		return nil
	})

	poller := tgb.NewPoller(handler, client)

	log.Println("Bot started")

	if err := poller.Run(context.Background()); err != nil {
		log.Fatal(err)
	}

}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
