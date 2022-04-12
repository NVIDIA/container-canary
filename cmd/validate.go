/*
* SPDX-FileCopyrightText: Copyright (c) <2022> NVIDIA CORPORATION & AFFILIATES. All rights reserved.
* SPDX-License-Identifier: Apache-2.0
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
* http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */

package cmd

import (
	"errors"
	"fmt"

	"github.com/nvidia/container-canary/internal/validator"
	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:           "validate [IMAGE]",
	Short:         "Validate a container against a platform",
	Long:          ``,
	Args:          imageArg,
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		file, err := cmd.Flags().GetString("file")
		if err != nil {
			return err
		}
		if file == "" {
			return errors.New("you must specify a manifest with '--file path/url'")
		}

		image := args[0]
		debug, err := cmd.Flags().GetBool("debug")
		if err != nil {
			return err
		}

		v, err := validator.Validate(image, file, cmd, debug)
		if err != nil {
			cmd.Printf("Error: %s\n", err.Error())
			return err
		}
		if !v {
			return errors.New("validation failed")
		}
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

	if validator.CheckImage(cmd, image, "docker") {
		return nil
	} else {
		return fmt.Errorf("no such image: %s", image)
	}
}

func init() {
	rootCmd.AddCommand(validateCmd)
	validateCmd.PersistentFlags().String("file", "", "Path or URL of a manifest to validate against.")
	validateCmd.PersistentFlags().Bool("debug", false, "Keep container running on failure for debugging.")

}
