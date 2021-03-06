package builds

import (
	"fmt"

	g "github.com/onsi/ginkgo"
	o "github.com/onsi/gomega"

	s2istatus "github.com/openshift/source-to-image/pkg/util/status"

	buildapi "github.com/openshift/origin/pkg/build/api"
	exutil "github.com/openshift/origin/test/extended/util"
)

var _ = g.Describe("[builds][Slow] update failure status", func() {
	defer g.GinkgoRecover()

	var (
		// convert the s2i failure cases to our own StatusReason
		reasonAssembleFailed  = buildapi.StatusReason(s2istatus.ReasonAssembleFailed)
		messageAssembleFailed = string(s2istatus.ReasonMessageAssembleFailed)
		postCommitHookFixture = exutil.FixturePath("testdata", "statusfail-postcommithook.yaml")
		fetchDockerSrc        = exutil.FixturePath("testdata", "statusfail-fetchsourcedocker.yaml")
		fetchS2ISrc           = exutil.FixturePath("testdata", "statusfail-fetchsources2i.yaml")
		builderImageFixture   = exutil.FixturePath("testdata", "statusfail-fetchbuilderimage.yaml")
		pushToRegistryFixture = exutil.FixturePath("testdata", "statusfail-pushtoregistry.yaml")
		failedAssembleFixture = exutil.FixturePath("testdata", "statusfail-failedassemble.yaml")
		binaryBuildDir        = exutil.FixturePath("testdata", "statusfail-assemble")
		oc                    = exutil.NewCLI("update-buildstatus", exutil.KubeConfigPath())
	)

	g.JustBeforeEach(func() {
		g.By("waiting for the builder service account")
		err := exutil.WaitForBuilderAccount(oc.KubeClient().ServiceAccounts(oc.Namespace()))
		o.Expect(err).NotTo(o.HaveOccurred())
	})

	g.Describe("Build status postcommit hook failure", func() {
		g.It("should contain the post commit hook failure reason and message", func() {
			err := oc.Run("create").Args("-f", postCommitHookFixture).Execute()
			o.Expect(err).NotTo(o.HaveOccurred())

			br, err := exutil.StartBuildAndWait(oc, "statusfail-postcommithook", "--build-loglevel=5")
			o.Expect(err).NotTo(o.HaveOccurred())
			br.AssertFailure()
			br.DumpLogs()

			build, err := oc.Client().Builds(oc.Namespace()).Get(br.Build.Name)
			o.Expect(err).NotTo(o.HaveOccurred())
			o.Expect(build.Status.Reason).To(o.Equal(buildapi.StatusReasonPostCommitHookFailed))
			o.Expect(build.Status.Message).To(o.Equal(buildapi.StatusMessagePostCommitHookFailed))
		})
	})

	g.Describe("Build status Docker fetch source failure", func() {
		g.It("should contain the Docker build fetch source failure reason and message", func() {
			err := oc.Run("create").Args("-f", fetchDockerSrc).Execute()
			o.Expect(err).NotTo(o.HaveOccurred())

			br, err := exutil.StartBuildAndWait(oc, "statusfail-fetchsourcedocker", "--build-loglevel=5")
			o.Expect(err).NotTo(o.HaveOccurred())
			br.AssertFailure()
			br.DumpLogs()

			build, err := oc.Client().Builds(oc.Namespace()).Get(br.Build.Name)
			o.Expect(err).NotTo(o.HaveOccurred())
			o.Expect(build.Status.Reason).To(o.Equal(buildapi.StatusReasonFetchSourceFailed))
			o.Expect(build.Status.Message).To(o.Equal(buildapi.StatusMessageFetchSourceFailed))
		})
	})

	g.Describe("Build status fetch S2I source failure", func() {
		g.It("should contain the S2I fetch source failure reason and message", func() {
			err := oc.Run("create").Args("-f", fetchS2ISrc).Execute()
			o.Expect(err).NotTo(o.HaveOccurred())

			br, err := exutil.StartBuildAndWait(oc, "statusfail-fetchsourcesourcetoimage", "--build-loglevel=5")
			o.Expect(err).NotTo(o.HaveOccurred())
			br.AssertFailure()
			br.DumpLogs()

			build, err := oc.Client().Builds(oc.Namespace()).Get(br.Build.Name)
			o.Expect(err).NotTo(o.HaveOccurred())
			o.Expect(build.Status.Reason).To(o.Equal(buildapi.StatusReasonFetchSourceFailed))
			o.Expect(build.Status.Message).To(o.Equal(buildapi.StatusMessageFetchSourceFailed))
		})
	})

	g.Describe("Build status fetch builder image failure", func() {
		g.It("should contain the fetch builder image failure reason and message", func() {
			err := oc.Run("create").Args("-f", builderImageFixture).Execute()
			o.Expect(err).NotTo(o.HaveOccurred())

			br, err := exutil.StartBuildAndWait(oc, "statusfail-builderimage", "--build-loglevel=5")
			o.Expect(err).NotTo(o.HaveOccurred())
			br.AssertFailure()
			br.DumpLogs()

			build, err := oc.Client().Builds(oc.Namespace()).Get(br.Build.Name)
			o.Expect(err).NotTo(o.HaveOccurred())
			o.Expect(build.Status.Reason).To(o.Equal(buildapi.StatusReasonPullBuilderImageFailed))
			o.Expect(build.Status.Message).To(o.Equal(buildapi.StatusMessagePullBuilderImageFailed))
		})
	})

	g.Describe("Build status push image to registry failure", func() {
		g.It("should contain the image push to registry failure reason and message", func() {
			err := oc.Run("create").Args("-f", pushToRegistryFixture).Execute()
			o.Expect(err).NotTo(o.HaveOccurred())

			br, err := exutil.StartBuildAndWait(oc, "statusfail-pushtoregistry", "--build-loglevel=5")
			o.Expect(err).NotTo(o.HaveOccurred())
			br.AssertFailure()
			br.DumpLogs()

			build, err := oc.Client().Builds(oc.Namespace()).Get(br.Build.Name)
			o.Expect(err).NotTo(o.HaveOccurred())
			o.Expect(build.Status.Reason).To(o.Equal(buildapi.StatusReasonPushImageToRegistryFailed))
			o.Expect(build.Status.Message).To(o.Equal(buildapi.StatusMessagePushImageToRegistryFailed))
		})
	})

	g.Describe("Build status failed assemble container", func() {
		g.It("should contain the failure reason related to an assemble script failing in s2i", func() {
			err := oc.Run("create").Args("-f", failedAssembleFixture).Execute()
			o.Expect(err).NotTo(o.HaveOccurred())

			br, err := exutil.StartBuildAndWait(oc, "statusfail-assemblescript", fmt.Sprintf("--from-dir=%s", binaryBuildDir), "--build-loglevel=5")
			o.Expect(err).NotTo(o.HaveOccurred())
			br.AssertFailure()
			br.DumpLogs()

			build, err := oc.Client().Builds(oc.Namespace()).Get(br.Build.Name)
			o.Expect(err).NotTo(o.HaveOccurred())
			o.Expect(build.Status.Reason).To(o.Equal(reasonAssembleFailed))
			o.Expect(build.Status.Message).To(o.Equal(messageAssembleFailed))
		})
	})
})
