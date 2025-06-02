package formatters

import (
	"evo-bot-go/internal/constants"
	"fmt"
)

// FormatHelpMessage generates the help message text with appropriate commands based on user permissions
func FormatHelpMessage(isAdmin bool) string {
	helpText := "<b>📋 Доступные команды</b>\n\n" +
		"<b>🏠 Базовые</b>\n" +
		"└ /start - Приветственное сообщение\n" +
		"└ /help - Показать список моих команд\n" +
		"└ /cancel - Принудительно отменяет любой диалог\n\n" +
		"<b>👤 Профиль</b>\n" +
		"└ /profile - Управление своим профилем, поиск профилей клубчан, публикация и обновление информации о себе в канале «Интро»\n\n" +
		"<b>🔍 Поиск</b>\n" +
		"└ /tools - Найти инструменты из канала «Инструменты»\n" +
		"└ /content - Найти видео из канала «Видео-контент»\n" +
		"└ /intro - Найти информацию об участниках клуба из канала «Интро» (умный поиск по профилям клубчан)\n\n" +
		"<b>📅 Мероприятия</b>\n" +
		"└ /events - Показать список предстоящих мероприятий\n" +
		"└ /topics - Просмотреть темы и вопросы к предстоящим мероприятиям\n" +
		"└ /topicAdd - Предложить тему или вопрос к предстоящему мероприятию\n\n"

	featuresDescription := "\n<b>🎲 Еженедельный Random Coffee:</b>\n" +
		"- Опрос для участия в случайных встречах на следующей неделе (обычно по пятницам).\n" +
		"- Используйте опрос, чтобы указать свое участие.\n" +
		"- Пары объявляются (обычно по понедельникам) после запуска команды администратором."

	helpText += featuresDescription

	helpText += "\n\n" + // Add some spacing before the link
		"<i>💡 <a href=\"https://t.me/c/2069889012/127/9470\">Открыть полное руководство</a></i>"

	if isAdmin {
		adminHelpText := "\n\n<b>🔐 Команды администратора</b>\n" +
			fmt.Sprintf("└ /%s - Начать мероприятие\n", constants.EventStartCommand) +
			fmt.Sprintf("└ /%s - Создать новое мероприятие\n", constants.EventSetupCommand) +
			fmt.Sprintf("└ /%s - Редактировать мероприятие\n", constants.EventEditCommand) +
			fmt.Sprintf("└ /%s - Удалить мероприятие\n", constants.EventDeleteCommand) +
			fmt.Sprintf("└ /%s - Просмотреть темы и вопросы к предстоящим мероприятиям *с возможностью удаления*\n", constants.ShowTopicsCommand) +
			fmt.Sprintf("└ /%s - Ручной запуск и объявление пар для еженедельного Random Coffee\n", constants.PairRandomCoffeeCommand) +
			fmt.Sprintf("└ /%s - Запустить новый опрос по кофейным встречам\n", constants.CoffeeRestartCommand) +
			fmt.Sprintf("└ /%s - Ввести код для авторизации TG-клиента (задом наперед)\n", constants.CodeCommand) +
			fmt.Sprintf("└ /%s - Тестирование саммаризации общения в клубе\n", constants.TrySummarizeCommand) +
			fmt.Sprintf("└ /%s - Управление профилями клубчан\n", constants.AdminProfilesCommand)

		helpText += adminHelpText
	}

	return helpText
}
