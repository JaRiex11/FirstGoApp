package main

import (
	"fmt"
	"time"
)

func main() {
	cur_t := time.Now()
    nextYearNum := cur_t.Year() + 1
    nextYearDate := time.Date(nextYearNum, time.January, 1, 0, 0, 0, 0, cur_t.Location())

    diff := nextYearDate.Sub(cur_t)

    daysLeft := int(diff.Hours() / 24)

    fmt.Printf("Количество дней до Нового Года: %d", daysLeft)
}
