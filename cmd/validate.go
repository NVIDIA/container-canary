package cmd

import (
	"errors"
	"fmt"

	"github.com/jacobtomlinson/containercanairy/internal/checks"
	"github.com/jacobtomlinson/containercanairy/internal/config"
	"github.com/jacobtomlinson/containercanairy/internal/validator"
	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate [IMAGE]",
	Short: "Validate a container against a platform",
	Long:  ``,
	Args:  imageArg,
	RunE: func(cmd *cobra.Command, args []string) error {
		file, _ := cmd.Flags().GetString("file")
		// TODO Support loading from Url
		// TODO Check file exists
		validatorConfig, _ := config.LoadValidatorFromFile(file)
		image := args[0]
		cmd.Printf("Validating %s against %s\n", image, validatorConfig.Name)
		v, err := validator.Validate(image, validatorConfig)
		if err != nil {
			cmd.Printf("Error: %s\n", err.Error())
			return fmt.Errorf("ERRORED")
		}
		if !v {
			return fmt.Errorf("FAILED")
		}
		cmd.Println("PASSED")
		return nil
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
	rootCmd.AddCommand(validateCmd)
	validateCmd.PersistentFlags().String("file", "", "Path or URL of a manifest to validate against.")

}