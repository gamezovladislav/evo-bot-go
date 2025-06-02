package tasks

import (
	"context"
	"log"
	"time"

	"evo-bot-go/internal/config"
	"evo-bot-go/internal/database/repositories"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

// RandomCoffeePollTask handles scheduling of random coffee polls
type RandomCoffeePollTask struct {
	config   *config.Config
	bot      *gotgbot.Bot
	pollRepo *repositories.RandomCoffeePollRepository
	stop     chan struct{}
}

// NewRandomCoffeePollTask creates a new random coffee poll task
func NewRandomCoffeePollTask(config *config.Config, bot *gotgbot.Bot, pollRepo *repositories.RandomCoffeePollRepository) *RandomCoffeePollTask {
	return &RandomCoffeePollTask{
		config:   config,
		bot:      bot,
		pollRepo: pollRepo,
		stop:     make(chan struct{}),
	}
}

// Start starts the random coffee poll task
func (t *RandomCoffeePollTask) Start() {
	if !t.config.RandomCoffeePollTaskEnabled {
		log.Println("Random Coffee Poll Task: Random coffee poll task is disabled")
		return
	}
	log.Printf("Random Coffee Poll Task: Starting random coffee poll task with time %02d:%02d UTC on %s",
		t.config.RandomCoffeePollTime.Hour(),
		t.config.RandomCoffeePollTime.Minute(),
		t.config.RandomCoffeePollDay.String())
	go t.run()
}

// Stop stops the random coffee poll task
func (t *RandomCoffeePollTask) Stop() {
	log.Println("Random Coffee Poll Task: Stopping random coffee poll task")
	close(t.stop)
}

// run runs the random coffee poll task
func (t *RandomCoffeePollTask) run() {
	nextRun := t.calculateNextRun()
	log.Printf("Random Coffee Poll Task: Next random coffee poll scheduled for: %v", nextRun)

	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-t.stop:
			return
		case now := <-ticker.C:
			if now.After(nextRun) {
				log.Println("Random Coffee Poll Task: Running scheduled random coffee poll")

				go func() {
					ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
					defer cancel()

					if err := t.sendRandomCoffeePoll(ctx); err != nil {
						log.Printf("Random Coffee Poll Task: Error sending random coffee poll: %v", err)
					}
				}()

				nextRun = t.calculateNextRun()
				log.Printf("Random Coffee Poll Task: Next random coffee poll scheduled for: %v", nextRun)
			}
		}
	}
}

// calculateNextRun calculates the next run time
func (t *RandomCoffeePollTask) calculateNextRun() time.Time {
	now := time.Now().UTC()
	targetHour := t.config.RandomCoffeePollTime.Hour()
	targetMinute := t.config.RandomCoffeePollTime.Minute()
	targetWeekday := t.config.RandomCoffeePollDay

	// Calculate days until target weekday
	daysUntilTarget := (int(targetWeekday) - int(now.Weekday()) + 7) % 7

	// Create target time for today
	targetTime := time.Date(now.Year(), now.Month(), now.Day(), targetHour, targetMinute, 0, 0, time.UTC)

	if daysUntilTarget == 0 && now.Before(targetTime) {
		// Today is target day and time hasn't passed yet
		return targetTime
	}

	// Either not target day or time has passed - schedule for next occurrence
	if daysUntilTarget == 0 {
		daysUntilTarget = 7 // Next week
	}

	return targetTime.AddDate(0, 0, daysUntilTarget)
}

// sendRandomCoffeePoll sends the random coffee poll
func (t *RandomCoffeePollTask) sendRandomCoffeePoll(ctx context.Context) error {
	chatID := t.config.SuperGroupChatID
	if chatID == 0 {
		log.Println("Random Coffee Poll Task: SuperGroupChatID is not configured. Skipping poll.")
		return nil
	}

	// Send the poll
	sentPollMsg, err := t.sendPoll(chatID)
	if err != nil {
		return err
	}

	// Save to database
	return t.savePollToDB(sentPollMsg)
}

// sendPoll sends the actual poll message
func (t *RandomCoffeePollTask) sendPoll(chatID int64) (*gotgbot.Message, error) {
	question := "📝 Готов ли ты участвовать в рандомных кофе-встречах на следующей неделе?\n\nКак это работает: в конце каждой недели я буду спрашивать здесь, хочешь ли ты участвовать во встречах. Если ответишь «да», то в понедельник тебя могут объединить в пару с другим участником для неформального общения!"
	options := []gotgbot.InputPollOption{
		{Text: "Да, участвую! ☕️"},
		{Text: "Нет, пропускаю эту неделю"},
	}
	opts := &gotgbot.SendPollOpts{
		IsAnonymous:           false,
		AllowsMultipleAnswers: false,
	}

	log.Printf("Random Coffee Poll Task: Sending poll to chat ID %d", chatID)
	sentPollMsg, err := t.bot.SendPoll(chatID, question, options, opts)
	if err != nil {
		return nil, err
	}

	log.Printf("Random Coffee Poll Task: Poll sent successfully. MessageID: %d, ChatID: %d", sentPollMsg.MessageId, sentPollMsg.Chat.Id)
	return sentPollMsg, nil
}

// savePollToDB saves the poll information to the database
func (t *RandomCoffeePollTask) savePollToDB(sentPollMsg *gotgbot.Message) error {
	if t.pollRepo == nil {
		log.Println("Random Coffee Poll Task: pollRepo is nil, skipping DB interaction.")
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

	log.Printf("Random Coffee Poll Task: Calculated WeekStartDate: %s (UTC)", weekStartDate.Format("2006-01-02"))

	newPollEntry := repositories.RandomCoffeePoll{
		MessageID:      sentPollMsg.MessageId,
		ChatID:         sentPollMsg.Chat.Id,
		TelegramPollID: sentPollMsg.Poll.Id,
		WeekStartDate:  weekStartDate,
	}

	pollID, err := t.pollRepo.CreatePoll(newPollEntry)
	if err != nil {
		log.Printf("Random Coffee Poll Task: Failed to save random coffee poll to DB: %v. Poll Message ID: %d", err, sentPollMsg.MessageId)
		return err
	}

	log.Printf("Random Coffee Poll Task: Random coffee poll saved to DB with ID: %d, Original MessageID: %d, WeekStartDate: %s",
		pollID, sentPollMsg.MessageId, weekStartDate.Format("2006-01-02"))

	return nil
}
