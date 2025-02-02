package tgbot

import (
	"context"
	"log"
	"path/filepath"
	"time"
	"townpoint_bot/config"
	ent "townpoint_bot/ent/generated"
	"townpoint_bot/ent/generated/point"
	"townpoint_bot/ent/generated/region"
	"townpoint_bot/ent/generated/town"
	"townpoint_bot/ent/generated/user"
	"townpoint_bot/ent/schema/common"
	"townpoint_bot/internal/services"
	"townpoint_bot/pkg/components"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gofrs/uuid"
)

// TOWN ACTIONS

type CreateTownAction struct {
	isActive      bool
	FieldName     string
	FieldRegionID uuid.UUID
}

func (a CreateTownAction) Clear() {
	a = CreateTownAction{}
}

type UpdateTownAction struct {
	isActive      bool
	FieldName     string
	FieldRegionID uuid.UUID
	FieldTownID   uuid.UUID
}

func (a UpdateTownAction) Clear() {
	a = UpdateTownAction{}
}

// REGION ACTIONS

type CreateRegionAction struct {
	isActive  bool
	FieldName string
}

func (a CreateRegionAction) Clear() {
	a = CreateRegionAction{}
}

type UpdateRegionAction struct {
	isActive      bool
	FieldName     string
	FieldRegionID uuid.UUID
}

func (a UpdateRegionAction) Clear() {
	a = UpdateRegionAction{}
}

// POINT ACTIONS

type CreatePointAction struct {
	isActive         bool
	FieldTownID      uuid.UUID
	FieldName        string
	FieldAddress     string
	FieldDescription string
	FieldPhone       string
	IsDownloadVideo  bool
}

func (a CreatePointAction) Clear() {
	a = CreatePointAction{}
}

type UpdatePointAction struct {
	isActive               bool
	FieldPointID           uuid.UUID
	FieldChangeName        bool
	FieldChangeAddress     bool
	FieldChangeDescription bool
	FieldChangePhone       bool
	FieldChangeVideo       bool
}

func (a UpdatePointAction) Clear() {
	a = UpdatePointAction{}
}

// ALL ADMINS ACTIONS

type AdminActions struct {
	CreateRegionAction
	UpdateRegionAction

	CreateTownAction
	UpdateTownAction

	CreatePointAction
	UpdatePointAction
}

type TownpointGraph struct {
	config *config.Config
	chatID int64
	api    *tgbotapi.BotAPI
	entity *ent.Client
	yadisk *services.Yadisk

	currentNode        *components.MarkupNode
	currentMessagedIDs []int
	isAdmin            bool
	adminActions       AdminActions
}

func NewTownpointGraph(config *config.Config, api *tgbotapi.BotAPI, entity *ent.Client, yadisk *services.Yadisk, chatID int64, userTGName string) *TownpointGraph {
	ctx := context.Background()

	isAdmin := false
	us, err := entity.User.Query().Where(user.TgID(userTGName)).Only(ctx)
	if err == nil && us != nil {
		if us.Role == user.RoleAdmin {
			isAdmin = true
		}
	}

	return &TownpointGraph{
		config:  config,
		chatID:  chatID,
		api:     api,
		entity:  entity,
		isAdmin: isAdmin,
		yadisk:  yadisk,
	}
}

func (g *TownpointGraph) Start() {
	g.currentNode = g.NewTownpointGraph(g.chatID)
	msg := g.currentNode.BuildMessage(g.chatID)
	_, err := g.api.Send(msg)
	if err != nil {
		log.Printf("[Graph] Error graph started %d: %s", g.chatID, err)
	}
}

func (g *TownpointGraph) Quit() {
	return
}

