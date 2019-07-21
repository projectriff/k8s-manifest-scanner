package cli

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"

	"github.com/ghodss/yaml"
	"github.com/projectriff/k8s-manifest-scanner/pkg/scan"

	"github.com/projectriff/cnab-k8s-installer-base/pkg/apis/kab/v1alpha1"
	"github.com/spf13/cobra"
)

type scanCmd struct {
	file        string
	kabManifest bool
	byteContent []byte
	dest        string
}

func NewScanCommand() *cobra.Command {
	sc := &scanCmd{}

	cmd := &cobra.Command{
		Use:   "scan",
		Short: "scans a kubernetes resource file for images",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			sc.file = args[1]
			return sc.run()
		},
	}

	f := cmd.Flags()
	f.StringVarP(&sc.dest, "output-file", "o", "", "File for output")
	f.BoolVarP(&sc.kabManifest, "kab-manifest", "k", false, "Input file is a kab manifest")
	return cmd
}

func (sc *scanCmd) run() error {
	var images []string

	if !sc.kabManifest {
		var err error
		images, err = sc.scanKubernetesResourceFileForImages(sc.file)
		if err != nil {
			return err
		}
	} else {
		_, err := sc.isKabManifest()
		if err != nil {
			return err
		}
		images, err = sc.scanKabManifestForImages()
		if err != nil {
			return err
		}
	}

	jsonImages, err := json.MarshalIndent(images, "", "    ")
	if err != nil {
		os.Exit(1)
	}
	if sc.dest != "" {
		return ioutil.WriteFile(sc.dest, jsonImages, 0644)
	}
	fmt.Println(string(jsonImages))

	return nil
}

func (sc *scanCmd) readFile() error {
	if sc.byteContent != nil {
		return nil
	}

	content, err := ioutil.ReadFile(sc.file)
	if err != nil {
		return err
	}
	sc.byteContent = content
	return nil
}

func (sc *scanCmd) isKabManifest() (bool, error) {
	err := sc.readFile()
	if err != nil {
		return false, err
	}

	kabMfst := &v1alpha1.Manifest{}
	err = yaml.Unmarshal(sc.byteContent, kabMfst)
	if err != nil {
		return false, nil
	}
	return true, nil
}

func (sc *scanCmd) scanKabManifestForImages() ([]string, error) {
	err := sc.readFile()
	if err != nil {
		return nil, err
	}

	kabMfst := &v1alpha1.Manifest{}
	err = yaml.Unmarshal(sc.byteContent, kabMfst)
	if err != nil {
		panic("previously un-marshalled successfully. should not happen")
	}

	images := map[string]struct{}{}
	err = kabMfst.VisitResources(func(res v1alpha1.KabResource) error {
		fmt.Fprintf(os.Stderr, "Scanning %s\n", res.Path)
		imgs, err := scan.ListImagesFromKubernetesManifest(res.Path, "")
		if err != nil {
			return err
		}
		for _, i := range imgs {
			images[i] = struct{}{}
		}
		return nil
	})
	return keys(images), err
}

func (sc *scanCmd) scanKubernetesResourceFileForImages(file string) ([]string, error) {
	images := map[string]struct{}{}
	imgs, err := scan.ListImagesFromKubernetesManifest(file, "")
	if err != nil {
		return []string{}, err
	}
	for _, i := range imgs {
		images[i] = struct{}{}
	}
	return keys(images), err
}

func keys(m map[string]struct{}) []string {
	ks := []string{}
	for k, _ := range m {
		ks = append(ks, k)
	}
	sort.Strings(ks)
	return ks
}
