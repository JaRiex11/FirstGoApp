package main

import (
	"fmt"
	"time"
)

func daysToNextYear(cur_t time.Time) int {
	nextYearNum := cur_t.Year() + 1
	nextYearDate := time.Date(nextYearNum, time.January, 1, 0, 0, 0, 0, cur_t.Location())

	diff := nextYearDate.Sub(cur_t)

	daysLeft := int(diff.Hours() / 24)
	return daysLeft
}

func main() {
	fmt.Print("Введите дату от которой считать дни (в формате ДД-ММ-ГГГГ): ")
	timeLayout := "02-01-2006"
	var input string
	fmt.Scan(&input)

	cur_t, err := time.Parse(timeLayout, input) // парсинг значения времени по шаблону
	if err != nil {
		panic(err)
	} // обработка возможной ошибки парсинга

	var daysLeft int
	daysLeft = daysToNextYear(cur_t)

	fmt.Printf("Количество дней до ближайшего Нового Года: %d", daysLeft)
}
