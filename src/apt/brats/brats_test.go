package brats_test

import (
	"path/filepath"

	"github.com/cloudfoundry/libbuildpack/cutlass"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Apt supply buildpack", func() {
	var app *cutlass.App
	AfterEach(func() { app = DestroyApp(app) })

	Context("Unbuilt buildpack (eg github)", func() {
		BeforeEach(func() {
			app = cutlass.New(filepath.Join(bpDir, "fixtures", "simple"))
			app.Buildpacks = []string{buildpacks.Unbuilt, "binary_buildpack"}
			app.SetEnv("BP_DEBUG", "1")
		})

		It("runs", func() {
			PushApp(app)
			Expect(app.Stdout.String()).To(ContainSubstring("-----> Download go 1.9"))

			Expect(app.Stdout.String()).To(ContainSubstring("Installing apt packages"))
			Expect(app.GetBody("/")).To(ContainSubstring("Ascii: ASCII 6/4 is decimal 100, hex 64"))
		})
	})

	// Context("as a supply buildpack", func() {
	// 	BeforeEach(func() {
	// 		app = cutlass.New(filepath.Join(bpDir, "fixtures", "simple"))
	// 		app.Buildpacks = []string{buildpacks.Cached, "binary_buildpack"}
	// 		app.SetEnv("BP_DEBUG", "1")
	// 	})

	// 	It("supplies apt packages to later buildpacks", func() {
	// 		PushApp(app)

	// 		Expect(app.Stdout.String()).To(ContainSubstring("Installing apt packages"))
	// 		Expect(app.GetBody("/")).To(ContainSubstring("Ascii: ASCII 6/4 is decimal 100, hex 64"))
	// 	})
	// })
})
