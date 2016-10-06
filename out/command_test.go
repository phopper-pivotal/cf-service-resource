package out_test

import (
	"errors"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/idahobean/cf-sb-resource"
	"github.com/idahobean/cf-sb-resource/out"
	"github.com/idahobean/cf-sb-resource/out/fakes"
)

var _ = Describe("Out Command", func() {
	var (
		cloudFoundry *fakes.FakePAAS
		request      out.Request
		command      *out.Command
	)

	BeforeEach(func() {
		cloudFoundry = &fakes.FakePAAS{}
		command = out.NewCommand(cloudFoundry)

		request = out.Request{
			Source: resource.Source{
				API:          "https://api.run.pivotal.io",
				Username:     "awesome@example.com",
				Password:     "hunter2",
				Organization: "secret",
				Space:        "volcano-base",
			},
			Params: out.Params{
				Repository: "idahobean/dummy:latest",
				CurrentAppName: "foobar",
			},
		}
	})

	Describe("running the command", func() {
		It("pushes an application into cloud foundry", func() {
			response, err := command.Run(request)
			Ω(err).ShouldNot(HaveOccurred())

			Ω(response.Version.Timestamp).Should(BeTemporally("~", time.Now(), time.Second))
			Ω(response.Metadata[0]).Should(Equal(
				resource.MetadataPair{
					Name:  "organization",
					Value: "secret",
				},
			))
			Ω(response.Metadata[1]).Should(Equal(
				resource.MetadataPair{
					Name:  "space",
					Value: "volcano-base",
				},
			))

			By("logging in")
			Ω(cloudFoundry.LoginCallCount()).Should(Equal(1))

			api, username, password, insecure := cloudFoundry.LoginArgsForCall(0)
			Ω(api).Should(Equal("https://api.run.pivotal.io"))
			Ω(username).Should(Equal("awesome@example.com"))
			Ω(password).Should(Equal("hunter2"))
			Ω(insecure).Should(Equal(false))

			By("targetting the organization and space")
			Ω(cloudFoundry.TargetCallCount()).Should(Equal(1))

			org, space := cloudFoundry.TargetArgsForCall(0)
			Ω(org).Should(Equal("secret"))
			Ω(space).Should(Equal("volcano-base"))

			By("pushing the app")
			Ω(cloudFoundry.PushAppCallCount()).Should(Equal(1))

			repository, currentAppName, memory, disk, healthCheck := cloudFoundry.PushAppArgsForCall(0)

			Ω(repository).Should(Equal("idahobean/dummy:latest"))
			Ω(currentAppName).Should(Equal("foobar"))
			Ω(memory).Should(Equal(""))
			Ω(disk).Should(Equal(""))
			Ω(healthCheck).Should(Equal(""))

		})

		Describe("handling any errors", func() {
			var expectedError error

			BeforeEach(func() {
				expectedError = errors.New("it all went wrong")
			})

			It("from logging in", func() {
				cloudFoundry.LoginReturns(expectedError)

				_, err := command.Run(request)
				Ω(err).Should(MatchError(expectedError))
			})

			It("from targetting an org and space", func() {
				cloudFoundry.TargetReturns(expectedError)

				_, err := command.Run(request)
				Ω(err).Should(MatchError(expectedError))
			})

			It("from pushing the application", func() {
				cloudFoundry.PushAppReturns(expectedError)

				_, err := command.Run(request)
				Ω(err).Should(MatchError(expectedError))
			})
		})

		It("lets people skip the certificate check", func() {
			request = out.Request{
				Source: resource.Source{
					API:           "https://api.run.pivotal.io",
					Username:      "awesome@example.com",
					Password:      "hunter2",
					Organization:  "secret",
					Space:         "volcano-base",
					SkipCertCheck: true,
				},
				Params: out.Params{
					Repository: "idahobean/dummy",
					CurrentAppName: "fooobaaar",
				},
			}

			_, err := command.Run(request)
			Ω(err).ShouldNot(HaveOccurred())

			By("logging in")
			Ω(cloudFoundry.LoginCallCount()).Should(Equal(1))

			_, _, _, insecure := cloudFoundry.LoginArgsForCall(0)
			Ω(insecure).Should(Equal(true))
		})

		It("deploy with fulfilled options", func() {
			request = out.Request{
				Source: resource.Source{
					API:          "https://api.foobar.cfapps.io",
					Username:     "foo@bar.baz",
					Password:     "foobarbaz",
					Organization: "orgg",
					Space:        "spacee",
				},
				Params: out.Params{
					Repository:     "idahobean/foo",
					CurrentAppName: "bar",
					Memory:         "64M",
					Disk:           "128M",
					HealthCheck:    "none",
				},
			}

			_, err := command.Run(request)
			Ω(err).ShouldNot(HaveOccurred())

			By("pushing the app")
			Ω(cloudFoundry.PushAppCallCount()).Should(Equal(1))

			repository, currentAppName, memory, disk, healthCheck := cloudFoundry.PushAppArgsForCall(0)
			Ω(repository).Should(Equal("idahobean/foo"))
			Ω(currentAppName).Should(Equal("bar"))
			Ω(memory).Should(Equal("64M"))
			Ω(disk).Should(Equal("128M"))
			Ω(healthCheck).Should(Equal("none"))
		})
	})
})
