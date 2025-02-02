package tgbot

import (
	"log"
	"townpoint_bot/config"
	"townpoint_bot/internal/interfaces"
	"townpoint_bot/internal/services"

	ent "townpoint_bot/ent/generated"

	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TGBot struct {
	interfaces.Service

	config *config.Config
	quit   chan bool
	api    *tgbotapi.BotAPI

	entity *ent.Client
	yadisk *services.Yadisk

	sessions map[int64]*TownpointGraph
}

func New(config *config.Config, entity *ent.Client, yadisk *services.Yadisk) *TGBot {
	bot, err := tgbotapi.NewBotAPI(config.TGbot.Token)
	if err != nil {
		// Abort if something is wrong
		log.Panic(err)
	}

	return &TGBot{
		api:      bot,
		quit:     make(chan bool, 1),
		sessions: make(map[int64]*TownpointGraph),
		entity:   entity,
		yadisk:   yadisk,

		config: config,
	}
}

func (b *TGBot) Stop() error {
	b.quit <- true
	return nil
}

func (b *TGBot) Start() error {
	u := tgbotapi.NewUpdate(0)

	// `updates` is a golang channel which receives telegram updates
	updates := b.api.GetUpdatesChan(u)

	for {
		select {
		// stop looping if ctx is cancelled
		case <-b.quit:
			return nil
		// receive update from channel and then handle it
		case update := <-updates:
			b.handleUpdate(update)
		}
	}
}

func (b *TGBot) handleUpdate(update tgbotapi.Update) {
	switch {
	// Handle messages
	case update.Message != nil:
		b.handleMessage(update.Message)

	// Handle button clicks
	case update.CallbackQuery != nil:
		chatID := update.CallbackQuery.Message.Chat.ID
		session, ok := b.sessions[chatID]
		if ok {
			session.Continue(update.CallbackQuery)
		}
	}
}

func (b *TGBot) handleMessage(message *tgbotapi.Message) {
	user := message.From
	text := message.Text

	if user == nil {
		return
	}

	// Print to console
	log.Printf("%s wrote %s", user.UserName, text)

	var err error
	if strings.HasPrefix(text, "/") {
		err = b.handleCommand(user.UserName, message.Chat.ID, text)
	} else {
		session, ok := b.sessions[message.Chat.ID]
		if ok {
			session.ContinueM(message)
		}
	}

	if err != nil {
		log.Printf("An error occured: %s", err.Error())
	}
}

// When we get a command, we react accordingly
func (b *TGBot) handleCommand(username string, chatId int64, command string) error {
	var err error

	switch command {
	case "/start":
		session, ok := b.sessions[chatId]
		if ok {
			session.Quit()
		}

		b.sessions[chatId] = NewTownpointGraph(b.config, b.api, b.entity, b.yadisk, chatId, username)
		b.sessions[chatId].Start()
	case "/about":
		text := ``

		msg := tgbotapi.NewMessage(chatId, text)

		b.api.Send(msg)

	}

	return err
}
