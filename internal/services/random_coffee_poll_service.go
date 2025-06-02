package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"evo-bot-go/internal/config"
	"evo-bot-go/internal/database/repositories"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

// RandomCoffeePollService handles the business logic for random coffee polls
type RandomCoffeePollService struct {
	config     *config.Config
	pollSender *PollSenderService
	pollRepo   *repositories.RandomCoffeePollRepository
}

// NewRandomCoffeePollService creates a new random coffee poll service
func NewRandomCoffeePollService(config *config.Config, pollSender *PollSenderService, pollRepo *repositories.RandomCoffeePollRepository) *RandomCoffeePollService {
	return &RandomCoffeePollService{
		config:     config,
		pollSender: pollSender,
		pollRepo:   pollRepo,
	}
}

// SendRandomCoffeePoll sends the random coffee poll
func (s *RandomCoffeePollService) SendRandomCoffeePoll(ctx context.Context) error {
	chatID := s.config.SuperGroupChatID
	if chatID == 0 {
		log.Println("Random Coffee Poll Service: SuperGroupChatID is not configured. Skipping poll.")
		return nil
	}

	if s.config.RandomCoffeeTopicID == 0 {
		return fmt.Errorf("Random Coffee Poll Service: RandomCoffeeTopicID is not configured")
	}

	// Send the poll
	question := "📝 Готов ли ты участвовать в рандомных кофе-встречах на следующей неделе?\n\nКак это работает: в конце каждой недели я буду спрашивать здесь, хочешь ли ты участвовать во встречах. Если ответишь «да», то в понедельник тебя могут объединить в пару с другим участником для неформального общения!"
	answers := []gotgbot.InputPollOption{
		{Text: "Да, участвую! ☕️"},
		{Text: "Нет, пропускаю эту неделю"},
	}
	options := &gotgbot.SendPollOpts{
		IsAnonymous:           false,
		AllowsMultipleAnswers: false,
		MessageThreadId:       int64(s.config.RandomCoffeeTopicID),
	}
	sentPollMsg, err := s.pollSender.SendPoll(chatID, question, answers, options)
	if err != nil {
		return err
	}

	// Save to database
	return s.savePollToDB(sentPollMsg)
}

// savePollToDB saves the poll information to the database
func (s *RandomCoffeePollService) savePollToDB(sentPollMsg *gotgbot.Message) error {
	if s.pollRepo == nil {
		log.Println("Random Coffee Poll Service: pollRepo is nil, skipping DB interaction.")
		return nil
	}

	// Calculate next Monday (week start date)
	now := time.Now().UTC()
	daysUntilMonday := (8 - int(now.Weekday())) % 7
	if daysUntilMonday == 0 {
		daysUntilMonday = 7 // Next Monday if today is Monday
	}

	weekStartDate := now.AddDate(0, 0, daysUntilMonday)
	weekStartDate = time.Date(weekStartDate.Year(), weekStartDate.Month(), weekStartDate.Day(), 0, 0, 0, 0, time.UTC)

	log.Printf("Random Coffee Poll Service: Calculated WeekStartDate: %s (UTC)", weekStartDate.Format("2006-01-02"))

	newPollEntry := repositories.RandomCoffeePoll{
		MessageID:      sentPollMsg.MessageId,
		TelegramPollID: sentPollMsg.Poll.Id,
		WeekStartDate:  weekStartDate,
	}

	pollID, err := s.pollRepo.CreatePoll(newPollEntry)
	if err != nil {
		log.Printf("Random Coffee Poll Service: Failed to save random coffee poll to DB: %v. Poll Message ID: %d", err, sentPollMsg.MessageId)
		return err
	}

	log.Printf("Random Coffee Poll Service: Random coffee poll saved to DB with ID: %d, Original MessageID: %d, WeekStartDate: %s",
		pollID, sentPollMsg.MessageId, weekStartDate.Format("2006-01-02"))

	return nil
}
