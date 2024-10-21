package tgbot

import (
	"kate_ritson_art_bot/config"
	"kate_ritson_art_bot/internal/interfaces"
	"log"

	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TGBot struct {
	interfaces.Service

	config *config.Config
	quit   chan bool
	api    *tgbotapi.BotAPI

	briefs    map[int64]*BriefGraph
	abstracts map[int64]*AbstractGraph
}

func New(config *config.Config) *TGBot {
	bot, err := tgbotapi.NewBotAPI(config.TGbot.Token)
	if err != nil {
		// Abort if something is wrong
		log.Panic(err)
	}

	return &TGBot{
		api:       bot,
		quit:      make(chan bool, 1),
		briefs:    make(map[int64]*BriefGraph),
		abstracts: make(map[int64]*AbstractGraph),

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
		brief, ok := b.briefs[chatID]
		if ok {
			brief.Continue(update.CallbackQuery)
		}

		abstr, ok := b.abstracts[chatID]
		if ok {
			abstr.Continue(update.CallbackQuery)
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
		err = b.handleCommand(message.Chat.ID, text)
	}

	if err != nil {
		log.Printf("An error occured: %s", err.Error())
	}
}

// When we get a command, we react accordingly
func (b *TGBot) handleCommand(chatId int64, command string) error {
	var err error

	switch command {
	case "/abstract":
		abstr, ok := b.abstracts[chatId]
		if ok {
			abstr.Quit()
		}

		b.abstracts[chatId] = NewAbstractGraph(b.config, b.api, chatId)
		b.abstracts[chatId].Start()
	case "/brief":
		game, ok := b.briefs[chatId]
		if ok {
			game.Quit()
		}

		b.briefs[chatId] = NewBriefGraph(b.config, b.api, chatId)
		b.briefs[chatId].Start()
	case "/about":
		text := `Лучшая художница Kate Ritson 🌺

Занимаюсь абстрактным искусством, наполняю картины самыми яркими чувствами.
Мой тг канал:
https://t.me/kateritson

Разработка логотипов совместно с клиентом, любые капризы за ваши деньги 

Мой инстаграм:
https://www.instagram.com/kate.ritson.art?igsh=czBtNWQ2aTluNm04


`

		msg := tgbotapi.NewMessage(chatId, text)

		b.api.Send(msg)

	}

	return err
}
