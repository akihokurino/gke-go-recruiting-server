package util

import (
	"unicode/utf8"

	"gke-go-recruiting-server/domain"
)

func ValidateTextRange(s string, min int, max int) error {
	length := utf8.RuneCountInString(s)
	if length < min || max < length {
		return domain.NewBadRequestErr(domain.BadRequestMsg)
	}
	return nil
}
