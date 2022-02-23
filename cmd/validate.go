package cmd

import (
	"errors"
	"fmt"
	"strings"

	canaryv1 "github.com/jacobtomlinson/containercanary/internal/apis/v1"
	"github.com/jacobtomlinson/containercanary/internal/config"
	"github.com/jacobtomlinson/containercanary/internal/validator"
	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate [IMAGE]",
	Short: "Validate a container against a platform",
	Long:  ``,
	Args:  imageArg,
	RunE: func(cmd *cobra.Command, args []string) error {
		file, err := cmd.Flags().GetString("file")
		if err != nil {
			return err
		}
		var validatorConfig *canaryv1.Validator

		if strings.Contains(file, "://") {
			validatorConfig, err = config.LoadValidatorFromURL(file)
		} else {
			validatorConfig, err = config.LoadValidatorFromFile(file)
		}
		if err != nil {
			return err
		}

		image := args[0]
		cmd.Printf("Validating %s against %s\n", image, validatorConfig.Name)
		v, err := validator.Validate(image, validatorConfig)
		if err != nil {
			cmd.Printf("Error: %s\n", err.Error())
			cmd.Println("ERRORED")
			return nil
		}
		if !v {
			cmd.Println("FAILED")
			return nil
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

	if validator.CheckImage(image, "docker") {
		return nil
	} else {
		return fmt.Errorf("no such image: %s", image)
	}
}

func init() {
	rootCmd.AddCommand(validateCmd)
	validateCmd.PersistentFlags().String("file", "", "Path or URL of a manifest to validate against.")

}
