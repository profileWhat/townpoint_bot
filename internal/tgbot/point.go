package tgbot

import (
	"context"
	"fmt"
	"townpoint_bot/ent/generated/point"
	"townpoint_bot/pkg/components"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gofrs/uuid"
)

func (g *TownpointGraph) GetDetailedPoint(pointNode *components.MarkupNode, pointID uuid.UUID) *components.MarkupNode {
	ctx := context.Background()
	point := g.entity.Point.Query().Where(point.ID(pointID)).OnlyX(ctx)

	pointAnswersRow := make([]*components.MarkupAnswer, 0)
	markupAnswer := &components.MarkupAnswer{
		Content: tgbotapi.NewInlineKeyboardButtonData("Назад", "Назад"),
		Next:    pointNode,
	}
	pointAnswersRow = append(pointAnswersRow, markupAnswer)

	pointAnswers := make([][]*components.MarkupAnswer, 0)
	pointAnswers = append(pointAnswers, pointAnswersRow)

	detailedPointNode := components.NewMarkupNode(components.MarkupParams{
		Text:    fmt.Sprintf("<b>Информация о точке</b>\n\n<b>Наименование:</b> %s\n<b>Адресс:</b> %s\n<b>Описание:</b> %s\n<b>Телефон:</b> %s\n", point.Name, point.Address, *point.Description, point.Phone),
		Answers: pointAnswers,
	})

	if len(point.Videos) > 0 {
		detailedPointNode.SetBeforeFunc(func() tgbotapi.Chattable {
			// fVideoPath := point.Videos[0]

			return tgbotapi.NewMessage(g.chatID, point.Videos[0].Path)
		})
	}

	return detailedPointNode
}

func (g *TownpointGraph) AddPointNode(townID uuid.UUID) *components.MarkupNode {
	AddRegionNode := components.NewMarkupNode(components.MarkupParams{
		Message: tgbotapi.NewMessage(g.chatID, "Чтобы добавить новую точку, введи имя точки:"),
		HandleFunc: func(data string) {
			g.adminActions.CreatePointAction.isActive = true
			g.adminActions.CreatePointAction.FieldTownID = townID
		},
	})

	return AddRegionNode
}

func (g *TownpointGraph) ChangeChooseFieldNode(pointID uuid.UUID) *components.MarkupNode {
	answers := make([][]*components.MarkupAnswer, 0)
	node := components.NewMarkupNode(components.MarkupParams{
		Text: "Выберите поле изменения",
	})

	nameAnswer := []*components.MarkupAnswer{
		&components.MarkupAnswer{
			Content: tgbotapi.NewInlineKeyboardButtonData("Имя точки", "Имя точки"),
			Next: components.NewMarkupNode(components.MarkupParams{
				Message: tgbotapi.NewMessage(g.chatID, "Введите новые данные"),
				HandleFunc: func(data string) {
					g.adminActions.UpdatePointAction.FieldChangeName = true
					g.adminActions.UpdatePointAction.isActive = true
					g.adminActions.UpdatePointAction.FieldPointID = pointID
				},
			}),
		},
	}
	answers = append(answers, nameAnswer)

	addressAnswer := []*components.MarkupAnswer{
		&components.MarkupAnswer{
			Content: tgbotapi.NewInlineKeyboardButtonData("Адрес точки", "Адрес точки"),
			Next: components.NewMarkupNode(components.MarkupParams{
				Message: tgbotapi.NewMessage(g.chatID, "Введите новые данные"),
				HandleFunc: func(data string) {
					g.adminActions.UpdatePointAction.FieldChangeAddress = true
					g.adminActions.UpdatePointAction.isActive = true
					g.adminActions.UpdatePointAction.FieldPointID = pointID
				},
			}),
		},
	}
	answers = append(answers, addressAnswer)

	descriptionAnswer := []*components.MarkupAnswer{
		&components.MarkupAnswer{
			Content: tgbotapi.NewInlineKeyboardButtonData("Описание точки", "Описание точки"),
			Next: components.NewMarkupNode(components.MarkupParams{
				Message: tgbotapi.NewMessage(g.chatID, "Введите новые данные"),
				HandleFunc: func(data string) {
					g.adminActions.UpdatePointAction.FieldChangeDescription = true
					g.adminActions.UpdatePointAction.isActive = true
					g.adminActions.UpdatePointAction.FieldPointID = pointID
				},
			}),
		},
	}
	answers = append(answers, descriptionAnswer)

	contactAnswer := []*components.MarkupAnswer{
		&components.MarkupAnswer{
			Content: tgbotapi.NewInlineKeyboardButtonData("Контакты точки", "Контакты точки"),
			Next: components.NewMarkupNode(components.MarkupParams{
				Message: tgbotapi.NewMessage(g.chatID, "Введите новые данные"),
				HandleFunc: func(data string) {
					g.adminActions.UpdatePointAction.FieldChangePhone = true
					g.adminActions.UpdatePointAction.isActive = true
					g.adminActions.UpdatePointAction.FieldPointID = pointID
				},
			}),
		},
	}
	answers = append(answers, contactAnswer)

	videoAnswer := []*components.MarkupAnswer{
		&components.MarkupAnswer{
			Content: tgbotapi.NewInlineKeyboardButtonData("Видео точки", "Видео точки"),
			Next: components.NewMarkupNode(components.MarkupParams{
				Message: tgbotapi.NewMessage(g.chatID, "Загрузите новое видео"),
				HandleFunc: func(data string) {
					g.adminActions.UpdatePointAction.FieldChangeVideo = true
					g.adminActions.UpdatePointAction.isActive = true
					g.adminActions.UpdatePointAction.FieldPointID = pointID
				},
			}),
		},
	}
	answers = append(answers, videoAnswer)

	node.SetAnswers(answers)
	return node
}

