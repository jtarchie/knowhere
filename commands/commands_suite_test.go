package commands_test

import (
	"log/slog"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCommands(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Commands Suite")
}

var _ = BeforeSuite(func() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(GinkgoWriter, nil)))
})
