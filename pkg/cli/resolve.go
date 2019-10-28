package cli

import (
	"fmt"
	"io/ioutil"

	"github.com/projectriff/k8s-manifest-scanner/pkg/scan"
	"github.com/spf13/cobra"
)

type resolveCmd struct {
	file string
	dest string
}

func NewResolveCommand() *cobra.Command {
	res := &resolveCmd{}

	cmd := &cobra.Command{
		Use:   "resolve",
		Short: "resolves the tags for images in a kubernetes resource file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			res.file = args[0]
			return res.run()
		},
	}

	f := cmd.Flags()
	f.StringVarP(&res.dest, "output-file", "o", "", "File for output")
	return cmd
}

func (sc *resolveCmd) run() error {
	result, err := scan.ResolveImagesFromKubernetesManifest(sc.file, "")
	if err != nil {
		return err
	}

	if sc.dest != "" {
		return ioutil.WriteFile(sc.dest, result, 0644)
	}

	fmt.Println(string(result))

	return nil
}
