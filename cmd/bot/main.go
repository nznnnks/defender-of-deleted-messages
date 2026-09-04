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

			displayName := ""

			if msg.From != nil {
				if msg.From.Username != "" {
					displayName = "@" + string(msg.From.Username)
				} else {
					displayName = msg.From.FirstName
				}
			}

			mediaType := "text"
			fileID := ""
			caption := ""

			if len(msg.Photo) > 0 {
				mediaType = "photo"

				photo := msg.Photo[len(msg.Photo)-1]

				fileID = string(photo.FileID)
				caption = msg.Caption
			}

			err := database.SaveMessage(
				db,
				int64(msg.Chat.ID),
				msg.ID,
				displayName,
				msg.Text,
				mediaType,
				fileID,
				caption,
			)

			if err != nil {
				log.Print(err)
			}
			log.Println("NEW MESSAGE")
			log.Printf("%+v\n", *msg)
			log.Printf("Message ID:%#v\n", msg.ID)
			log.Printf("Text:%#v\n", msg.Text)
		}

		//УДАЛЕНО СОО
		if u.DeletedBusinessMessages != nil {

			del := u.DeletedBusinessMessages

			for _, messageID := range del.MessageIDs {
				username, text, mediaType, fileID, capiton, err := database.GetMessage(
					db,
					int64(del.Chat.ID),
					messageID,
				)

				if err != nil {
					log.Println(err)
					continue
				}

				notification := ""

				if mediaType == "photo" {
					log.Printf("Deleted photo: %s", fileID)
					notification = fmt.Sprintf(
						"🗑️ Фото удалено пользователем %s\n\n%s",
						username,
						capiton,
					)

					err = client.SendMessage(
						tg.ChatID(chatID),
						notification,
					).DoVoid(ctx)

					if err != nil {
						log.Println(err)
					}

					err = client.SendPhoto(
						tg.ChatID(chatID),
						tg.FileArg{
							FileID: tg.FileID(fileID),
						},
					).DoVoid(ctx)

				} else {
					notification = fmt.Sprintf(
						"🗑️ Сообщение удалено пользователем %s\n\n%s",
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
				}

				log.Printf("Deleted message: %s", text)
			}
		}

		//ИСПРАВЛЕНО СОО
		if u.EditedBusinessMessage != nil {

			msg := u.EditedBusinessMessage

			username, oldText, mediaType, fileID, capiton, err := database.GetMessage(
				db,
				int64(msg.Chat.ID),
				msg.ID,
			)

			if err != nil {
				log.Println(err)
				return nil
			}

			newText := msg.Text

			notification := ""

			if oldText != newText {

				if mediaType == "photo" {
					log.Printf("Deleted photo: %s", fileID)
					notification = fmt.Sprintf(
						"🗑️ Фото изменено пользователем %s\n\n%s",
						username,
						capiton,
					)

					err = client.SendMessage(
						tg.ChatID(chatID),
						notification,
					).DoVoid(ctx)

					if err != nil {
						log.Println(err)
					}

					err = client.SendPhoto(
						tg.ChatID(chatID),
						tg.FileArg{
							FileID: tg.FileID(fileID),
						},
					).DoVoid(ctx)

					if err != nil {
						log.Println(err)
					}
				} else {
					notification = fmt.Sprintf("✏️ Сообщение изменено пользователем %s\n\n%s\n\n%s",
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
