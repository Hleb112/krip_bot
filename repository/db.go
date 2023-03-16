package repository

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
	"strings"
	"time"
)

var host = os.Getenv("HOST")
var port = os.Getenv("PORT")
var user = os.Getenv("USER")
var password = os.Getenv("PASSWORD")
var dbname = os.Getenv("DBNAME")
var sslmode = os.Getenv("SSLMODE")

var dbInfo = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, dbname, sslmode)

// Собираем данные полученные ботом
func CollectData(username string, chatid int64, message string, answer []string) error {

	//Подключаемся к БД
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return err
	}
	defer db.Close()

	//Конвертируем срез с ответом в строку
	answ := strings.Join(answer, ", ")

	//Создаем SQL запрос
	data := `INSERT INTO users(username, chat_id, message, answer) VALUES($1, $2, $3, $4);`

	//Выполняем наш SQL запрос
	if _, err = db.Exec(data, `@`+username, chatid, message, answ); err != nil {
		return err
	}

	return nil
}

// Создаем таблицу users в БД при подключении к ней
func CreateTable() error {

	//Подключаемся к БД
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return err
	}
	defer db.Close()

	//Создаем таблицу users
	if _, err = db.Exec(`CREATE TABLE users(ID SERIAL PRIMARY KEY, TIMESTAMP TIMESTAMP DEFAULT CURRENT_TIMESTAMP, USERNAME TEXT, CHAT_ID INT, MESSAGE TEXT, ANSWER TEXT);`); err != nil {
		return err
	}

	return nil
}

func GetNumberOfUsers(username string) (int64, error) {
	fmt.Println(username, "- воспользовался командой number_of_users в:", time.Now())

	var count int64

	//Подключаемся к БД
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return 0, err
	}
	defer db.Close()

	//Отправляем запрос в БД для подсчета числа уникальных пользователей
	row := db.QueryRow("SELECT COUNT(DISTINCT username) FROM users;")
	err = row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func GetListOfRequests(username string) ([]string, int, error) {
	request := make([]string, 5)
	var count int
	fmt.Println(username, "- воспользовался командой requests в:", time.Now())
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return []string{"Ошибка подкдючения к базе", ":("}, 0, err
	}
	defer db.Close()

	//Отправляем запрос в БД для просмотра истории
	rows, err := db.Query(`SELECT (MESSAGE,TIMESTAMP) FROM users WHERE USERNAME = $1`, "@"+username)
	defer rows.Close()

	for rows.Next() {
		count += 1
		var item string
		err := rows.Scan(&item)
		if err != nil {
			log.Fatal(err)
		}

		request = append(request, ";", item)
	}
	return request, count, nil
}
