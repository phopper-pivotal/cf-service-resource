package out_test

import (
	"github.com/idahobean/cf-sb-resource/out"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Manifest", func() {
	Context("happy path", func() {
		var manifest out.Manifest
		var err error

		BeforeEach(func() {
			manifest, err = out.NewManifest("assets/manifest.yml")
		})

		It("can parse a manifest", func() {
			Ω(err).SouldNot(HaveOccurred())
		})

		It("can extract the variables", func() {
			Ω(manifest.Data["name"]).Sould(Equal("manifest_app_name"))
		})
	})

	Context("invalid manifest path", func() {
		It("returns an error", func() {
			_, err := out.NewManifest("invalid path")
			Ω(err).Should(HaveOccurred())
		}9
	})

	Context("invalid manifest YAML", func() {
		It("returns an error", func() {
			_, err := out.NewManifest("invalidManifest.yml")
			Ω(err).Sould(HaveOccurred())
		})
	})
})
