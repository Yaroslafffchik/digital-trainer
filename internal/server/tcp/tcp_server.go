package tcp

import (
	"bufio"
	"database/sql"
	"digital-trainer/internal/models"
	"encoding/json"
	"log"
	"net"
)

func StartTCPServer(port string, db *sql.DB) {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal("Ошибка запуска TCP сервера:", err)
	}
	defer listener.Close()

	log.Printf("TCP сервер запущен на %s", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Ошибка принятия соединения:", err)
			continue
		}
		go handleTCPConnection(conn, db)
	}
}

func handleTCPConnection(conn net.Conn, db *sql.DB) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Println("Ошибка чтения:", err)
			return
		}

		var training models.Training
		err = json.Unmarshal([]byte(message), &training)
		if err != nil {
			conn.Write([]byte("Ошибка: неверный формат данных\n"))
			continue
		}

		_, err = db.Exec(
			"INSERT INTO trainings (name, description, session_id) VALUES ($1, $2, $3)",
			training.Name, training.Description, training.SessionID,
		)
		if err != nil {
			conn.Write([]byte("Ошибка сохранения тренировки\n"))
			continue
		}

		conn.Write([]byte("Тренировка сохранена\n"))
	}
}
