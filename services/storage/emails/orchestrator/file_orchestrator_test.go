package orchestrator

import (
	"encoding/csv"
	"log"
	"os"
	"storage/config"
	"storage/emails/messages"
	"testing"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"
)

func TestAddEmail(t *testing.T) {
	mockConf := loadTestConf("TEST_ADD_EMAIL_STORAGE_PATH", "../../.env.test")
	clearTestFile(mockConf)
	defer clearTestFile(mockConf)
	fileOrchestrator := NewFileOrchestrator(mockConf)
	testEmail := "test@test.com"
	testFile, err := os.Open(mockConf.EmailStoragePath)
	if err != nil {
		log.Fatalf(errors.Wrap(err, "Can not load test file").Error())
	}
	defer testFile.Close()
	csvReader := csv.NewReader(testFile)

	if err = fileOrchestrator.WriteEmail(messages.Email{Value: testEmail}); err != nil {
		t.Fatalf(`%q: %v`, "write email failed", err)
	}

	data, err := csvReader.ReadAll()
	if err != nil {
		log.Fatalf(errors.Wrap(err, "Can not load data from test file").Error())
	}
	if data[0][0] != testEmail {
		t.Fatalf(`%q: %q != %q`, "write email failed", data[0][0], testEmail)
	}
}

func TestGetAllEmails(t *testing.T) {
	realEmails := []string{
		"test1@gmail.com",
		"test2@gmail.com",
		"test3@gmail.com",
	}
	length := 3
	mockConf := loadTestConf("TEST_GET_EMAILS_STORAGE_PATH", "../../.env.test")
	fileOrchestrator := NewFileOrchestrator(mockConf)

	emails, err := fileOrchestrator.GetAllRecords()

	if err != nil {
		t.Fatalf(`%q: %q`, "can not get all emails", err)
	}
	if len(emails) != length {
		t.Fatalf(`%q: %q != %q`, "get emails get wrong count of emails", length, len(emails))
	}
	for i, email := range emails {
		if email.Value != realEmails[i] {
			t.Fatalf(`%q: %q != %q`, "got wrong email", length, len(emails))
		}
	}
}

func TestReturnErrorWithWrongConf(t *testing.T) {
	fileOrchestrator := NewFileOrchestrator(config.Config{EmailStoragePath: "wrong email path"})

	_, err := fileOrchestrator.GetAllRecords()
	if err == nil {
		t.Fatalf("err == nil wile getting all emails")
	}

	err = fileOrchestrator.WriteEmail(messages.Email{Value: "Test"})
	if err == nil {
		t.Fatalf("err == nil while adding email")
	}
}

func loadTestConf(path string, testEnvPath string) config.Config {
	conf := config.Config{}
	err := godotenv.Load(testEnvPath)
	if err != nil {
		log.Fatalf(errors.Wrap(err, "Can not load test config").Error())
	}

	conf.EmailStoragePath = os.Getenv(path)

	return conf
}

func clearTestFile(conf config.Config) {
	if err := os.Truncate(conf.EmailStoragePath, 0); err != nil {
		log.Fatalf(errors.Wrap(err, "Can not load test config").Error())
	}
}
