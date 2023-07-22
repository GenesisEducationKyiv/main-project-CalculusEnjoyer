package architecture

import (
	"testing"

	"github.com/matthewmcnew/archtest"
)

func TestServiceShouldNotDependOnOther(t *testing.T) {
	archtest.Package(t, "email/...").
		ShouldNotDependOn("api/...", "currency/...", "storage/...")
}

func TestEmailDispatcherShouldNotDependOnServiceImplementation(t *testing.T) {
	archtest.Package(t, "email/dispatcher/executor/...").ShouldNotDependOn("email/dispatcher")
}

func TestDomainShouldNotDependOnAnything(t *testing.T) {
	archtest.Package(t, "email/domain/...").ShouldNotDependOn("...")
}
