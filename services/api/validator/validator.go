package validator

type EmailValidator interface {
	Validate(email string) bool
}
