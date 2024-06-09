package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/Synkov2102/tarotbot/tarot"
	_ "github.com/lib/pq" // Драйвер для PostgreSQL
)

var db *sql.DB

func InitDB() {
	// Строка подключения к базе данных PostgreSQL
	// Замените значения на свои
	connStr := fmt.Sprintf("user=%s dbname=%s password=%s host=%s sslmode=disable", os.Getenv("POSTGRES_USR"), os.Getenv("POSTGRES_DB_NAME"), os.Getenv("POSTGRES_PASS"), os.Getenv("POSTGRES_HOST"))

	// Подключение к базе данных
	dbConn, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	// Проверка подключения к базе данных
	err = dbConn.Ping()
	if err != nil {
		log.Fatal(err)
	}

	db = dbConn
}

func GetDB() *sql.DB {
	return db
}

func CloseDB() {
	if db != nil {
		err := db.Close()
		if err != nil {
			log.Println("Ошибка при закрытии соединения с базой данных:", err)
		} else {
			log.Println("Соединение с базой данных успешно закрыто.")
		}
	}
}

// Функция для получения массива карт из базы данных
func GetTarotCards() ([]tarot.TarotCard, error) {
	rows, err := db.Query("SELECT name, suit, number, img_url FROM tarot_cards")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cards []tarot.TarotCard
	for rows.Next() {
		var card tarot.TarotCard
		err := rows.Scan(&card.Name, &card.Suit, &card.Number, &card.ImgURL)
		if err != nil {
			return nil, err
		}
		cards = append(cards, card)
	}

	// Проверка на ошибки, возникшие во время итерации
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return cards, nil
}
