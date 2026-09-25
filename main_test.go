package main

import (
	"testing"
	"time"
)

func TestDaysToNextYear(t *testing.T) {
	// Определяем структуру для одной строки таблицы тестов
	type testCase struct {
		name     string    // Описание тест-кейса
		input    time.Time // Входное значение времени
		expected int       // Ожидаемое количество дней
	}

	// Набор тестовых сценариев, покрывающий обычные и граничные случаи
	tests := []testCase{
		// начало календарного года
		{
			name:     "Начало обычного года (1 января)",
			input:    time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC),
			expected: 365, // До 1 января следующего года ровно 365 дней
		},
		{
			name:     "Начало високосного года (1 января)",
			input:    time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC),
			expected: 366, // В високосном году на 1 день больше
		},

		// Конец календарного года
		{
			name:     "Последний день года (31 декабря)",
			input:    time.Date(2023, time.December, 31, 0, 0, 0, 0, time.UTC),
			expected: 1, // До 1 января остался 1 день
		},

		// Даты до и после 29 февраля в високосном году (2024)
		{
			name:     "До 29 февраля (28 февраля високосного года)",
			input:    time.Date(2024, time.February, 28, 0, 0, 0, 0, time.UTC),
			expected: 308, // 366 дней - 58 прошедших дней (31 в янв + 27 в фев)
		},
		{
			name:     "Сам день 29 февраля",
			input:    time.Date(2024, time.February, 29, 0, 0, 0, 0, time.UTC),
			expected: 307,
		},
		{
			name:     "После 29 февраля (1 марта високосного года)",
			input:    time.Date(2024, time.March, 1, 0, 0, 0, 0, time.UTC),
			expected: 306,
		},

		// Те же даты в невисокосном году (2023) для сравнения
		{
			name:     "28 февраля обычного года",
			input:    time.Date(2023, time.February, 28, 0, 0, 0, 0, time.UTC),
			expected: 307,
		},
		{
			name:     "1 марта обычного года",
			input:    time.Date(2023, time.March, 1, 0, 0, 0, 0, time.UTC),
			expected: 306,
		},
	}

	// Запуск таблицы тестов в цикле
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := daysToNextYear(tc.input)

			if actual != tc.expected {
				t.Errorf("Для даты %s ожидали %d дней, но получили %d",
					tc.input.Format("2006-01-02"), tc.expected, actual)
			}
		})
	}
}
