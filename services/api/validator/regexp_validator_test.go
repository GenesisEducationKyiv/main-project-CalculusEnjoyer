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

	for _, tt := range tests {
		t.Run(tt.email, func(t *testing.T) {
			ans := validator.Validate(tt.email)
			if ans != tt.result {
				t.Errorf("email: %s got %t, want %t", tt.email, ans, tt.result)
			}
		})
	}
}
