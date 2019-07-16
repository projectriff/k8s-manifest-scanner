package main

import (
	"fmt"
	"os"

	"github.com/projectriff/k8s-manifest-scanner/pkg/cli"
)

func main() {

	cmd := cli.NewScanCommand()
	err := cmd.Execute()
	if err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(1)
	}
}