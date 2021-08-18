package util

import (
	"regexp"
	"strings"
)

func normalizeCommand(command string) string {
	pattern := regexp.MustCompile(`[\s.,]+`)
	return pattern.ReplaceAllString(
		strings.TrimSpace(strings.ToLower(command)),
		"",
	)
}

func removeDiagonals(command string) string {
	diagonals := regexp.MustCompile(`[1379]+`)
	return diagonals.ReplaceAllString(command, "")
}

func normalizeCompare(i string, j string) bool {
	return normalizeCommand(i) == normalizeCommand(j)
}

func sameMove(command1 string, command2 string) bool {
	normalizedEqual := normalizeCompare(command1, command2)
	// Try again with diagonals removed for HCF motions
	if !normalizedEqual && len(command1) > 3 && len(command2) > 3 {
		normalizedEqual = normalizeCompare(removeDiagonals(command1), removeDiagonals(command2))
	}
	return normalizedEqual
}
