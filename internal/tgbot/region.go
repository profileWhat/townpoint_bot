package tgbot

import (
	"context"
	"townpoint_bot/ent/generated/region"
	"townpoint_bot/pkg/components"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gofrs/uuid"
)

func (g *TownpointGraph) AddRegionNode() *components.MarkupNode {
	AddRegionNode := components.NewMarkupNode(components.MarkupParams{
		Message: tgbotapi.NewMessage(g.chatID, "Чтобы добавить новый регион, введи имя региона:"),
		HandleFunc: func(data string) {
			g.adminActions.CreateRegionAction.isActive = true
		},
	})

	return AddRegionNode
}

func (g *TownpointGraph) ChangeRegionNode() *components.MarkupNode {
	ctx := context.Background()
	regions := g.entity.Region.Query().AllX(ctx)

	regionAnswers := make([][]*components.MarkupAnswer, 0)
	regionNode := components.NewMarkupNode(components.MarkupParams{
		Text: "Выберите регион для изменения",
	})

	for _, rg := range regions {
		markupAnswer := &components.MarkupAnswer{
			Content: tgbotapi.NewInlineKeyboardButtonData(rg.Name, rg.ID.String()),
			Next: components.NewMarkupNode(components.MarkupParams{
				Message: tgbotapi.NewMessage(g.chatID, "Введите новое название региона"),
				HandleFunc: func(data string) {
					g.adminActions.UpdateRegionAction.FieldRegionID = rg.ID
					g.adminActions.UpdateRegionAction.isActive = true
				},
			}),
		}

		regionAnswers = append(regionAnswers, []*components.MarkupAnswer{
			markupAnswer,
		})
	}

	regionNode.SetAnswers(regionAnswers)

	return regionNode
}

func (g *TownpointGraph) DeleteRegionNode() *components.MarkupNode {
	ctx := context.Background()
	regions := g.entity.Region.Query().AllX(ctx)

	regionAnswers := make([][]*components.MarkupAnswer, 0)
	regionNode := components.NewMarkupNode(components.MarkupParams{
		Text: "Выберите регион для удаления",
	})

	for _, rg := range regions {
		markupAnswer := &components.MarkupAnswer{
			Content: tgbotapi.NewInlineKeyboardButtonData(rg.Name, rg.ID.String()),
			Next: components.NewMarkupNode(components.MarkupParams{
				Message: tgbotapi.NewMessage(g.chatID, "Регион удален"),
				HandleFunc: func(data string) {
					_, err := g.entity.Region.Delete().Where(region.ID(uuid.FromStringOrNil(data))).Exec(ctx)
					if err != nil {
						g.api.Send(tgbotapi.NewMessage(g.chatID, "Ошибка удаления региона - нельзя удалить регион, если в нем есть города и точки"))
					}
				},
			}),
		}

		regionAnswers = append(regionAnswers, []*components.MarkupAnswer{
			markupAnswer,
		})
	}

	regionNode.SetAnswers(regionAnswers)

	return regionNode
}

func (g *TownpointGraph) GetRegionsNode() *components.MarkupNode {
	text := `Привет, здесь вы можете просмотреть точки различных городов. Отобранные Дамиром и его командой.

Чтобы нам с тобой выбрать точку мечты необходимо выбрать область`

	if g.isAdmin {
		text += `

Я вижу ты Админ, это сообщение видишь только ты. Ты сможешь добавлять и изменять точки, города и регионы. Удачи ;)`
	}

	regionNode := components.NewMarkupNode(components.MarkupParams{
		Text: text,
	})

	ctx := context.Background()
	regions := g.entity.Region.Query().AllX(ctx)

	regionAnswers := make([][]*components.MarkupAnswer, 0)

	if g.isAdmin {
		regionAnswers = append(regionAnswers, []*components.MarkupAnswer{
			&components.MarkupAnswer{
				Content: tgbotapi.NewInlineKeyboardButtonData("Добавить", "Добавить"),
				Next:    g.AddRegionNode(),
			},
			&components.MarkupAnswer{
				Content: tgbotapi.NewInlineKeyboardButtonData("Удалить", "Удалить"),
				Next:    g.DeleteRegionNode(),
			},
			&components.MarkupAnswer{
				Content: tgbotapi.NewInlineKeyboardButtonData("Изменить", "Изменить"),
				Next:    g.ChangeRegionNode(),
			},
		})
	}

	for _, region := range regions {
		townsNode := g.GetTownNode(regionNode, region.ID)
		markupAnswer := &components.MarkupAnswer{
			Content: tgbotapi.NewInlineKeyboardButtonData(region.Name, region.ID.String()),
		}
		markupAnswer.Next = townsNode

		regionAnswers = append(regionAnswers, []*components.MarkupAnswer{
			markupAnswer,
		})
	}

	regionNode.SetAnswers(regionAnswers)
	return regionNode
}
