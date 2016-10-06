package out_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"

	"github.com/idahobean/cf-sb-resource"
	"github.com/idahobean/cf-sb-resource/out"
)

var _ = Describe("Out", func() {
	var (
		tmpDir  string
		cmd     *exec.Cmd
		request out.Request
	)

	BeforeEach(func() {
		var err error

		tmpDir, err = ioutil.TempDir("", "cf_resource_out")
		Ω(err).ShouldNot(HaveOccurred())

		request = out.Request{
			Source: resource.Source{
				API:           "https://api.run.pivotal.io",
				Username:      "awesome@example.com",
				Password:      "hunter2",
				Organization:  "org",
				Space:         "space",
				SkipCertCheck: true,
			},
			Params: out.Params{
				Repository:     "foobar/foofoo:latest",
				CurrentAppName: "foobar-app",
				Memory:         "2G",
				Disk:           "256M",
				HealthCheck:    "none",
			},
		}
	})

	JustBeforeEach(func() {
		assetsPath, err := filepath.Abs("assets")
		Ω(err).ShouldNot(HaveOccurred())

		stdin := &bytes.Buffer{}

		err = json.NewEncoder(stdin).Encode(request)
		Ω(err).ShouldNot(HaveOccurred())

		cmd = exec.Command(binPath, tmpDir)
		cmd.Stdin = stdin
		cmd.Dir = tmpDir

		newEnv := []string{}
		for _, envVar := range os.Environ() {
			if strings.HasPrefix(envVar, "PATH=") {
				newEnv = append(newEnv, fmt.Sprintf("PATH=%s:%s", assetsPath, os.Getenv("PATH")))
			} else {
				newEnv = append(newEnv, envVar)
			}
		}

		cmd.Env = newEnv
	})

	AfterEach(func() {
		err := os.RemoveAll(tmpDir)
		Ω(err).ShouldNot(HaveOccurred())
	})

	Context("when option is fulfilled", func() {
		It("pushes an application to cloud foundry", func() {
			session, err := gexec.Start(
				cmd,
				GinkgoWriter,
				GinkgoWriter,
			)
			Ω(err).ShouldNot(HaveOccurred())

			Eventually(session).Should(gexec.Exit(0))

			var response out.Response
			err = json.Unmarshal(session.Out.Contents(), &response)
			Ω(err).ShouldNot(HaveOccurred())

			Ω(response.Version.Timestamp).Should(BeTemporally("~", time.Now(), time.Second))

			// shim outputs arguments
			Ω(session.Err).Should(gbytes.Say("cf api https://api.run.pivotal.io --skip-ssl-validation"))
			Ω(session.Err).Should(gbytes.Say("cf auth awesome@example.com hunter2"))
			Ω(session.Err).Should(gbytes.Say("cf target -o org -s space"))
			Ω(session.Err).Should(gbytes.Say("cf push foobar-app -o foobar/foofoo:latest -m 2G -k 256M -u none"))
			// color should be always
			Ω(session.Err).Should(gbytes.Say("CF_COLOR=true"))
		})
	})

	Context("when required option is empty", func() {
		Context("repository is empty", func() {
			BeforeEach(func() {
				request.Params.Repository = ""
			})

			It("return an error", func() {
				session, err := gexec.Start(
					cmd,
					GinkgoWriter,
					GinkgoWriter,
				)
				Ω(err).ShouldNot(HaveOccurred())

				Eventually(session).Should(gexec.Exit(1))

				errMsg := fmt.Sprintf("error parameter required: repository")
				Ω(session.Err).Should(gbytes.Say(errMsg))
			})
		})

		Context("currentAppName is empty", func() {
			BeforeEach(func() {
				request.Params.CurrentAppName = ""
			})

			It("return an error", func() {
				session, err := gexec.Start(
					cmd,
					GinkgoWriter,
					GinkgoWriter,
				)
				Ω(err).ShouldNot(HaveOccurred())

				Eventually(session).Should(gexec.Exit(1))

				errMsg := fmt.Sprintf("error parameter required: current_app_name")
				Ω(session.Err).Should(gbytes.Say(errMsg))
			})
		})
	})
})
