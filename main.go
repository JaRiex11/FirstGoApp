package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"time"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

type HealthResponse struct {
	Status string `json:"status"`
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(rec, r)

		slog.Info("http request",
			"method", r.Method,
			"path", r.URL.Path,
			"query", r.URL.RawQuery,
			"status", rec.status,
			"duration_ms", time.Since(start).Milliseconds(),
			"remote", r.RemoteAddr,
		)
	})
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: "Разрешен только метод GET"})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(HealthResponse{Status: "ok"})
}

func daysToNextYear(cur_t time.Time) int {
	// Нормализуем до начала дня, чтобы считать именно календарные дни
	cur_t = time.Date(cur_t.Year(), cur_t.Month(), cur_t.Day(), 0, 0, 0, 0, cur_t.Location())

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

// routes возвращает HTTP-обработчик со всеми маршрутами.
// Это удобно использовать и в main, и в интеграционных тестах.
func routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/days", handleDaysToNewYear) // Регистрируем маршрут API
	mux.HandleFunc("/healthz", handleHealth)         // Дополнительный endpoint
	return loggingMiddleware(mux)                    // Обернули для логирования
}

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	addr := ":8080"
	slog.Info("server starting", "addr", addr)

	if err := http.ListenAndServe(addr, routes()); err != nil {
		slog.Error("server failed", "err", err)
	}
}
