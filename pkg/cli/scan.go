package cli

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/projectriff/k8s-manifest-scanner/pkg/scan"
	"github.com/spf13/cobra"
)

type scanCmd struct {
	file string
	dest string
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
	return cmd
}

func (sc *scanCmd) run() error {
	var images []string

	var err error
	images, err = scan.ListSortedImagesFromKubernetesManifest(sc.file, "")
	if err != nil {
		return err
	}

	jsonImages, err := json.MarshalIndent(images, "", "    ")
	if err != nil {
		os.Exit(1)
	}
	if sc.dest != "" {
		return ioutil.WriteFile(sc.dest, jsonImages, 0644)
	} else {
		fmt.Println(string(jsonImages))
	}

	return nil
}
