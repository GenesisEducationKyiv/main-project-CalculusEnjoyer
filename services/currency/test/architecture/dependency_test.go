package architecture

import (
	"testing"

	"github.com/matthewmcnew/archtest"
)

func TestServiceShouldNotDependOnOther(t *testing.T) {
	archtest.Package(t, "currency/...").
		ShouldNotDependOn("api/...", "storage/...", "email/...")
}

func TestRateProvidersShouldNotDependOnServiceImplementation(t *testing.T) {
	archtest.Package(t, "currency/rate/providers/crypto/...").ShouldNotDependOn("currency/rate")
}

func TestDomainShouldNotDependOnAnything(t *testing.T) {
	archtest.Package(t, "currency/domain/...").Ignoring("time").ShouldNotDependOn("...")
}
