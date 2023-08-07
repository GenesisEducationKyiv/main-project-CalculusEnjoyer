package architecture

import (
	"testing"

	"github.com/matthewmcnew/archtest"
)

func TestServiceShouldNotDependOnOtherExceptProtoSchemas(t *testing.T) {
	archtest.Package(t, "api/...").
		Ignoring("storage/transport/proto", "currency/transport/proto", "email/transport/proto").
		ShouldNotDependOn("currency/...", "storage/...", "email/...")
}

func TestControllersShouldNotDependOnServiceImplementation(t *testing.T) {
	archtest.Package(t, "api/rest/...").ShouldNotDependOn("api/service/...")
}

func TestServicesShouldNotDependOnImplementation(t *testing.T) {
	archtest.Package(t, "api/service/...").ShouldNotDependOn("api/grpc/...")
}

func TestDomainShouldNotDependOnAnything(t *testing.T) {
	archtest.Package(t, "api/domain/...").Ignoring("time").ShouldNotDependOn("...")
}
