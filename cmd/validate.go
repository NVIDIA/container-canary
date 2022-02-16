package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/jacobtomlinson/containercanairy/internal/checks"
	"github.com/jacobtomlinson/containercanairy/internal/platforms"
	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate a container against a platform",
	Long:  ``,
}

var validateKubeflowCmd = &cobra.Command{
	Use:   "kubeflow [IMAGE]",
	Short: "Validate a container against kubeflow",
	Long:  ``,
	Args:  imageArg,
	Run: func(cmd *cobra.Command, args []string) {
		if !platforms.ValidateKubeflow(args[0]) {
			fmt.Println("FAILED")
			os.Exit(1)
		}
		fmt.Println("PASSED")
	},
}

func imageArg(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return errors.New("requires an image argument")
	}

	if len(args) > 1 {
		return errors.New("too many arguments")
	}

	image := args[0]

	if checks.CheckImage(image, "docker") {
		return nil
	} else {
		return fmt.Errorf("no such image: %s", image)
	}
}

func init() {
	validateCmd.AddCommand(validateKubeflowCmd)

	rootCmd.AddCommand(validateCmd)
}
