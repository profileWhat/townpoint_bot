package tgbot

import (
	"context"
	"townpoint_bot/ent/generated/town"
	"townpoint_bot/pkg/components"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gofrs/uuid"
)

func (g *TownpointGraph) AddTownNode(regionID uuid.UUID) *components.MarkupNode {
	AddRegionNode := components.NewMarkupNode(components.MarkupParams{
		Message: tgbotapi.NewMessage(g.chatID, "Теперь, введи имя города или поселка:"),
		HandleFunc: func(data string) {
			g.adminActions.CreateTownAction.FieldRegionID = regionID
			g.adminActions.CreateTownAction.isActive = true
		},
	})

	return AddRegionNode
}

func (g *TownpointGraph) ChangeTownNode(regionID uuid.UUID) *components.MarkupNode {
	ctx := context.Background()
	towns := g.entity.Town.Query().Where(town.RegionID(regionID)).AllX(ctx)

	townAnswers := make([][]*components.MarkupAnswer, 0)
	townNode := components.NewMarkupNode(components.MarkupParams{
		Text: "Выберите город или поселок для изменения",
	})

	for _, tw := range towns {
		markupAnswer := &components.MarkupAnswer{
			Content: tgbotapi.NewInlineKeyboardButtonData(tw.Name, tw.ID.String()),
			Next: components.NewMarkupNode(components.MarkupParams{
				Message: tgbotapi.NewMessage(g.chatID, "Введите новое название города или поселка"),
				HandleFunc: func(data string) {
					g.adminActions.UpdateTownAction.FieldTownID = tw.ID
					g.adminActions.UpdateTownAction.isActive = true
				},
			}),
		}

		townAnswers = append(townAnswers, []*components.MarkupAnswer{
			markupAnswer,
		})
	}

	townNode.SetAnswers(townAnswers)

	return townNode
}

func (g *TownpointGraph) DeleteTownNode(regionID uuid.UUID) *components.MarkupNode {
	ctx := context.Background()
	towns := g.entity.Town.Query().Where(town.RegionID(regionID)).AllX(ctx)

	townAnswers := make([][]*components.MarkupAnswer, 0)
	townNode := components.NewMarkupNode(components.MarkupParams{
		Text: "Выберите город или поселок для удаления",
	})

	for _, tw := range towns {
		markupAnswer := &components.MarkupAnswer{
			Content: tgbotapi.NewInlineKeyboardButtonData(tw.Name, tw.ID.String()),
			Next: components.NewMarkupNode(components.MarkupParams{
				Message: tgbotapi.NewMessage(g.chatID, "Город или поселок удален"),
				HandleFunc: func(data string) {
					_, err := g.entity.Town.Delete().Where(town.ID(uuid.FromStringOrNil(data))).Exec(ctx)
					if err != nil {
						g.api.Send(tgbotapi.NewMessage(g.chatID, "Ошибка удаления города - нельзя удалить город, если в нем точки"))
					}
				},
			}),
		}

		townAnswers = append(townAnswers, []*components.MarkupAnswer{
			markupAnswer,
		})
	}

	townNode.SetAnswers(townAnswers)

	return townNode
}

func (g *TownpointGraph) GetTownNode(regionNode *components.MarkupNode, regionID uuid.UUID) *components.MarkupNode {
	townNode := components.NewMarkupNode(components.MarkupParams{
		Text: `Далее необходимо выбрать город или поселок`,
	})

	ctx := context.Background()
	towns := g.entity.Town.Query().Where(town.RegionID(regionID)).AllX(ctx)

	townAnswers := make([][]*components.MarkupAnswer, 0)

	if g.isAdmin {
		townAnswers = append(townAnswers, []*components.MarkupAnswer{
			&components.MarkupAnswer{
				Content: tgbotapi.NewInlineKeyboardButtonData("Добавить", "Добавить"),
				Next:    g.AddTownNode(regionID),
			},
			&components.MarkupAnswer{
				Content: tgbotapi.NewInlineKeyboardButtonData("Удалить", "Удалить"),
				Next:    g.DeleteTownNode(regionID),
			},
			&components.MarkupAnswer{
				Content: tgbotapi.NewInlineKeyboardButtonData("Изменить", "Изменить"),
				Next:    g.ChangeTownNode(regionID),
			},
		})
	}

	for _, town := range towns {
		pointNode := g.GetPointNode(townNode, town.ID)
		markupAnswer := &components.MarkupAnswer{
			Content: tgbotapi.NewInlineKeyboardButtonData(town.Name, town.ID.String()),
		}
		markupAnswer.Next = pointNode

		townAnswers = append(townAnswers, []*components.MarkupAnswer{
			markupAnswer,
		})
	}
	townAnswers = append(townAnswers, []*components.MarkupAnswer{
		&components.MarkupAnswer{
			Content: tgbotapi.NewInlineKeyboardButtonData("Назад", "Назад"),
			Next:    regionNode,
		}})

	townNode.SetAnswers(townAnswers)

	return townNode
}
