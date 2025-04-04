package utils

// FormatHelpMessage generates the help message text with appropriate commands based on user permissions
func FormatHelpMessage(isAdmin bool) string {
	helpText := "<b>📋 Доступные команды</b>\n\n" +
		"<b>🏠 Основные</b>\n" +
		"└ /start - Приветственное сообщение\n" +
		"└ /help - Показать список моих команд\n\n" +
		"<b>🔍 Поиск</b>\n" +
		"└ /tools - Найти инструменты из канала «Инструменты»\n" +
		"└ /content - Найти видео из канала «Видео-контент»\n\n" +
		"<b>📅 Мероприятия</b>\n" +
		"└ /events - Показать список предстоящих мероприятий\n" +
		"└ /topics - Просмотреть темы и вопросы к предстоящим мероприятиям\n" +
		"└ /topicAdd - Предложить тему или вопрос к предстоящему мероприятию\n\n" +
		"<i>💡 <a href=\"https://t.me/c/2069889012/127/9470\">Открыть полное руководство</a></i>"

	if isAdmin {
		adminHelpText := "\n\n<b>🔐 Команды администратора</b>\n" +
			"└ /eventEdit - Редактировать мероприятие\n" +
			"└ /eventSetup - Создать новое мероприятие\n" +
			"└ /eventDelete - Удалить мероприятие\n" +
			"└ /eventFinish - Отметить мероприятие как завершенное\n" +
			"└ /showTopics - Просмотреть темы и вопросы к предстоящим мероприятиям с возможностью удаления\n" +
			"└ /code - Ввести код для авторизации TG-клиента (задом наперед)\n" +
			"└ /trySummarize - Тестирование саммаризации общения в клубе\n"

		helpText += adminHelpText
	}

	return helpText
}
