package components

import (
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ContentFunc func() tgbotapi.Chattable
type HandleFunc func(data string)
type HandleAfkFunc func(node *MarkupNode)

type MarkupAnswer struct {
	Reflect bool
	Next    *MarkupNode
	Content tgbotapi.InlineKeyboardButton
}

type MarkupNode struct {
	text    string
	answers [][]*MarkupAnswer

	remindLater    *time.Duration
	remindLaterMsg tgbotapi.Chattable

	remindIfAFK     *time.Duration
	remindIfAFKText string

	next    *MarkupNode
	prev    *MarkupNode
	message tgbotapi.Chattable

	maxRepeatNumber int
	afterRepeat     *MarkupNode

	afterFunc  ContentFunc
	handleFunc HandleFunc
	beforeFunc ContentFunc
	afkFunc    HandleAfkFunc

	// iternal logic fields
	remindAFKquit chan bool
	quitAFK       chan bool
	quitLater     chan bool
	repeatNumber  int
}

func (n *MarkupNode) SetAnswers(answers [][]*MarkupAnswer) {
	n.answers = answers
}

func (n *MarkupNode) GetAnswers() [][]*MarkupAnswer {
	return n.answers
}

func (n *MarkupNode) Quit() {
	n.quitAFK <- true
	n.quitLater <- true
}

func (n *MarkupNode) BuildMessage(chatID int64) tgbotapi.Chattable {
	msg := tgbotapi.NewMessage(chatID, n.text)
	msg.ReplyMarkup = n.Build()
	msg.ParseMode = tgbotapi.ModeHTML
	return msg
}

func CreateAFKNode(current *MarkupNode, text string) *MarkupNode {
	afkAnswersRow := make([]*MarkupAnswer, 0)
	afkAnswersRow = append(afkAnswersRow,
		&MarkupAnswer{
			Content: tgbotapi.NewInlineKeyboardButtonData("Да продолжаем", "Да продолжаем"),
			Next:    current,
		},
		&MarkupAnswer{
			Content: tgbotapi.NewInlineKeyboardButtonURL("Нужна помощь", "https://t.me/ev_do_live"),
		},
		&MarkupAnswer{
			Content: tgbotapi.NewInlineKeyboardButtonURL("Не интересно", "https://docs.google.com/forms/d/e/1FAIpQLSeyx3buP40jI2FdDxjqvb0bvVCx5NST0GixZKd1d3ShCl8Pyw/viewform?usp=sf_link"),
		},
	)

	afkAnswers := make([][]*MarkupAnswer, 0)
	afkAnswers = append(afkAnswers, afkAnswersRow)
	return &MarkupNode{
		text: `Похоже, что мы остановились в процессе пути. 
Получается ли у тебя найти ответы?`,
		answers: afkAnswers,
	}
}

type MarkupParams struct {
	Text    string
	Answers [][]*MarkupAnswer

	RemindLater    *time.Duration
	RemindLaterMsg tgbotapi.Chattable

	RemindIfAFK     *time.Duration
	RemindIfAFKText string

	Next    *MarkupNode
	Prev    *MarkupNode
	Message tgbotapi.Chattable

	AfterFunc  ContentFunc
	HandleFunc HandleFunc
	BeforeFunc ContentFunc
	AfkFunc    HandleAfkFunc

	MaxRepeatNumber int
	AfterRepeat     *MarkupNode
}

func NewMarkupNode(params MarkupParams) *MarkupNode {
	return &MarkupNode{
		text:            params.Text,
		answers:         params.Answers,
		remindLater:     params.RemindLater,
		remindLaterMsg:  params.RemindLaterMsg,
		remindIfAFK:     params.RemindIfAFK,
		remindIfAFKText: params.RemindIfAFKText,
		next:            params.Next,
		prev:            params.Prev,
		message:         params.Message,
		maxRepeatNumber: params.MaxRepeatNumber,
		afterRepeat:     params.AfterRepeat,

		afterFunc:  params.AfterFunc,
		handleFunc: params.HandleFunc,
		beforeFunc: params.BeforeFunc,
		afkFunc:    params.AfkFunc,

		remindAFKquit: make(chan bool, 1),
		quitAFK:       make(chan bool, 1),
		quitLater:     make(chan bool, 1),
		repeatNumber:  1,
	}
}

func (m *MarkupNode) Build() tgbotapi.InlineKeyboardMarkup {
	rows := make([][]tgbotapi.InlineKeyboardButton, 0)
	for _, keyboardRow := range m.answers {
		row := make([]tgbotapi.InlineKeyboardButton, 0)
		for _, answer := range keyboardRow {
			row = append(row, answer.Content)
		}
		rows = append(rows, row)
	}

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

func (m *MarkupNode) Continue(query *tgbotapi.CallbackQuery, api *tgbotapi.BotAPI) *MarkupNode {
	message := query.Message

	// prev node internal logic
	if m.prev != nil && m.prev.remindIfAFK != nil {
		m.prev.remindAFKquit <- true
	}

	// iternal node logic
	m.repeatNumber++

	// getting next node
	isReflect := false
	var nextNode *MarkupNode

answersLoop:
	for _, row := range m.answers {
		for _, answer := range row {
			if answer.Content.CallbackData != nil && query.Data == *answer.Content.CallbackData {
				nextNode = answer.Next
				if answer.Reflect {
					isReflect = true
				}
				break answersLoop
			}
		}
	}

	if isReflect {
		nextNode = m
	}

	if nextNode == nil {
		nextNode = m.next
	}

	if nextNode == nil {
		return nil
	}
	nextNode.prev = m

	if nextNode.maxRepeatNumber != 0 && nextNode.maxRepeatNumber < nextNode.repeatNumber {
		nextNode = nextNode.afterRepeat
	}

	// iternal logic before send message
	if m.handleFunc != nil {
		m.handleFunc(query.Data)
	}

	var beforeMsg tgbotapi.Chattable
	if nextNode.beforeFunc != nil {
		beforeMsg = nextNode.beforeFunc()
	}

	var afterMsg tgbotapi.Chattable
	if m.afterFunc != nil {
		afterMsg = m.afterFunc()
	}

	// creating message on markup or primary message
	var iMsg tgbotapi.Chattable

	if nextNode.answers != nil {
		iMsg = nextNode.BuildMessage(message.Chat.ID)
	}

	if nextNode.message != nil {
		iMsg = nextNode.message
	}

	if iMsg == nil {
		return nil
	}

	//sending message with internal logic
	if nextNode.remindLater != nil {
		api.Send(nextNode.remindLaterMsg)
		timer := time.NewTimer(*nextNode.remindLater)
		go func() {
			for {
				select {
				case <-timer.C:
					api.Send(iMsg)
					return
				case <-m.quitLater:
					return
				}
			}
		}()
		return nextNode
	}

	if m.remindIfAFK != nil {
		go func() {
			timer := time.NewTimer(*m.remindIfAFK)
			for {
				select {
				case <-m.quitAFK:
					return
				case <-timer.C:
					if beforeMsg != nil {
						api.Send(beforeMsg)
					}
					nextNode = CreateAFKNode(nextNode, nextNode.remindIfAFKText)
					api.Send(nextNode.BuildMessage(message.Chat.ID))
					m.afkFunc(nextNode)
					return
				case <-m.remindAFKquit:
					return
				}
			}
		}()
	}

	if afterMsg != nil {
		_, err := api.Send(afterMsg)
		if err != nil {
			log.Println(err)
		}
	}

	if beforeMsg != nil {
		_, err := api.Send(beforeMsg)
		if err != nil {
			log.Println(err)
		}
	}

	_, err := api.Send(iMsg)
	if err != nil {
		log.Println(err)
	}

	return nextNode
}
