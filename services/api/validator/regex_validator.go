package validator

import (
	"regexp"
)

var DefaultEmailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

type RegexEmailValidator struct {
	regexp regexp.Regexp
}

func NewRegexValidator(regexp regexp.Regexp) *RegexEmailValidator {
	return &RegexEmailValidator{
		regexp: regexp,
	}
}

func (v *RegexEmailValidator) Validate(email string) bool {
	return v.regexp.MatchString(email)
}
