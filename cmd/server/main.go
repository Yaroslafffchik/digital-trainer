package main

import (
	"digital-trainer/internal/config"
	"digital-trainer/internal/db"
	"digital-trainer/internal/handler"
	"digital-trainer/internal/server/tcp"
	"digital-trainer/internal/server/udp"
	"digital-trainer/internal/server/websocket"
	"log"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки конфигурации:", err)
	}

	dbConn, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Ошибка подключения к базе данных:", err)
	}
	defer dbConn.Close()

	// Запуск TCP сервера
	go tcp.StartTCPServer(cfg.TCPPort, dbConn)

	// Запуск UDP сервера
	go udp.StartUDPServer(cfg.UDPPort, dbConn)

	// Запуск WebSocket сервера
	go websocket.StartWebSocketServer(cfg.WebSocketPort, dbConn)

	// Запуск HTTP сервера
	handler.StartHTTPServer(cfg.HTTPPort, dbConn)
}
