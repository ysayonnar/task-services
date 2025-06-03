package utils

import "strings"

func IsEmailValid(email string) bool {
	if len(email) < 3 || len(email) > 200 {
		return false
	}

	idx := strings.Index(email, "@")
	return idx > 0 && idx < len(email)-1
}