func (g *TownpointGraph) downloadVideo(mess *tgbotapi.Message, pointID uuid.UUID) {
	ctx := context.Background()
	g.api.Send(tgbotapi.NewMessage(g.chatID, "Началась загрузка видео, вам отпишет бот по готовности"))
	f, err := g.api.GetFile(tgbotapi.FileConfig{
		FileID: mess.Video.FileID,
	})
	if err != nil {
		g.api.Send(tgbotapi.NewMessage(g.chatID, "Что-то не так с загрузкой с телеграмма, обратитесь к разработчику: "+err.Error()))
		return
	}

	urlLink := f.Link(g.api.Token)
	path := filepath.Join(g.config.Yadisk.Path, mess.Video.FileName)
	upR, err := g.yadisk.Client.API.UploadByURL(path, urlLink)
	if err != nil {
		g.api.Send(tgbotapi.NewMessage(g.chatID, "Что-то не так с загрузкой на yadisk, обратитесь к разработчику: "+err.Error()))
		return
	}

	for {
		time.Sleep(3 * time.Second)
		s, err := g.yadisk.Client.API.GetOperationStatus(upR.Href)
		if err != nil {
			g.api.Send(tgbotapi.NewMessage(g.chatID, "Что-то не так с получением статуса загрузки на yadisk, обратитесь к разработчику: "+err.Error()))
			return
		}
		if s.IsFailed() {
			g.api.Send(tgbotapi.NewMessage(g.chatID, "Что-то не так с загрузкой на yadisk, обратитесь к разработчику: Операция Failed"))
			return
		}
		if s.IsSuccess() {
			break
		}
	}

	op, err := g.yadisk.Client.API.Publish(path)
	if err != nil {
		g.api.Send(tgbotapi.NewMessage(g.chatID, "Что-то не так с публикацией на yadisk, обратитесь к разработчику: "+err.Error()))
		return
	}

	info, err := g.yadisk.Client.API.PublishInfo(op.Href)
	if err != nil {
		g.api.Send(tgbotapi.NewMessage(g.chatID, "Что-то не так с инфо о публикации на yadisk, обратитесь к разработчику: "+err.Error()))
		return
	}
	err = g.entity.Point.Update().SetVideos([]common.PointVideo{
		common.PointVideo{
			Path: info.PublicUrl,
		},
	}).Where(point.ID(pointID)).Exec(ctx)
	if err != nil {
		g.api.Send(tgbotapi.NewMessage(g.chatID, "Что-то не так с сохранением видео в БД, обратитесь к разработчику: "+err.Error()))
		return
	}
	g.api.Send(tgbotapi.NewMessage(g.chatID, "Видео для точки: "+g.adminActions.CreatePointAction.FieldName+" опубликовано"))
}

