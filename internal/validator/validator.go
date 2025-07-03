package validator

import (
	"context"
	"regexp"
	"strings"
	"unicode/utf8"
)

type Validator interface {
	Valid(context.Context) Evaluator
}

var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type Evaluator map[string]string

// AddFieldError adds a validation error message for a specific field if it doesn't already exist.
func (e *Evaluator) AddFieldError(field, message string) {
	if *e == nil {
		*e = make(Evaluator)
	}

	if _, exists := (*e)[field]; !exists {
		(*e)[field] = message
	}
}

// CheckField adds a validation error for a field if the provided condition is false.
func (e *Evaluator) CheckField(ok bool, field, message string) {
	if !ok {
		e.AddFieldError(field, message)
	}
}

// NotBlank returns true if the string is not empty or whitespace only.
func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

// MaxChars returns true if the string contains at most max runes (Unicode characters).
func MaxChars(value string, max int) bool {
	return utf8.RuneCountInString(value) <= max
}

// MinChars returns true if the string contains at least min runes (Unicode characters).
func MinChars(value string, min int) bool {
	return utf8.RuneCountInString(value) >= min
}

// Matches returns true if the string matches the provided regular expression.
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}
