package emails

import (
	"log"
	"os"
	"storage/config"
	"storage/emails/messages"
	"storage/emails/orchestrator"
	"storage/serror"
	"testing"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"
)

const (
	test1 = "test1@gmail.com"
	test2 = "test2@gmail.com"
	test3 = "test3@gmail.com"
)

func TestGetAllEmails(t *testing.T) {
	resetTestFile(t)
	service := getFileService()
	expectedEmails := []string{test1, test2, test3}

	actualEmails, err := service.GetAllEmails()
	if err != nil {
		t.Fatalf(errors.Wrap(err, "can not get all emails").Error())
	}

	for i, email := range expectedEmails {
		if actualEmails[i].Value != email {
			t.Fatalf(`%q: %q != %q`, "got wrong email", email, actualEmails[i].Value)
		}
	}
	resetTestFile(t)
}

func TestAddEmail(t *testing.T) {
	const emailToAdd = "iWantToBeAdded@please.com"
	resetTestFile(t)
	service := getFileService()

	err := service.AddEmail(messages.Email{Value: emailToAdd})
	if err != nil {
		t.Fatalf(errors.Wrap(err, "can not add email").Error())
	}

	wroteEmails, err := service.GetAllEmails()
	if err != nil {
		t.Fatalf(errors.Wrap(err, "can not get all emails").Error())
	}

	if !contains(wroteEmails, emailToAdd) {
		t.Fatalf(errors.Wrap(err, "email was not added").Error())
	}
	resetTestFile(t)
}

func TestAddSameEmailTwice(t *testing.T) {
	const emailToAdd = "iWantToBeAddedOnlyOnce@please.com"
	resetTestFile(t)
	service := getFileService()

	err := service.AddEmail(messages.Email{Value: emailToAdd})
	if err != nil {
		t.Fatalf(errors.Wrap(err, "can not add email").Error())
	}
	err = service.AddEmail(messages.Email{Value: emailToAdd})
	if !errors.Is(err, serror.ErrEmailAlreadyExists) {
		t.Fatalf(errors.Wrap(err, "adding the same email twice did not produce an error").Error())
	}

	wroteEmails, err := service.GetAllEmails()
	if err != nil {
		t.Fatalf(errors.Wrap(err, "can not get all emails").Error())
	}
	amount := containsAmount(wroteEmails, emailToAdd)
	if amount != 1 {
		t.Fatalf(`%q: %q`, "added an email not once but", amount)
	}
	resetTestFile(t)
}

func getFileService() StorageService {
	conf := loadTestConf()
	return NewService(orchestrator.NewFileOrchestrator(conf))
}

func resetTestFile(t *testing.T) {
	conf := loadTestConf()
	if err := os.Truncate(conf.EmailStoragePath, 0); err != nil {
		t.Errorf(errors.Wrap(err, "can not reset test file").Error())
	}

	fileOrchestrator := orchestrator.NewFileOrchestrator(conf)
	emails := []string{test1, test2, test3}
	for _, email := range emails {
		err := fileOrchestrator.WriteEmail(messages.Email{Value: email})
		if err != nil {
			t.Errorf(errors.Wrap(err, "can not reset test file").Error())
		}
	}
}

func loadTestConf() config.Config {
	conf := config.Config{}
	err := godotenv.Load("../.env.test")
	if err != nil {
		log.Fatalf(errors.Wrap(err, "Can not load test config").Error())
	}

	conf.EmailStoragePath = os.Getenv("TEST_FILE_SERVICE_INTEGRATION")

	return conf
}

func contains(emails []messages.Email, target string) bool {
	for _, v := range emails {
		if v.Value == target {
			return true
		}
	}

	return false
}

func containsAmount(emails []messages.Email, target string) int {
	base := 0
	for _, v := range emails {
		if v.Value == target {
			base++
		}
	}

	return base
}