func (g *TownpointGraph) ContinueM(mess *tgbotapi.Message) {
	if !g.isAdmin {
		return
	}

	ctx := context.Background()

	// REGIONS

	if g.adminActions.CreateRegionAction.isActive {
		if g.adminActions.CreateRegionAction.FieldName == "" {
			g.adminActions.CreateRegionAction.FieldName = mess.Text
		}

		err := g.entity.Region.Create().SetName(g.adminActions.CreateRegionAction.FieldName).SetStatus(region.StatusCreated).Exec(ctx)
		if err != nil {
			g.api.Send(tgbotapi.NewMessage(g.chatID, "Ошибка создания региона: "+err.Error()))
		} else {
			g.api.Send(tgbotapi.NewMessage(g.chatID, "Регион успешно создан"))
		}
		g.adminActions.CreateRegionAction.Clear()
		return
	}

	if g.adminActions.UpdateRegionAction.isActive {
		if g.adminActions.UpdateRegionAction.FieldName == "" {
			g.adminActions.UpdateRegionAction.FieldName = mess.Text
		}

		err := g.entity.Region.Update().SetName(g.adminActions.UpdateRegionAction.FieldName).Where(region.ID(g.adminActions.UpdateRegionAction.FieldRegionID)).Exec(ctx)
		if err != nil {
			g.api.Send(tgbotapi.NewMessage(g.chatID, "Ошибка изменения региона: "+err.Error()))
		} else {
			g.api.Send(tgbotapi.NewMessage(g.chatID, "Регион успешно изменен"))
		}
		g.adminActions.UpdateRegionAction.Clear()
		return
	}

	// TOWNS

	if g.adminActions.CreateTownAction.isActive {
		if g.adminActions.CreateTownAction.FieldName == "" {
			g.adminActions.CreateTownAction.FieldName = mess.Text
		}

		err := g.entity.Town.Create().SetName(g.adminActions.CreateTownAction.FieldName).SetRegionID(g.adminActions.CreateTownAction.FieldRegionID).SetStatus(town.StatusCreated).Exec(ctx)
		if err != nil {
			g.api.Send(tgbotapi.NewMessage(g.chatID, "Ошибка создания города: "+err.Error()))
		} else {
			g.api.Send(tgbotapi.NewMessage(g.chatID, "Город успешно создан"))
		}
		g.adminActions.CreateTownAction.Clear()
		return
	}

	if g.adminActions.UpdateTownAction.isActive {
		if g.adminActions.UpdateTownAction.FieldName == "" {
			g.adminActions.UpdateTownAction.FieldName = mess.Text
		}

		err := g.entity.Town.Update().SetName(g.adminActions.UpdateTownAction.FieldName).Where(town.ID(g.adminActions.UpdateTownAction.FieldTownID)).Exec(ctx)
		if err != nil {
			g.api.Send(tgbotapi.NewMessage(g.chatID, "Ошибка изменения города: "+err.Error()))
		} else {
			g.api.Send(tgbotapi.NewMessage(g.chatID, "Город успешно изменен"))
		}
		g.adminActions.UpdateTownAction.Clear()
		return
	}

	// POINTS

	if g.adminActions.CreatePointAction.isActive {
		if g.adminActions.CreatePointAction.FieldName == "" {
			g.adminActions.CreatePointAction.FieldName = mess.Text

			g.api.Send(tgbotapi.NewMessage(g.chatID, "Теперь введите Адрес точки: "))
			return
		}

		if g.adminActions.CreatePointAction.FieldAddress == "" {
			g.adminActions.CreatePointAction.FieldAddress = mess.Text
			g.api.Send(tgbotapi.NewMessage(g.chatID, "Теперь введите Описание точки: "))
			return
		}

		if g.adminActions.CreatePointAction.FieldDescription == "" {
			g.adminActions.CreatePointAction.FieldDescription = mess.Text

			g.api.Send(tgbotapi.NewMessage(g.chatID, "Теперь введите Телефон, контакты точки: "))
			return
		}

		if g.adminActions.CreatePointAction.FieldPhone == "" {
			g.adminActions.CreatePointAction.FieldPhone = mess.Text

			g.api.Send(tgbotapi.NewMessage(g.chatID, "Теперь загрузите видео или напишите любой текст если не хотите загружать"))
			return
		}

		pointID, _ := uuid.NewV4()
		err := g.entity.Point.Create().
			SetID(pointID).
			SetName(g.adminActions.CreatePointAction.FieldName).
			SetAddress(g.adminActions.CreatePointAction.FieldAddress).
			SetDescription(g.adminActions.CreatePointAction.FieldDescription).
			SetPhone(g.adminActions.FieldPhone).
			SetTownID(g.adminActions.CreatePointAction.FieldTownID).
			SetStatus(point.StatusCreated).
			Exec(ctx)
		if err != nil {
			g.api.Send(tgbotapi.NewMessage(g.chatID, "Ошибка создания точки: "+err.Error()))
		} else {
			g.api.Send(tgbotapi.NewMessage(g.chatID, "Точка успешно создана"))
		}

		if !g.adminActions.CreatePointAction.IsDownloadVideo {
			if mess.Video != nil {
				go g.downloadVideo(mess, pointID)
			}
		}

		g.adminActions.CreatePointAction.Clear()
		return
	}

	if g.adminActions.UpdatePointAction.isActive {
		if g.adminActions.UpdatePointAction.FieldChangeName {
			err := g.entity.Point.Update().SetName(mess.Text).Where(point.ID(g.adminActions.UpdatePointAction.FieldPointID)).Exec(ctx)
			if err != nil {
				g.api.Send(tgbotapi.NewMessage(g.chatID, "Ошибка изменения точки: "+err.Error()))
			}
		}

		if g.adminActions.UpdatePointAction.FieldChangeDescription {
			err := g.entity.Point.Update().SetDescription(mess.Text).Where(point.ID(g.adminActions.UpdatePointAction.FieldPointID)).Exec(ctx)
			if err != nil {
				g.api.Send(tgbotapi.NewMessage(g.chatID, "Ошибка изменения точки: "+err.Error()))
			}
		}

		if g.adminActions.UpdatePointAction.FieldChangeAddress {
			err := g.entity.Point.Update().SetAddress(mess.Text).Where(point.ID(g.adminActions.UpdatePointAction.FieldPointID)).Exec(ctx)
			if err != nil {
				g.api.Send(tgbotapi.NewMessage(g.chatID, "Ошибка изменения точки: "+err.Error()))
			}
		}

		if g.adminActions.UpdatePointAction.FieldChangePhone {
			err := g.entity.Point.Update().SetPhone(mess.Text).Where(point.ID(g.adminActions.UpdatePointAction.FieldPointID)).Exec(ctx)
			if err != nil {
				g.api.Send(tgbotapi.NewMessage(g.chatID, "Ошибка изменения точки: "+err.Error()))
			}
		}

		if g.adminActions.UpdatePointAction.FieldChangeVideo {
			go g.downloadVideo(mess, g.adminActions.UpdatePointAction.FieldPointID)
		}

		g.api.Send(tgbotapi.NewMessage(g.chatID, "точка изменена"))
		g.adminActions.UpdatePointAction.Clear()
	}
}

func (g *TownpointGraph) Continue(query *tgbotapi.CallbackQuery) {
	for _, messageID := range g.currentMessagedIDs {
		if messageID == query.Message.MessageID {
			return
		}
	}
	g.currentMessagedIDs = append(g.currentMessagedIDs, query.Message.MessageID)
	g.currentNode = g.currentNode.Continue(query, g.api)
}

func (g *TownpointGraph) NewTownpointGraph(chatID int64) *components.MarkupNode {
	// laterDuration := time.Hour * 24 * 3
	// afkDuration := time.Hour

	// ==================================================================================================================================================
	// region node

	regionNode := g.GetRegionsNode()
	return regionNode
}
