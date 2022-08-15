package vuploader

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramUploader struct {
	Bot    *tgbotapi.BotAPI
	ChatId int64
}

func NewTelegramUploader(token string, chatId int64) (*TelegramUploader, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	bot.Debug = true
	return &TelegramUploader{
		Bot:    bot,
		ChatId: chatId,
	}, nil
}
func (t *TelegramUploader) ChatIdGetter() error {
	updateConfig := tgbotapi.NewUpdate(0)

	// Tell Telegram we should wait up to 30 seconds on each request for an
	// update. This way we can get information just as quickly as making many
	// frequent requests without having to send nearly as many.
	updateConfig.Timeout = 30

	// Start polling Telegram for updates.
	updates := t.Bot.GetUpdatesChan(updateConfig)

	// Let's go through each update that we're getting from Telegram.
	for update := range updates {
		// Telegram can send many types of updates depending on what your Bot
		// is up to. We only want to look at messages for now, so we can
		// discard any other updates.
		if update.Message == nil {
			continue
		}
		if strings.Contains(update.Message.Text, "/start") {
			// Now that we know we've gotten a new message, we can construct a
			// reply! We'll take the Chat ID and Text from the incoming message
			// and use it to create a new message.
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("%d", update.Message.Chat.ID))
			// We'll also say that this message is a reply to the previous message.
			// For any other specifications than Chat ID or Text, you'll need to
			// set fields on the `MessageConfig`.
			msg.ReplyToMessageID = update.Message.MessageID

			// Okay, we're sending our message off! We don't care about the message
			// we just sent, so we'll discard it.
			if _, err := t.Bot.Send(msg); err != nil {
				// Note that panics are a bad way to handle errors. Telegram can
				// have service outages or network errors, you should retry sending
				// messages or more gracefully handle failures.
				return err
			}
		}

	}
	return nil
}
func (t *TelegramUploader) SendMessage(message string) error {
	msg := tgbotapi.NewMessage(t.ChatId, message)
	msg.ParseMode = tgbotapi.ModeHTML
	sended, err := t.Bot.Send(msg)
	if err != nil {
		return err
	}
	t.printMessage(sended)
	return nil
}
func (t *TelegramUploader) UploadFiles(filePathes []string, message string) error {
	files := []interface{}{}
	l := len(filePathes)
	for index, path := range filePathes {
		file := tgbotapi.NewInputMediaDocument(tgbotapi.FilePath(path))
		if l == index+1 {
			message = tgbotapi.EscapeText(tgbotapi.ModeMarkdown, message)
			file.ParseMode = tgbotapi.ModeHTML
			file.Caption = message
		}
		files = append(files, file)
	}
	fileMessage := tgbotapi.NewMediaGroup(t.ChatId, files)
	sended, err := t.Bot.Send(fileMessage)
	if err != nil {
		return err
	}
	t.printMessage(sended)
	return nil
}

func (t *TelegramUploader) printMessage(message interface{}) {
	bytes, err := json.MarshalIndent(message, "", " ")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Message sent: %s", string(bytes))
}
