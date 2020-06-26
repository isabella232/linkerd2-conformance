package specs

import (
	"testing"

	"github.com/linkerd/linkerd2-conformance/specs/inject"
	"github.com/linkerd/linkerd2-conformance/specs/install"
	"github.com/linkerd/linkerd2-conformance/specs/uninstall"
	"github.com/linkerd/linkerd2-conformance/utils"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

func RunAllSpecs(t *testing.T) {

	c := utils.TestConfig
	globalControlPlane := c.GlobalControlPlane.Enable
	u := c.GlobalControlPlane.Uninstall
	// A single top-level wrapper Describe is required to prevent
	// the specs from being run in a random order
	// The Describe message is intentionally left empty
	// as it only serves to prevent randomization of specs
	_ = ginkgo.Describe("", func() {
		// The pre-flight and primary tests
		// have been organized into 2 separate
		// Describe blocks because the primary
		// Describe block holds the logic for
		// global/local control plane installation
		// We do not want that to conflict or interfere
		// with the pre-flight test Describe block

		// pre-flight checks
		// leave Describe message empty for the sake of
		// clarity in results. Test descriptions can be found in the
		// body of the spec runner function
		_ = ginkgo.Describe("", func() {
			_ = install.RunInstallSpec()
			if !globalControlPlane {
				_ = uninstall.RunUninstallSpec()
			}
		})

		// Run primary tests
		_ = ginkgo.Describe("", func() {
			h := utils.TestHelper

			// Install and uninstall a control plane
			// before and after each It block
			// This means, while writing a new test suite,
			// each major feature must be wrapped in its own
			// It block
			if !globalControlPlane {
				_ = ginkgo.BeforeEach(func() {
					utils.InstallLinkerdControlPlane(h)
				})

				_ = ginkgo.AfterEach(func() {
					utils.UninstallLinkerdControlPlane(h)
				})

			}
			_ = inject.RunInjectSpec()

			if globalControlPlane && u { // should always be at the bottom
				_ = uninstall.RunUninstallSpec()
			}
		})
	})

	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "linkerd2 conformance validation")
}
