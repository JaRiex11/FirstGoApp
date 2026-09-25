package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func daysToNextYear(cur_t time.Time) int {
	nextYearNum := cur_t.Year() + 1
	nextYearDate := time.Date(nextYearNum, time.January, 1, 0, 0, 0, 0, cur_t.Location())

	diff := nextYearDate.Sub(cur_t)

	daysLeft := int(diff.Hours() / 24)
	return daysLeft
}

// Структуры для машиночитаемого ответа в формате JSON
type Response struct {
	DaysLeft int `json:"days_left"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

// HTTP-обработчик (принимает запрос, проверяет параметры, отдает JSON)
func handleDaysToNewYear(w http.ResponseWriter, r *http.Request) {
	// Устанавливаем заголовок ответа, что мы возвращаем именно JSON
	w.Header().Set("Content-Type", "application/json")

	// Проверяем, что метод именно GET
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Разрешен только метод GET"})
		return
	}

	// Получаем параметр date из URL
	dateStr := r.URL.Query().Get("date")
	var targetTime time.Time

	if dateStr == "" {
		// если запрос без указанной даты, берем текущее время сервера
		targetTime = time.Now()
	} else {
		// Явный формат передачи даты
		var err error
		targetTime, err = time.Parse("02-01-2006", dateStr)
		if err != nil {
			// некорректные данные обрабатываются, возвращаем статус 400
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Неверный формат даты. Используйте ДД-ММ-ГГГГ"})
			return
		}
	}

	// Собственно, считаем дни
	days := daysToNextYear(targetTime)

	// формируем ответ
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{DaysLeft: days})
}

func main() {
	// Регистрируем маршрут API
	http.HandleFunc("/api/days", handleDaysToNewYear)

	fmt.Println("Сервер запущен на порту :8080...")
	// Запускаем веб-сервер на порту 8080
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
