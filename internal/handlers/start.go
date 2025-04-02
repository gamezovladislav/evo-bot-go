package handlers

import (
	"evo-bot-go/internal/config"
	"evo-bot-go/internal/constants"
	"evo-bot-go/internal/utils"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

type startHandler struct {
	config *config.Config
}

func NewStartHandler(config *config.Config) ext.Handler {
	h := &startHandler{
		config: config,
	}

	return handlers.NewCommand(constants.StartCommand, h.handleCommand)
}

func (h *startHandler) handleCommand(b *gotgbot.Bot, ctx *ext.Context) error {
	// Only proceed if this is a private chat
	if !utils.CheckPrivateChatType(b, ctx) {
		return nil
	}

	userName := ""
	if ctx.EffectiveUser.FirstName != "" {
		userName = ctx.EffectiveUser.FirstName
	}

	greeting := "Приветствую"
	if userName != "" {
		greeting += ", *" + userName + "*"
	}
	greeting += "! 🎩"

	// Check if user is a member of the club
	isClubMember := utils.IsUserClubMember(b, ctx.EffectiveMessage, h.config)

	var message string
	if isClubMember {
		// Message for club members
		message = greeting + "\n\n" +
			"Я — *Дженкинс Вебстер*, потомственный дворецкий и верный помощник клуба _\"Эволюция Кода\"_ 🧐\n\n" +
			"Рад видеть тебя среди участников нашего клуба! Я готов помочь тебе во всех твоих начинаниях. 🤵\n\n" +
			"Воспользуйся командой /help, чтобы узнать о всех моих возможностях."
	} else {
		// Message for non-members
		message = greeting + "\n\n" +
			"Я — *Дженкинс Вебстер*, потомственный дворецкий и верный помощник клуба _\"Эволюция Кода\"_ 🧐\n\n" +
			"Позвольте предложить тебе присоединиться к нашему изысканному сообществу разработчиков и разработчиц, " +
			"где я буду рад служить тебе всеми своими возможностями и ресурсами. 🤵\n\n" +
			"👉 [Жду тебя в клубе!](https://web.tribute.tg/l/get-started)"
	}

	utils.SendLoggedMarkdownReply(b, ctx.EffectiveMessage, message, nil)

	return nil
}
