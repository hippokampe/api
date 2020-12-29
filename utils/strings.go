package utils

import "strings"

func CleanString(str string) string {
	return strings.Trim(str, "\t\n ")
}
