package main

import (
	"fmt"
	"krip_bot/repository"
	_ "krip_bot/repository"
	_ "krip_bot/search"
	bot "krip_bot/telegram-bot"
	"os"
	"time"
)

func main() {
	if os.Getenv("CREATE_TABLE") == "yes" {
		if os.Getenv("DB_SWITCH") == "on" {
			if err := repository.CreateTable(); err != nil {
				fmt.Sprintf("не возможно создать базу, возможно она уже существует")
			}
		}
	}
	time.Sleep(30 * time.Second)
	bot.TelegramBot()
}
