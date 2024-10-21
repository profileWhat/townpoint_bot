package tgbot

import (
	"kate_ritson_art_bot/config"
	"kate_ritson_art_bot/pkg/components"
	"log"
	"math/rand"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type AbstractGraph struct {
	config *config.Config
	chatID int64
	api    *tgbotapi.BotAPI

	currentNode        *components.MarkupNode
	currentPath        string
	currentMessagedIDs []int
}

func NewAbstractGraph(config *config.Config, api *tgbotapi.BotAPI, chatID int64) *AbstractGraph {
	return &AbstractGraph{
		config: config,
		chatID: chatID,
		api:    api,
	}
}

func (g *AbstractGraph) Quit() {
	return
}

func (g *AbstractGraph) setAfkNode(node *components.MarkupNode) {
	g.currentNode = node
}

func (g *AbstractGraph) Start() {
	g.currentNode = g.NewAbstractGraph(g.chatID)
	msg := g.currentNode.BuildMessage(g.chatID)
	_, err := g.api.Send(msg)
	if err != nil {
		log.Printf("[Game] Error game started %d: %s", g.chatID, err)
	}
}

func (g *AbstractGraph) Continue(query *tgbotapi.CallbackQuery) {
	for _, messageID := range g.currentMessagedIDs {
		if messageID == query.Message.MessageID {
			return
		}
	}
	g.currentMessagedIDs = append(g.currentMessagedIDs, query.Message.MessageID)
	g.currentNode = g.currentNode.Continue(query, g.api)
}

func (g *AbstractGraph) LoadAbstractPicture() tgbotapi.Chattable {
	files, err := os.ReadDir(g.currentPath)
	if err != nil {
		log.Fatal(err)
	}
	cardNumber := rand.Intn(len(files) - 1)
	cardFile := files[cardNumber]

	file, _ := os.Open(g.currentPath + cardFile.Name())
	reader := tgbotapi.FileReader{Name: cardFile.Name(), Reader: file}

	photo := tgbotapi.NewPhoto(g.chatID, reader)

	return photo
}

func (g *AbstractGraph) CreateAbstractPictureNode(pic config.AbstractPicture, back *components.MarkupNode) *components.MarkupNode {
	g.currentPath = pic.Path
	// ==================================================================================================================================================
	// picture node
	pictureAnswersRow := make([]*components.MarkupAnswer, 0)
	pictureAnswersRow = append(pictureAnswersRow,
		&components.MarkupAnswer{
			Content: tgbotapi.NewInlineKeyboardButtonURL("Связаться для покупки ❤️", "https://t.me/Kate_ritson"),
		},
		&components.MarkupAnswer{
			Content: tgbotapi.NewInlineKeyboardButtonData("Назад", "Назад"),
			Next:    back,
		},
	)

	pictureAnswers := make([][]*components.MarkupAnswer, 0)
	pictureAnswers = append(pictureAnswers, pictureAnswersRow)
	pictureNode := components.NewMarkupNode(components.MarkupParams{
		Text:       pic.Description,
		Answers:    pictureAnswers,
		BeforeFunc: g.LoadAbstractPicture,
	})

	return pictureNode
}

func (g *AbstractGraph) CreateAbstractRows(chooseNode *components.MarkupNode) [][]*components.MarkupAnswer {
	rows := make([][]*components.MarkupAnswer, 0)

	for _, abstract := range g.config.Source.AbstractPictures {
		rows = append(rows, []*components.MarkupAnswer{
			{
				Content: tgbotapi.NewInlineKeyboardButtonData(abstract.Name, abstract.Name),
				Next:    g.CreateAbstractPictureNode(abstract, chooseNode),
			},
		})
	}

	return rows
}

func (g *AbstractGraph) NewAbstractGraph(chatID int64) *components.MarkupNode {
	// laterDuration := time.Hour * 24 * 3
	// afkDuration := time.Hour

	// ==================================================================================================================================================
	// start node

	startNode := components.NewMarkupNode(components.MarkupParams{
		Text: `Привет, здесь вы можете увидеть, проникнутся глубиной и купить мои картины ❤️.`,
	})

	startNode.SetAnswers(g.CreateAbstractRows(startNode))

	return startNode
}
