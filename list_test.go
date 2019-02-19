package scan_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/pfs/pkg/scan"
)

var _ = Describe("ListImages", func() {
	var (
		res     string
		baseDir string
		images  []string
		err     error
	)

	BeforeEach(func() {
		wd, err := os.Getwd()
		Expect(err).NotTo(HaveOccurred())
		baseDir = filepath.Join(wd, "fixtures")
	})

	JustBeforeEach(func() {
		images, err = scan.ListImages(res, baseDir)
	})

	Context("when the resource file names an image directly", func() {
		BeforeEach(func() {
			res = "simple.yaml"
		})

		It("should list the image", func() {
			Expect(err).NotTo(HaveOccurred())
			Expect(images).To(ConsistOf("gcr.io/knative-releases/x/y"))
		})
	})

	Context("when the resource file names images as arguments", func() {
		BeforeEach(func() {
			res = "arg.yaml"
		})

		It("should list the images", func() {
			Expect(err).NotTo(HaveOccurred())
			Expect(images).To(ConsistOf("gcr.io/knative-releases/a/b", "gcr.io/knative-releases/c/d"))
		})
	})

	Context("when the resource file names an image using a key ending in 'Image'", func() {
		BeforeEach(func() {
			res = "suffix.yaml"
		})

		It("should list the images", func() {
			Expect(err).NotTo(HaveOccurred())
			Expect(images).To(ConsistOf("gcr.io/knative-releases/e/f", "k8s.gcr.io/fluentd-elasticsearch:v2.0.4"))
		})
	})

	Context("when the resource file contains block scalars containing images", func() {
		BeforeEach(func() {
			res = "block.yaml"
		})

		It("should list the images", func() {
			Expect(err).NotTo(HaveOccurred())
			Expect(images).To(ConsistOf("docker.io/istio/proxy_init:1.0.1"))
		})
	})

	Context("when using a realistic resource file", func() {
		BeforeEach(func() {
			res = "complex.yaml"
		})

		It("should list the images in the resource file", func() {
			Expect(err).NotTo(HaveOccurred())
			Expect(images).To(ConsistOf("gcr.io/knative-releases/github.com/knative/build/cmd/controller@sha256:3981b19105aabf3ed66db38c15407dc7accf026f4f4703d7e0ca7986ffd37d99",
				"gcr.io/knative-releases/github.com/knative/build/cmd/creds-init@sha256:b5dff24742c5c8ac4673dc991e3f960d11b58efdf751d26c54ec5144c48eef30",
				"gcr.io/knative-releases/github.com/knative/build/cmd/git-init@sha256:fe0d19e5da3fc9e7da20abc13d032beafcc283358a8325188dced62536a66e54",
				"gcr.io/knative-releases/github.com/knative/build/cmd/webhook@sha256:b9a97b7d360e10e540edfc9329e4f1c01832e58bf57d5dddea5c3a664f64bfc6",
				"docker.io/istio/proxyv2:0.8.0",
				"gcr.io/knative-releases/github.com/knative/serving/cmd/activator@sha256:e83258dd5858c8b1e92dbd413d0857ad2b22a7c4215ed911f256f68e2972f362",
				"gcr.io/knative-releases/github.com/knative/serving/cmd/autoscaler@sha256:76222399addc02454db9837ea3ff54bae29849168586051a9d0180daa2c1a805",
				"gcr.io/knative-releases/github.com/knative/serving/cmd/queue@sha256:99c841aa72c2928d1bf333348a848b5afb182715a2a0441da6282c86d4be807e",
				"k8s.gcr.io/fluentd-elasticsearch:v2.0.4",
				"gcr.io/knative-releases/github.com/knative/serving/cmd/controller@sha256:28db335f18cbd2a015fd218b9c7ce30b4366898fa3728a7f6dab6537991de028",
				"gcr.io/knative-releases/github.com/knative/serving/cmd/webhook@sha256:50ea89c48f8890fbe0cee336fc5cbdadcfe6884afbe5977db5d66892095b397d",
			))
		})
	})

	Context("when using a simple parameterized resource file", func() {
		BeforeEach(func() {
			res = "parameterized.yaml"
		})

		It("should list the images in the resource file", func() {
			Expect(err).NotTo(HaveOccurred())
			Expect(images).To(ConsistOf("packs/run",
				"packs/base",
				"projectriff/buildpack",
				"projectriff/buildpack",
				"packs/util",
				"packs/util",
			))
		})
	})

	Context("when using a more complex parameterized resource file", func() {
		BeforeEach(func() {
			res = "parameterized-2.yaml"
		})

		It("should list the image", func() {
			Expect(err).NotTo(HaveOccurred())
			Expect(images).To(ConsistOf("packs/run:v3alpha2", "projectriff/builder:0.2.0-snapshot-ci-63cd05079e1f", "ubuntu:18.04"))
		})
	})

	Context("when the resource file is not found", func() {
		BeforeEach(func() {
			res = "nosuch.yaml"
		})

		It("should return a suitable error", func() {
			Expect(os.IsNotExist(err)).To(BeTrue())
		})
	})

	Context("when the resource file contains invalid YAML", func() {
		BeforeEach(func() {
			res = "invalid.yaml"
		})

		It("should return a suitable error", func() {
			Expect(err).To(MatchError(HavePrefix("error parsing resource file")))
		})
	})
})
