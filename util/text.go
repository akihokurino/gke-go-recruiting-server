package util

import (
	"unicode/utf8"

	"gke-go-sample/domain"
)

func ValidateTextRange(s string, min int, max int) error {
	length := utf8.RuneCountInString(s)
	if length < min || max < length {
		return domain.NewBadRequestErr(domain.BadRequestMsg)
	}
	return nil
}
