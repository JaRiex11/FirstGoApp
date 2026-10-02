package main

import (
	"testing"
	"time"
)

func BenchmarkDaysToNextYear(b *testing.B) {
	inputs := []time.Time{
		time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2024, time.February, 29, 12, 0, 0, 0, time.UTC),
		time.Date(2023, time.December, 31, 23, 59, 59, 0, time.UTC),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = daysToNextYear(inputs[i%len(inputs)])
	}
}
