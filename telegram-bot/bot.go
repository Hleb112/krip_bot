package bot

import (
	"fmt"
	"github.com/Syfaro/telegram-bot-api"
	"krip_bot/repository"
	"krip_bot/search"
	encoder "krip_bot/url-encode"
	"os"
	"reflect"
)

func TelegramBot() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TOKEN"))
	if err != nil {
		panic(err)
	}
	fmt.Sprintf("бот запущен")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	fmt.Sprintf("время обновления: %s", u.Timeout)

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if reflect.TypeOf(update.Message.Text).Kind() == reflect.String && update.Message.Text != "" {

			switch update.Message.Text {
			case "/start":

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Привет,я могу найти что угодно на Wikipedia. Пожалуйста используйте Английский язык для запросов!")
				bot.Send(msg)

			case "/number_of_users":

				if os.Getenv("DB_SWITCH") == "on" {

					num, err := repository.GetNumberOfUsers(update.Message.Chat.UserName)
					if err != nil {

						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка баз данных")
						bot.Send(msg)
					}

					ans := fmt.Sprintf("Cтолько пользователей воспользовались моими функциями - %d ", num)

					msg := tgbotapi.NewMessage(update.Message.Chat.ID, ans)
					bot.Send(msg)
				} else {

					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "База данных не подключена, невозможно ответить")
					bot.Send(msg)
				}

			case "/requests":

				if os.Getenv("DB_SWITCH") == "on" {

					req, count, err := repository.GetListOfRequests(update.Message.Chat.UserName)
					if err != nil {

						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка баз данных")
						bot.Send(msg)
					}

					ans := fmt.Sprintf("Ваша история поиска: %s", req)

					msg := tgbotapi.NewMessage(update.Message.Chat.ID, ans)
					bot.Send(msg)

					ans2 := fmt.Sprintf(" Колличество запросов: %d", count)

					msg2 := tgbotapi.NewMessage(update.Message.Chat.ID, ans2)
					bot.Send(msg2)
				} else {

					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "База данных не подключена, невозможно ответить")
					bot.Send(msg)
				}
			default:
				language := os.Getenv("LANGUAGE")
				ms, _ := encoder.UrlEncoded(update.Message.Text)
				url := ms
				request := "https://" + language + ".wikipedia.org/w/api.php?action=opensearch&search=" + url + "&limit=3&origin=*&format=json"

				message := search.WikipediaAPI(request)

				if os.Getenv("DB_SWITCH") == "on" {
					if err := repository.CollectData(update.Message.Chat.UserName, update.Message.Chat.ID, update.Message.Text, message); err != nil {

						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка базы, но бот работает.")
						bot.Send(msg)
					}
				}

				for _, val := range message {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, val)
					bot.Send(msg)
				}
			}
		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Используйте слова для поиска")
			bot.Send(msg)
		}
	}
}
