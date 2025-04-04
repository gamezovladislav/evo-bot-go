package utils

import (
	"log"
	"math/rand"
	"time"

	"evo-bot-go/internal/config"
	"evo-bot-go/internal/constants"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

// CheckAdminPermissions checks if the user has admin permissions and returns an appropriate error response
// Returns true if user has permission, false otherwise
func CheckAdminPermissions(b *gotgbot.Bot, ctx *ext.Context, superGroupChatID int64, commandName string) bool {
	msg := ctx.EffectiveMessage

	if !IsUserAdminOrCreator(b, msg.From.Id, superGroupChatID) {
		if _, err := msg.Reply(b, "Эта команда доступна только администраторам.", nil); err != nil {
			log.Printf("Failed to send admin-only message: %v", err)
		}
		log.Printf("User %d tried to use %s without admin rights", msg.From.Id, commandName)
		return false
	}

	return true
}

// CheckPrivateChatType checks if the command is used in a private chat and returns an appropriate error response
// Returns true if used in private chat, false otherwise
func CheckPrivateChatType(b *gotgbot.Bot, ctx *ext.Context) bool {
	msg := ctx.EffectiveMessage

	if msg.Chat.Type != constants.PrivateChatType {
		if _, err := ReplyWithCleanupAfterDelayWithPingMessage(
			b,
			msg,
			"*Прошу прощения* 🧐\n\nЭта команда работает только в _личной беседе_ со мной. "+
				"Напишите мне в ЛС, и я с удовольствием помогу (я тебя там пинганул, если мы общались ранее)."+
				"\n\nДанное сообщение и твоя команда будут автоматически удалены через 10 секунд.",
			10, &gotgbot.SendMessageOpts{
				ParseMode: "Markdown",
			}); err != nil {
			log.Printf("Failed to send private-only message: %v", err)
		}
		return false
	}

	return true
}

// ReplyWithCleanupAfterDelay replies to a message and then deletes both the reply and the original message after the specified delay
// Returns the sent message and any error that occurred during sending
func ReplyWithCleanupAfterDelayWithPingMessage(b *gotgbot.Bot, msg *gotgbot.Message, text string, delaySeconds int, opts *gotgbot.SendMessageOpts) (*gotgbot.Message, error) {
	// Reply to the message
	sentMsg, err := msg.Reply(b, text, opts)
	if err != nil {
		log.Printf("Failed to send reply: %v", err)
		return nil, err
	}

	// Send a random greeting message
	greetings := []string{
		"Ping!",
		"Hi!",
		"Ку!",
		"Приветы!",
		"Дзень добры!",
		"Пинг!",
	}
	randomGreeting := greetings[rand.Intn(len(greetings))]
	_, greetErr := b.SendMessage(msg.From.Id, randomGreeting, nil)
	if greetErr != nil {
		log.Printf("Failed to send greeting message: %v", greetErr)
	}

	// Start a goroutine to delete both messages after the delay
	go func() {
		time.Sleep(time.Duration(delaySeconds) * time.Second)

		// Delete the reply message
		_, replyErr := sentMsg.Delete(b, nil)
		if replyErr != nil {
			log.Printf("Failed to delete reply message after delay: %v", replyErr)
		}

		// Delete the original message
		_, origErr := msg.Delete(b, nil)
		if origErr != nil {
			log.Printf("Failed to delete original message after delay: %v", origErr)
		}
	}()

	return sentMsg, nil
}

func CheckClubMemberPermissions(b *gotgbot.Bot, msg *gotgbot.Message, config *config.Config, commandName string) bool {
	if !IsUserClubMember(b, msg, config) {
		if _, err := msg.Reply(b, "Эта команда доступна только участникам клуба.", nil); err != nil {
			log.Printf("Failed to send club-only message: %v", err)
		}
		log.Printf("User %d tried to use %s without club member rights", msg.From.Id, commandName)
		return false
	}

	return true
}

// CheckAdminAndPrivateChat combines permission and chat type checking for admin-only private commands
// Returns true if all checks pass, false otherwise
func CheckAdminAndPrivateChat(b *gotgbot.Bot, ctx *ext.Context, superGroupChatID int64, commandName string) bool {
	if !CheckAdminPermissions(b, ctx, superGroupChatID, commandName) {
		return false
	}

	if !CheckPrivateChatType(b, ctx) {
		return false
	}

	return true
}
