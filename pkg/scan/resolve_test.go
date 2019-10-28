package scan

import (
	"fmt"
	"log"
	"net/http/httptest"
	"strings"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/registry"
	"github.com/google/go-containerregistry/pkg/v1/random"
	"github.com/google/go-containerregistry/pkg/v1/remote"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ResolveImagesFromKubernetesManifest", func() {
	var server *httptest.Server
	var digestRef string
	var tagRef string

	BeforeEach(func() {
		logger := log.New(GinkgoWriter, "[registry] ", log.LstdFlags)

		r := registry.New(registry.Logger(logger))
		server = httptest.NewServer(r)

		tagRef = strings.TrimPrefix(server.URL, "http://") + "/foo:latest"
		tag, err := name.NewTag(tagRef)
		Expect(err).ToNot(HaveOccurred(), "unable to create tag")

		i, err := random.Image(1024, 1)
		Expect(err).ToNot(HaveOccurred(), "unable to make random image")

		err = remote.Write(tag, i)
		Expect(err).ToNot(HaveOccurred(), "unable to upload random image")

		ri, err := remote.Image(tag)
		Expect(err).ToNot(HaveOccurred(), "unable read image")

		digest, err := ri.Digest()
		Expect(err).ToNot(HaveOccurred(), "get digest")

		digestRef = fmt.Sprintf("%v@%v", tag.Repository, digest)

	})

	AfterEach(func() {
		server.Close()
	})

	It("resolves image tags to digest", func() {
		yaml := "image: " + tagRef
		bytes, err := resolveImagesFromKubernetesManifest([]byte(yaml))

		Expect(err).ToNot(HaveOccurred())

		result := strings.TrimSpace(string(bytes))
		Expect(result).To(Equal("image: " + digestRef))
	})
})
