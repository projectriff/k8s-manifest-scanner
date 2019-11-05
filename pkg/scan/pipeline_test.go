package scan_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"io/ioutil"

	. "github.com/projectriff/k8s-manifest-scanner/pkg/scan"
	"gopkg.in/yaml.v3"
)

var _ = Describe("MatchImageKey", func() {
	It("matches map nodes with image key", func() {
		res := MatchImageKey(toYAML("image: hello")).ToArray()

		Expect(res).To(HaveLen(1))
		Expect(res[0].Value).To(Equal("hello"))
	})

	It("matches map nodes with an '-image' suffix", func() {
		res := MatchImageKey(toYAML("some-image: hello")).ToArray()

		Expect(res).To(HaveLen(1))
		Expect(res[0].Value).To(Equal("hello"))
	})

	It("matches map nodes with an 'Image' suffix", func() {
		res := MatchImageKey(toYAML("someImage: hello")).ToArray()

		Expect(res).To(HaveLen(1))
		Expect(res[0].Value).To(Equal("hello"))
	})

	It("ignores image values that are parameterized (starts with $)", func() {
		res := MatchImageKey(toYAML("image: $hello")).ToArray()
		Expect(res).To(BeEmpty())
	})

	It("ignores images values that aren't strings", func() {
		res := MatchImageKey(toYAML("image: 1.0")).ToArray()
		Expect(res).To(BeEmpty())
	})

	It("doesn't match arrays", func() {
		res := MatchImageKey(toYAML("[]")).ToArray()
		Expect(res).To(BeEmpty())
	})
})

var _ = Describe("MatchArgsMap", func() {
	It("matches image flags in 'args' sequences", func() {
		doc := loadYAML("testdata/arg.yaml")
		res := MatchArgsMap(doc).ToArray()

		Expect(res).To(HaveLen(2))
		Expect(res[0].Value).To(Equal("gcr.io/knative-releases/a/b"))
		Expect(res[1].Value).To(Equal("gcr.io/knative-releases/c/d"))
	})
})

var _ = Describe("MatchTemplateDefaults", func() {
	It("matches template parameter defaults", func() {
		doc := loadYAML("testdata/parameterized-2.yaml")
		res := MatchTemplateDefaults(doc).ToArray()

		Expect(res).To(HaveLen(2))
		Expect(res[0].Value).To(Equal("packs/run:v3alpha2"))
		Expect(res[1].Value).To(Equal("projectriff/builder:0.2.0-snapshot-ci-63cd05079e1f"))
	})
})

func loadYAML(path string) *yaml.Node {
	bytes, err := ioutil.ReadFile(path)
	Expect(err).ToNot(HaveOccurred())
	return toYAML(string(bytes))
}

func toYAML(s string) *yaml.Node {
	var node yaml.Node

	err := yaml.Unmarshal([]byte(s), &node)
	Expect(err).ToNot(HaveOccurred())

	return &node
}
