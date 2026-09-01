package main

import (
	"context"
	"fmt"
	"log"
	"messagesdefender/database"
	"os"
	"strconv"

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

		if u.Message != nil {
			log.Printf(
				"MY CHAT ID: %d\n",
				u.Message.Chat.ID,
			)
		}

		chatID, err := strconv.ParseInt(
			os.Getenv("CHAT_ID"),
			10,
			64,
		)

		checkErr(err)

		if u.BusinessMessage != nil {
			msg := u.BusinessMessage

			//ОТДЕЛЬНАЯ ФУНКЦИЯ
			displayName := ""

			if msg.From != nil {
				if msg.From.Username != "" {
					displayName = "@" + string(msg.From.Username)
				} else {
					displayName = msg.From.FirstName
				}
			}

			err := database.SaveMessage(
				db,
				int64(msg.Chat.ID),
				msg.ID,
				displayName,
				msg.Text,
			)

			if err != nil {
				log.Print(err)
			}
			log.Println("NEW MESSAGE")
			log.Printf("Message ID:%#v\n", msg.ID)
			log.Printf("Text:%#v\n", msg.Text)
		}

		if u.DeletedBusinessMessages != nil {

			del := u.DeletedBusinessMessages

			for _, messageID := range del.MessageIDs {
				username, text, err := database.GetMessage(
					db,
					int64(del.Chat.ID),
					messageID,
				)

				if err != nil {
					log.Println(err)
					continue
				}

				notification := fmt.Sprintf(
					"🗑️ Сообщение было удалено пользователем %s\n\nТекст:\n%s",
					username,
					text,
				)

				err = client.SendMessage(
					tg.ChatID(chatID),
					notification,
				).DoVoid(ctx)

				if err != nil {
					log.Println(err)
				}

				log.Printf("Deleted message: %s", text)
			}
		}

		if u.EditedBusinessMessage != nil {

			msg := u.EditedBusinessMessage
			username := msg.From.Username

			if username == "" {
				username = tg.Username(msg.From.FirstName)
			}

			//ИСПРАВИТЬ ЮЗЕРНЕЙМ
			oldText, err := database.GetMessage(
				db,
				int64(msg.Chat.ID),
				msg.ID,
			)
			if err != nil {
				log.Println(err)
				return nil
			}

			newText := msg.Text

			if oldText != newText {
				notification := fmt.Sprintf("✏️ Сообщение изменено пользователем %s\n\nБыло:\n%s\n\nСтало\n%s",
					username,
					oldText,
					newText,
				)

				err = client.SendMessage(
					tg.ChatID(chatID),
					notification,
				).DoVoid(ctx)

				if err != nil {
					log.Println(err)
				}

				err = database.UpdateMessage(
					db,
					int64(msg.Chat.ID),
					msg.ID,
					newText,
				)

				if err != nil {
					log.Println(err)
				}

			}
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
