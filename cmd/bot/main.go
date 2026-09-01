package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"

	tg "github.com/mr-linch/go-tg"
	"github.com/mr-linch/go-tg/tgb"
)

func main() {
	err := godotenv.Load()
	if err != nil{
		log.Fatal(err)
	}

	token := os.Getenv("BOT_TOKEN")

	client := tg.New(token)

	handler := tgb.HandlerFunc(func(
		ctx context.Context,
		update *tgb.Update,
	) error{
		log.Printf("%#v\n", update.Update)
		return nil
	})

	poller := tgb.NewPoller(handler, client)

	log.Println("Bot started")

	if err := poller.Run(context.Background()); err != nil{
		log.Fatal(err)
	}

}