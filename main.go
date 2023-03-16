package main

import (
	"fmt"
	"github.com/Syfaro/telegram-bot-api"
	"krip_bot/repository"
	_ "krip_bot/repository"
	"krip_bot/search"
	_ "krip_bot/search"
	encoder "krip_bot/url-encode"
	"os"
	"reflect"
	"time"
)

func telegramBot() {
	//docker build -t krip_bot -f dockerfile .
	//Создаем бота
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TOKEN"))
	if err != nil {
		panic(err)
	}
	fmt.Sprintf("бот запущен")

	//Устанавливаем время обновления
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	fmt.Sprintf("время обновления: %s", u.Timeout)

	//Получаем обновления от бота
	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		//Проверяем что от пользователья пришло именно текстовое сообщение
		if reflect.TypeOf(update.Message.Text).Kind() == reflect.String && update.Message.Text != "" {

			switch update.Message.Text {
			case "/start":

				//Отправлем сообщение
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Привет,я могу найти что угодно на Wikipedia. Пожалуйста используйте Английский язык для запросов!")
				bot.Send(msg)

			case "/number_of_users":

				if os.Getenv("DB_SWITCH") == "on" {

					//Присваиваем количество пользоватьелей использовавших бота в num переменную
					num, err := repository.GetNumberOfUsers(update.Message.Chat.UserName)
					if err != nil {

						//Отправлем сообщение
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка баз данных")
						bot.Send(msg)
					}

					//Создаем строку которая содержит колличество пользователей использовавших бота
					ans := fmt.Sprintf("Cтолько пользователей воспользовались моими функциями - %d ", num)

					//Отправлем сообщение
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, ans)
					bot.Send(msg)
				} else {

					//Отправлем сообщение
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "База данных не подключена, невозможно ответить")
					bot.Send(msg)
				}

			case "/requests":

				if os.Getenv("DB_SWITCH") == "on" {

					req, count, err := repository.GetListOfRequests(update.Message.Chat.UserName)
					if err != nil {

						//Отправлем сообщение
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка баз данных")
						bot.Send(msg)
					}

					ans := fmt.Sprintf("Ваша история поиска: %s", req)

					//Отправлем сообщение
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, ans)
					bot.Send(msg)

					ans2 := fmt.Sprintf(" Колличество запросов: %d", count)

					//Отправлем сообщение
					msg2 := tgbotapi.NewMessage(update.Message.Chat.ID, ans2)
					bot.Send(msg2)
				} else {

					//Отправлем сообщение
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "База данных не подключена, невозможно ответить")
					bot.Send(msg)
				}
			default:

				//Устанавливаем язык для поиска в википедии
				language := os.Getenv("LANGUAGE")

				//Создаем url для поиска
				ms, _ := encoder.UrlEncoded(update.Message.Text)

				url := ms
				request := "https://" + language + ".wikipedia.org/w/api.php?action=opensearch&search=" + url + "&limit=3&origin=*&format=json"

				//Присваем данные среза с ответом в переменную message
				message := search.WikipediaAPI(request)

				if os.Getenv("DB_SWITCH") == "on" {

					//Отправляем username, chat_id, message, answer в БД
					if err := repository.CollectData(update.Message.Chat.UserName, update.Message.Chat.ID, update.Message.Text, message); err != nil {

						//Отправлем сообщение
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка базы, но бот работает.")
						bot.Send(msg)
					}
				}

				//Проходим через срез и отправляем каждый элемент пользователю
				for _, val := range message {

					//Отправлем сообщение
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, val)
					bot.Send(msg)
				}
			}
		} else {

			//Отправлем сообщение
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Используйте слова для поиска")
			bot.Send(msg)
		}
	}
}

func main() {

	time.Sleep(1 * time.Minute)

	//Создаем таблицу
	if os.Getenv("CREATE_TABLE") == "yes" {

		if os.Getenv("DB_SWITCH") == "on" {

			if err := repository.CreateTable(); err != nil {
				fmt.Sprintf("не возможно создать базу, возможно она уже существует")
			}
		}
	}

	time.Sleep(1 * time.Minute)

	//Вызываем бота
	telegramBot()
}
