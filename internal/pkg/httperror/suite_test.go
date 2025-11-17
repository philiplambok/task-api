package httperror

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestHttpError(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "HttpError Suite")
}
