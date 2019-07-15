package cli

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/ghodss/yaml"
	"github.com/projectriff/ciu/pkg/scan"
	"github.com/projectriff/cnab-k8s-installer-base/pkg/apis/kab/v1alpha1"
	"github.com/spf13/cobra"
)

type scanCmd struct {
	file string
	byteContent []byte
	dest string
}

func NewScanCommand() *cobra.Command {
	sc := &scanCmd{}
	cmd := &cobra.Command {
		Use: "scan",
		Short: "scans urls for images",
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			sc.file = args[1]
			return sc.run()
		},
	}

	f := cmd.Flags()
	f.StringVarP(&sc.dest, "output-file", "o", "images.json", "Save images to file path")
	return cmd
}

func (sc *scanCmd) run() error {
	isKabMfst, err := sc.isKabManifest()
	if err != nil {
		return err
	}
	if isKabMfst {
		images, err := sc.scanKabManifestForImages()
		if err != nil {
			return err
		}
		mfstBytes, err := json.MarshalIndent(images, "", "    ")
		if err != nil {
			os.Exit(1)
		}
		err = ioutil.WriteFile(sc.dest, mfstBytes, 0644)
		if err != nil {
			return err
		}
	}

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

	images := []string{}
	err = kabMfst.VisitResources(func(res v1alpha1.KabResource) error {
		tmpImgs, err := scan.ListImages(res.Path, "")
		if err != nil {
			return err
		}
		images = append(images, tmpImgs...)
		return nil
	})
	return images, err
}
