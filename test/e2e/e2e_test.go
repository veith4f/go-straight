package e2e

import (
	"os/exec"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/format"

	"github.com/veith4f/go-straight/test/utils"
)

var _ = Describe("go-straight", Ordered, func() {

	// do not truncate outputs
	format.MaxLength = 0

	SetDefaultEventuallyTimeout(2 * time.Second)
	SetDefaultEventuallyPollingInterval(time.Second)

	Context("go-straight", func() {
		It("should run successfully", func() {
			By("validating that the container executes successfully")
			verifyContainerRuns := func(g Gomega) {
				// Get the name of the controller-manager pod
				cmd := exec.Command("docker", "run",
					"node-647ee1368442ecd1a315c673.ps-xaas.io/pluscontainer/go-straight:latest",
					"/go-straight")

				runOutput, err := utils.Run(cmd)
				g.Expect(err).NotTo(HaveOccurred(), "Failed to run docker container")
				Expect(runOutput).NotTo(BeEmpty())
				Expect(runOutput).To(ContainSubstring("Fib(10): 89"))
			}
			Eventually(verifyContainerRuns).Should(Succeed())
		})
	})
})
