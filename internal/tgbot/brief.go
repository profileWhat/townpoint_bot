package tgbot

import (
	"kate_ritson_art_bot/config"
	"kate_ritson_art_bot/pkg/components"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BriefGraph struct {
	config *config.Config
	chatID int64
	api    *tgbotapi.BotAPI

	currentNode        *components.MarkupNode
	currentMessagedIDs []int
}

func NewBriefGraph(config *config.Config, api *tgbotapi.BotAPI, chatID int64) *BriefGraph {
	return &BriefGraph{
		config: config,
		chatID: chatID,
		api:    api,
	}
}

func (g *BriefGraph) Quit() {
	return
}

func (g *BriefGraph) setAfkNode(node *components.MarkupNode) {
	g.currentNode = node
}

func (g *BriefGraph) Start() {
	g.currentNode = g.NewBriefGraph(g.chatID)
	msg := g.currentNode.BuildMessage(g.chatID)
	_, err := g.api.Send(msg)
	if err != nil {
		log.Printf("[Game] Error game started %d: %s", g.chatID, err)
	}
}

func (g *BriefGraph) Continue(query *tgbotapi.CallbackQuery) {
	for _, messageID := range g.currentMessagedIDs {
		if messageID == query.Message.MessageID {
			return
		}
	}
	g.currentMessagedIDs = append(g.currentMessagedIDs, query.Message.MessageID)
	g.currentNode = g.currentNode.Continue(query, g.api)
}

func (g *BriefGraph) NewBriefGraph(chatID int64) *components.MarkupNode {
	// laterDuration := time.Hour * 24 * 3
	// afkDuration := time.Hour

	// ==================================================================================================================================================
	// start node
	startAnswersRow := make([]*components.MarkupAnswer, 0)
	startAnswersRow = append(startAnswersRow,
		&components.MarkupAnswer{
			Content: tgbotapi.NewInlineKeyboardButtonData("Продолжим", "Продолжим"),
		},
	)

	startAnswers := make([][]*components.MarkupAnswer, 0)
	startAnswers = append(startAnswers, startAnswersRow)
	startNode := components.NewMarkupNode(components.MarkupParams{
		Text: `Привет, здесь вы можете составить заявку на разработку вашего личного логотипа.

Чтобы нам с тобой разработать логотип необходимо пройти первоначальный бриф для определения логотипа ваешй мечты.`,
		Answers: startAnswers,
	})

	return startNode
}
