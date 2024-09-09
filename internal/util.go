package internal

import (
	"fmt"
	"strconv"
	"time"
)

func GetNewCode(prefix string, dateStr string, currentCode string) (string, error) {
	// Define the length of the numeric part
	numericPartLength := 4

	// Check the date format
	// Parse the string using the full RFC3339 format
	parsedDate, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return "", fmt.Errorf("GetNewCode invalid date format: %v", err)
	}
	formattedDate := parsedDate.Format("20060102")
	var newCode string
	if currentCode == "" {
		// Create a new code
		newCode = fmt.Sprintf("%s%s%0*d", prefix, formattedDate, numericPartLength, 1)
	} else {
		// Extract the numeric part from the current code
		if len(currentCode) <= len(prefix)+len(formattedDate) {
			return "", fmt.Errorf("current code is too short")
		}
		extractedPrefix := currentCode[:len(prefix)]
		currentNumericPart := currentCode[len(prefix)+len(formattedDate):]

		// Ensure the prefix and date match
		if extractedPrefix != prefix {
			return "", fmt.Errorf("prefix or date mismatch")
		}

		// Parse the numeric part as an integer
		num, err := strconv.Atoi(currentNumericPart)
		if err != nil {
			return "", fmt.Errorf("failed to parse numeric part: %v", err)
		}

		// Increment the numeric part
		num++

		// Format the new numeric part with leading zeros
		newNumericPart := fmt.Sprintf("%0*d", numericPartLength, num)

		// Combine the prefix, date, and new numeric part
		newCode = fmt.Sprintf("%s%s%s", prefix, formattedDate, newNumericPart)
	}

	return newCode, nil
}

func Find[T any](slice []T, condition func(T) bool) *T {
	for i := range slice {
		if condition(slice[i]) {
			return &slice[i]
		}
	}
	return nil
}
