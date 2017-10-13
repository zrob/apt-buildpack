package brats_test

import (
	"io/ioutil"
	"os/exec"
	"path/filepath"

	"github.com/cloudfoundry/libbuildpack"
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
		})

		It("runs", func() {
			PushApp(app)
			Expect(app.Stdout.String()).To(ContainSubstring("-----> Download go 1.9"))

			Expect(app.Stdout.String()).To(ContainSubstring("Installing apt packages"))
			Expect(app.GetBody("/")).To(ContainSubstring("Ascii: ASCII 6/4 is decimal 100, hex 64"))
		})
	})

	Context("deploying an app with an updated version of the same buildpack", func() {
		var bpName string
		BeforeEach(func() {
			bpName = "brats_apt_changing_" + cutlass.RandStringRunes(6)

			app = cutlass.New(filepath.Join(bpDir, "fixtures", "simple"))
			app.Buildpacks = []string{bpName + "_buildpack", "binary_buildpack"}
		})
		AfterEach(func() {
			Expect(cutlass.DeleteBuildpack(bpName)).To(Succeed())
		})

		FIt("prints useful warning message to stdout", func() {
			Expect(cutlass.CreateOrUpdateBuildpack(bpName, buildpacks.CachedFile)).To(Succeed())
			PushApp(app)
			Expect(app.Stdout.String()).ToNot(ContainSubstring("buildpack version changed from"))

			Expect(libbuildpack.CopyFile(buildpacks.CachedFile, filepath.Join("/tmp/buildpack29.zip"))).To(Succeed())
			Expect(ioutil.WriteFile("/tmp/VERSION", []byte("NewVerson"), 0644)).To(Succeed())
			Expect(exec.Command("zip", "-d", "/tmp/buildpack29.zip", "VERSION").Run()).To(Succeed())
			Expect(exec.Command("zip", "-j", "-u", "/tmp/buildpack29.zip", "/tmp/VERSION").Run()).To(Succeed())

			Expect(cutlass.CreateOrUpdateBuildpack(bpName, "/tmp/buildpack29.zip")).To(Succeed())
			PushApp(app)
			Expect(app.Stdout.String()).To(ContainSubstring("buildpack version changed from"))
		})
	})
})
