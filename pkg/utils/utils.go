package utils

import (
	"go-parser-lambda/pkg/shared"
)

func FindMaxDate(markets []shared.Product) string {

	// HashTable for max date
	hashTableDate := make(map[string]int)
	var maxDate string
	maxCount := 0

	for _, market := range markets {
		resDate := market.Date

		// Counter dates
		if _, ok := hashTableDate[resDate]; ok {
			hashTableDate[resDate] = hashTableDate[resDate] + 1
		} else {
			hashTableDate[resDate] = 1
		}
	}

	for date, count := range hashTableDate {
		if count > maxCount {
			maxCount = count
			maxDate = date
		}
	}

	return maxDate
}
