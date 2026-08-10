/*
* SPDX-FileCopyrightText: Copyright (c) <2026> NVIDIA CORPORATION & AFFILIATES. All rights reserved.
* SPDX-License-Identifier: Apache-2.0
 */

package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/nvidia/container-canary/internal/schema"
	"github.com/spf13/cobra"
)

var schemaPath string

var schemaValidateCmd = &cobra.Command{
	Use:   "schema-validate manifest [manifest...]",
	Short: "Validate manifests against the published Validator schema",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		for _, manifest := range args {
			if err := schema.ValidateFile(schemaPath, manifest); err != nil {
				return err
			}
		}
		fmt.Fprintln(cmd.OutOrStdout(), "schema validation passed")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(schemaValidateCmd)
	schemaValidateCmd.Flags().StringVar(&schemaPath, "schema", filepath.FromSlash(schema.Path), "Path to the Validator JSON Schema")
}
