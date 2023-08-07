package validator

import "testing"

func TestNewRegexValidator(t *testing.T) {
	validator := NewRegexValidator(*DefaultEmailRegex)

	var tests = []struct {
		email  string
		result bool
	}{
		{"test@gmail.com", true},
		{"testgmail.com", false},
		{"test@gmail", false},
		{"@gmail.com", false},
		{"", false},
	}

	for _, testEmails := range tests {
		t.Run(testEmails.email, func(t *testing.T) {
			validationResult := validator.Validate(testEmails.email)
			if validationResult != testEmails.result {
				t.Errorf("email: %s got %t, want %t", testEmails.email, validationResult, testEmails.result)
			}
		})
	}
}