func (g *TownpointGraph) ChangePointNode(townID uuid.UUID) *components.MarkupNode {
	ctx := context.Background()
	points := g.entity.Point.Query().Where(point.TownID(townID)).AllX(ctx)

	pointAnswers := make([][]*components.MarkupAnswer, 0)
	pointNode := components.NewMarkupNode(components.MarkupParams{
		Text: "Выберите точку для изменения",
	})

	for _, pt := range points {
		markupAnswer := &components.MarkupAnswer{
			Content: tgbotapi.NewInlineKeyboardButtonData(pt.Name, pt.ID.String()),
			Next:    g.ChangeChooseFieldNode(pt.ID),
		}

		pointAnswers = append(pointAnswers, []*components.MarkupAnswer{
			markupAnswer,
		})
	}

	pointNode.SetAnswers(pointAnswers)

	return pointNode
}

func (g *TownpointGraph) DeletePointNode(townID uuid.UUID) *components.MarkupNode {
	ctx := context.Background()
	points := g.entity.Point.Query().Where(point.TownID(townID)).AllX(ctx)

	pointAnswers := make([][]*components.MarkupAnswer, 0)
	pointNode := components.NewMarkupNode(components.MarkupParams{
		Text: "Выберите точку для удаления",
	})

	for _, pn := range points {
		markupAnswer := &components.MarkupAnswer{
			Content: tgbotapi.NewInlineKeyboardButtonData(pn.Name, pn.ID.String()),
			Next: components.NewMarkupNode(components.MarkupParams{
				Message: tgbotapi.NewMessage(g.chatID, "Город или поселок удален"),
				HandleFunc: func(data string) {
					_, err := g.entity.Point.Delete().Where(point.ID(uuid.FromStringOrNil(data))).Exec(ctx)
					if err != nil {
						g.api.Send(tgbotapi.NewMessage(g.chatID, "Ошибка удаления точки: "+err.Error()))
					}
				},
			}),
		}

		pointAnswers = append(pointAnswers, []*components.MarkupAnswer{
			markupAnswer,
		})
	}

	pointNode.SetAnswers(pointAnswers)

	return pointNode
}

func (g *TownpointGraph) GetPointNode(townNode *components.MarkupNode, townID uuid.UUID) *components.MarkupNode {
	pointNode := components.NewMarkupNode(components.MarkupParams{
		Text: `Далее необходимо выбрать точку`,
	})

	ctx := context.Background()
	points := g.entity.Point.Query().Where(point.TownID(townID)).AllX(ctx)

	pointAnswers := make([][]*components.MarkupAnswer, 0)

	if g.isAdmin {
		pointAnswers = append(pointAnswers, []*components.MarkupAnswer{
			&components.MarkupAnswer{
				Content: tgbotapi.NewInlineKeyboardButtonData("Добавить", "Добавить"),
				Next:    g.AddPointNode(townID),
			},
			&components.MarkupAnswer{
				Content: tgbotapi.NewInlineKeyboardButtonData("Удалить", "Удалить"),
				Next:    g.DeletePointNode(townID),
			},
			&components.MarkupAnswer{
				Content: tgbotapi.NewInlineKeyboardButtonData("Изменить", "Изменить"),
				Next:    g.ChangePointNode(townID),
			},
		})
	}

	for _, point := range points {
		detailedPointNode := g.GetDetailedPoint(pointNode, point.ID)
		markupAnswer := &components.MarkupAnswer{
			Content: tgbotapi.NewInlineKeyboardButtonData(point.Name, point.ID.String()),
		}
		markupAnswer.Next = detailedPointNode

		pointAnswers = append(pointAnswers, []*components.MarkupAnswer{
			markupAnswer,
		})
	}
	pointAnswers = append(pointAnswers, []*components.MarkupAnswer{
		&components.MarkupAnswer{
			Content: tgbotapi.NewInlineKeyboardButtonData("Назад", "Назад"),
			Next:    townNode,
		}})

	pointNode.SetAnswers(pointAnswers)

	return pointNode
}
