package scan

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/pivotal/go-ape/pkg/furl"
	"gopkg.in/yaml.v3"
)

func ResolveImagesFromKubernetesManifest(res, baseDir string) ([]byte, error) {
	contents, err := furl.Read(res, baseDir)
	if err != nil {
		return nil, err
	}

	return resolveImagesFromKubernetesManifest(contents)
}

func resolveImagesFromKubernetesManifest(contents []byte) ([]byte, error) {
	d := yaml.NewDecoder(bytes.NewReader(contents))

	var docNodes []*yaml.Node
	images := make(map[string]string)

	for {
		var doc yaml.Node
		err := d.Decode(&doc)

		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		docNodes = append(docNodes, &doc)

		for _, node := range SearchImageNodes(&doc) {
			val := strings.TrimSpace(node.Value)

			if digestRef, ok := images[val]; ok {
				node.Value = digestRef
				continue
			}

			ref, err := name.ParseReference(val)
			if err != nil {
				return nil, fmt.Errorf("parsing reference %q: %v", val, err)
			}

			desc, err := remote.Get(ref, remote.WithAuthFromKeychain(authn.DefaultKeychain))
			if err != nil {
				return nil, fmt.Errorf("error fetching resource %q: %v", ref, err)
			}

			digestRef, err := name.NewDigest(fmt.Sprintf("%s@%s", ref.Context(), desc.Digest))
			if err != nil {
				return nil, fmt.Errorf("error fetching digest: %v", err)
			}

			images[node.Value] = digestRef.String()
			node.Value = digestRef.String()
		}
	}

	buf := &bytes.Buffer{}
	e := yaml.NewEncoder(buf)

	for _, doc := range docNodes {
		err := e.Encode(doc)
		if err != nil {
			return nil, fmt.Errorf("failed to encode output: %v", err)
		}
	}

	return buf.Bytes(), nil
}
