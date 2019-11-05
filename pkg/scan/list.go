package scan

import (
	"bytes"
	"io"
	"sort"

	"github.com/pivotal/go-ape/pkg/furl"
	"gopkg.in/yaml.v3"
)

func ListSortedImagesFromKubernetesManifest(res string, baseDir string) ([]string, error) {
	images, err := ListImagesFromKubernetesManifest(res, baseDir)
	if err != nil {
		return nil, err
	}
	sort.Strings(images)
	return images, nil
}

func ListImagesFromKubernetesManifest(res string, baseDir string) ([]string, error) {
	contents, err := furl.Read(res, baseDir)
	if err != nil {
		return nil, err
	}
	return ListImagesFromContent(contents)
}

func ListSortedImagesFromContent(contents []byte) ([]string, error) {
	images, err := ListImagesFromContent(contents)
	if err != nil {
		return nil, err
	}
	sort.Strings(images)
	return images, nil
}

func ListImagesFromContent(contents []byte) ([]string, error) {
	var images []string

	d := yaml.NewDecoder(bytes.NewReader(contents))

	for {
		var doc yaml.Node
		err := d.Decode(&doc)

		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		for _, node := range SearchImageNodes(&doc) {
			images = append(images, node.Value)
		}
	}

	return images, nil
}
