package architecture

import (
	"testing"

	"github.com/matthewmcnew/archtest"
)

func TestServiceShouldNotDependOnOther(t *testing.T) {
	archtest.Package(t, "storage/...").
		ShouldNotDependOn("api/...", "currency/...", "email/...")
}

func TestFileOrchestratorShouldNotDependOnStorageService(t *testing.T) {
	archtest.Package(t, "storage/orchestrator/...").ShouldNotDependOn("storage/email")
}

func TestDomainShouldNotDependOnAnything(t *testing.T) {
	archtest.Package(t, "storage/domain").ShouldNotDependOn("...")
}
