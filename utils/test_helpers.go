package utils

import (
	"encoding/json"
	"fmt"

	"github.com/linkerd/linkerd2/testutil"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

type CheckOutput struct {
	Success    bool `json:"success"`
	Categories []struct {
		CategoryName string `json:"categoryName"`
		Checks       []struct {
			Result string `json:"result"`
			Error  string `json:"error"`
		}
	}
}

func getFailedChecks(r *CheckOutput) string {
	err := "The following errors were detected:\n"

	for _, c := range r.Categories {
		for _, check := range c.Checks {
			if check.Result == "error" {
				err = fmt.Sprintf("%s\n%s", err, check.Error)
			}
		}
	}

	return err
}

// RunCheck rus `linkerd check`
func RunCheck(h *testutil.TestHelper, pre bool) {

	var checkResult *CheckOutput

	cmd := []string{
		"check",
		"--expected-version",
		h.GetVersion(),
		"-o",
		"json",
	}

	if pre {
		cmd = append(cmd, "--pre")
		ginkgo.By("Running pre-installation checks")
	} else {
		ginkgo.By("Running post-installation checks")
	}

	out, stderr, err := h.LinkerdRun(cmd...)
	gomega.Expect(err).Should(gomega.BeNil(), stderr)

	ginkgo.By("Validating `check` output")
	err = json.Unmarshal([]byte(out), &checkResult)
	gomega.Expect(err).Should(gomega.BeNil(), Err(err))

	gomega.Expect(checkResult.Success).Should(gomega.BeTrue(), getFailedChecks(checkResult))
}

func InstallLinkerdControlPlane(h *testutil.TestHelper, withHA bool) {
	ginkgo.By(fmt.Sprintf("Installing linkerd control plane with HA: %v", withHA))
	RunCheck(h, true) // run pre checks

	if err := h.CheckIfNamespaceExists(h.GetLinkerdNamespace()); err == nil {
		ginkgo.Skip(fmt.Sprintf("linkerd control plane already exists in namespace %s", h.GetLinkerdNamespace()))
	}

	// TODO: Uncomment while writing Helm tests
	// ginkgo.By("verifying if Helm release is empty")
	// gomega.Expect(h.GetHelmReleaseName()).To(gomega.Equal(""))

	cmd := "install"
	args := []string{
		"--controller-log-level", "debug",
		"--proxy-log-level", "warn,linkerd2_proxy=debug",
		"--proxy-version", h.GetVersion(),
	}

	if withHA {
		args = append(args, "--ha")
	}

	if h.GetClusterDomain() != "cluster.local" {
		args = append(args, "--cluster-domain", h.GetClusterDomain())
	}

	exec := append([]string{cmd}, args...)

	ginkgo.By("Running `linkerd install`")
	out, stderr, err := h.LinkerdRun(exec...)
	gomega.Expect(err).Should(gomega.BeNil(), stderr)

	ginkgo.By("Applying control plane manifests")
	out, err = h.KubectlApply(out, "")
	gomega.Expect(err).Should(gomega.BeNil(), Err(err))

	RunCheck(h, false) // run post checks
}

func UninstallLinkerdControlPlane(h *testutil.TestHelper) {
	ginkgo.By("Uninstalling linkerd control plane")
	cmd := "install"
	args := []string{
		"--ignore-cluster",
	}

	exec := append([]string{cmd}, args...)

	ginkgo.By("Gathering control plane manifests")
	out, stderr, err := h.LinkerdRun(exec...)
	gomega.Expect(err).Should(gomega.BeNil(), stderr)

	args = []string{"delete", "-f", "-"}

	ginkgo.By("Deleting resources from the cluster")
	out, err = h.Kubectl(out, args...)
	gomega.Expect(err).Should(gomega.BeNil(), Err(err))

	RunCheck(h, true) // run pre checks
}
