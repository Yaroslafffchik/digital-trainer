package websocket

import (
	"database/sql"
	"digital-trainer/internal/models"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Разрешить все источники (для упрощения разработки)
	},
}

func StartWebSocketServer(port string, db *sql.DB) {
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handleWebSocket(w, r, db)
	})

	// Генерация челленджей каждые 1-5 минут
	go generateChallenges(db)

	log.Printf("WebSocket сервер запущен на %s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func handleWebSocket(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Ошибка установки WebSocket соединения:", err)
		return
	}
	defer conn.Close()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Ошибка чтения WebSocket сообщения:", err)
			break
		}

		var challenge models.Challenge
		err = json.Unmarshal(msg, &challenge)
		if err != nil {
			log.Println("Ошибка декодирования челленджа:", err)
			continue
		}

		// Сохранение челленджа в базе данных
		_, err = db.Exec(
			"INSERT INTO challenges (name, description, created_at) VALUES ($1, $2, CURRENT_TIMESTAMP)",
			challenge.Name, challenge.Description,
		)
		if err != nil {
			log.Println("Ошибка сохранения челленджа:", err)
			continue
		}

		// Отправка ответа клиенту
		response := map[string]string{"status": "Челлендж получен"}
		conn.WriteJSON(response)
	}
}

func generateChallenges(db *sql.DB) {
	challenges := []models.Challenge{
		{Name: "Спринт на 5 км", Description: "Пробегите 5 км за минимальное время!"},
		{Name: "100 отжиманий", Description: "Выполните 100 отжиманий за один подход!"},
		{Name: "Планка 5 минут", Description: "Держите планку 5 минут без перерыва!"},
	}

	for {
		// Случайный интервал от 1 до 5 минут
		interval := time.Duration(1+(rand.Intn(4))) * time.Minute
		time.Sleep(interval)

		challenge := challenges[rand.Intn(len(challenges))]
		_, err := db.Exec(
			"INSERT INTO challenges (name, description, created_at) VALUES ($1, $2, CURRENT_TIMESTAMP)",
			challenge.Name, challenge.Description,
		)
		if err != nil {
			log.Println("Ошибка генерации челленджа:", err)
			continue
		}

		log.Printf("Сгенерирован челлендж: %s", challenge.Name)
	}
}
