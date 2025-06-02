package udp

import (
	"database/sql"
	"digital-trainer/internal/models"
	"encoding/json"
	"log"
	"net"
)

func StartUDPServer(port string, db *sql.DB) {
	addr, err := net.ResolveUDPAddr("udp", port)
	if err != nil {
		log.Fatal("Ошибка разрешения UDP адреса:", err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatal("Ошибка запуска UDP сервера:", err)
	}
	defer conn.Close()

	log.Printf("UDP сервер запущен на %s", port)

	buffer := make([]byte, 1024)
	for {
		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Println("Ошибка чтения UDP:", err)
			continue
		}

		var metrics models.Metrics
		err = json.Unmarshal(buffer[:n], &metrics)
		if err != nil {
			log.Println("Ошибка декодирования метрик:", err)
			continue
		}

		_, err = db.Exec(
			"INSERT INTO metrics (session_id, pulse, speed, timestamp) VALUES ($1, $2, $3, CURRENT_TIMESTAMP)",
			metrics.SessionID, metrics.Pulse, metrics.Speed,
		)
		if err != nil {
			log.Println("Ошибка сохранения метрик:", err)
			continue
		}

		log.Printf("Метрика сохранена от %s: %+v", addr.String(), metrics)
	}
}
