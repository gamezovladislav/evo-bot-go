package privatehandlers

import (
	"evo-bot-go/internal/config"
	"evo-bot-go/internal/constants"
	"evo-bot-go/internal/database/repositories"
	"evo-bot-go/internal/services"
	"evo-bot-go/internal/utils"
	"fmt"
	"log"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

type eventsHandler struct {
	config               *config.Config
	eventRepository      *repositories.EventRepository
	messageSenderService services.MessageSenderService
}

func NewEventsHandler(
	config *config.Config,
	eventRepository *repositories.EventRepository,
	messageSenderService services.MessageSenderService,
) ext.Handler {
	h := &eventsHandler{
		config:               config,
		eventRepository:      eventRepository,
		messageSenderService: messageSenderService,
	}

	return handlers.NewCommand(constants.EventsCommand, h.handleCommand)
}

func (h *eventsHandler) handleCommand(b *gotgbot.Bot, ctx *ext.Context) error {
	msg := ctx.EffectiveMessage

	// Only proceed if this is a private chat
	if !utils.CheckPrivateChatType(b, ctx) {
		return nil
	}

	// Check if user is a club member
	if !utils.CheckClubMemberPermissions(b, msg, h.config, constants.EventsCommand) {
		return nil
	}

	// Get actual events to show
	events, err := h.eventRepository.GetLastActualEvents(10) // Fetch last 10 actual events
	if err != nil {
		h.messageSenderService.Reply(b, msg, "Ошибка при получении списка мероприятий.", nil)
		log.Printf("EventsHandler: Error during events retrieval: %v", err)
		return nil
	}

	if len(events) == 0 {
		h.messageSenderService.Reply(b, msg, "На данный момент нет актуальных мероприятий.", nil)
		return nil
	}

	// Format and display event list
	formattedEvents := utils.FormatEventListForUsersWithoutIds(
		events,
		"📋 Список ближайших мероприятий",
	)
	formattedEvents += fmt.Sprintf("\nИспользуй команду /%s, если хочешь предложить темы и вопросы к этим мероприятиям, либо команду /%s для просмотра уже добавленных тем и вопросов.", constants.TopicAddCommand, constants.TopicsCommand)
	formattedEvents += "\n\nА вот ссылка на [клубный календарь](https://itbeard.com/s/evo-calendar), который можно добавить к себе и удобно следить всеми мероприятиями клуба."
	h.messageSenderService.ReplyMarkdown(b, msg, formattedEvents, nil)

	return nil
}
