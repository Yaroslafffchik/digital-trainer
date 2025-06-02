package handler

import (
	"database/sql"
	"digital-trainer/internal/models"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net"
	"net/http"
	"text/template"
)

func StartHTTPServer(port string, db *sql.DB) {
	r := mux.NewRouter()

	// Обслуживание статических файлов
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Шаблоны HTML
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t, _ := template.ParseFiles("static/templates/index.html")
		t.Execute(w, nil)
	})
	r.HandleFunc("/trainings", func(w http.ResponseWriter, r *http.Request) {
		t, _ := template.ParseFiles("static/templates/trainings.html")
		t.Execute(w, nil)
	})
	r.HandleFunc("/challenges", func(w http.ResponseWriter, r *http.Request) {
		t, _ := template.ParseFiles("static/templates/challenges.html")
		t.Execute(w, nil)
	})
	r.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		t, _ := template.ParseFiles("static/templates/metrics.html")
		t.Execute(w, nil)
	})

	// API для получения и добавления тренировок
	r.HandleFunc("/api/trainings", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			var training models.Training
			err := json.NewDecoder(r.Body).Decode(&training)
			if err != nil {
				http.Error(w, "Неверный формат данных", http.StatusBadRequest)
				return
			}

			conn, err := net.Dial("tcp", ":1234")
			if err != nil {
				http.Error(w, "Ошибка подключения к TCP серверу", http.StatusInternalServerError)
				return
			}
			defer conn.Close()

			data, err := json.Marshal(training)
			if err != nil {
				http.Error(w, "Ошибка сериализации данных", http.StatusInternalServerError)
				return
			}

			_, err = conn.Write(append(data, '\n'))
			if err != nil {
				http.Error(w, "Ошибка отправки данных", http.StatusInternalServerError)
				return
			}

			response := make([]byte, 1024)
			n, err := conn.Read(response)
			if err != nil {
				http.Error(w, "Ошибка чтения ответа", http.StatusInternalServerError)
				return
			}

			w.Write(response[:n])
		} else {
			rows, err := db.Query("SELECT id, name, description, session_id FROM trainings")
			if err != nil {
				http.Error(w, "Ошибка получения тренировок", http.StatusInternalServerError)
				return
			}
			defer rows.Close()

			var trainings []models.Training
			for rows.Next() {
				var t models.Training
				err := rows.Scan(&t.ID, &t.Name, &t.Description, &t.SessionID)
				if err != nil {
					http.Error(w, "Ошибка сканирования данных", http.StatusInternalServerError)
					return
				}
				trainings = append(trainings, t)
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(trainings)
		}
	})
	// API для получения челленджей
	r.HandleFunc("/api/challenges", func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, name, description, created_at FROM challenges")
		if err != nil {
			http.Error(w, "Ошибка получения челленджей", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var challenges []models.Challenge
		for rows.Next() {
			var c models.Challenge
			err := rows.Scan(&c.ID, &c.Name, &c.Description, &c.CreatedAt)
			if err != nil {
				http.Error(w, "Ошибка сканирования данных", http.StatusInternalServerError)
				return
			}
			challenges = append(challenges, c)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(challenges)
	})

	// API для получения и добавления метрик
	r.HandleFunc("/api/metrics", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			var metric models.Metrics
			err := json.NewDecoder(r.Body).Decode(&metric)
			if err != nil {
				http.Error(w, "Неверный формат данных", http.StatusBadRequest)
				return
			}

			conn, err := net.Dial("udp", ":1235")
			if err != nil {
				http.Error(w, "Ошибка подключения к UDP серверу", http.StatusInternalServerError)
				return
			}
			defer conn.Close()

			data, err := json.Marshal(metric)
			if err != nil {
				http.Error(w, "Ошибка сериализации данных", http.StatusInternalServerError)
				return
			}

			_, err = conn.Write(data)
			if err != nil {
				http.Error(w, "Ошибка отправки данных", http.StatusInternalServerError)
				return
			}

			w.Write([]byte("Метрика успешно добавлена"))
		} else {
			rows, err := db.Query("SELECT id, session_id, pulse, speed, timestamp FROM metrics ORDER BY timestamp DESC LIMIT 50")
			if err != nil {
				http.Error(w, "Ошибка получения метрик", http.StatusInternalServerError)
				return
			}
			defer rows.Close()

			var metrics []models.Metrics
			for rows.Next() {
				var m models.Metrics
				err := rows.Scan(&m.ID, &m.SessionID, &m.Pulse, &m.Speed, &m.Timestamp)
				if err != nil {
					http.Error(w, "Ошибка сканирования данных", http.StatusInternalServerError)
					return
				}
				metrics = append(metrics, m)
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(metrics)
		}
	})

	log.Printf("HTTP сервер запущен на %s", port)
	log.Fatal(http.ListenAndServe(port, r))
}
