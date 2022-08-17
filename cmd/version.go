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
	"fmt"

	"github.com/nvidia/container-canary/internal"
	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of containercanary",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Container Canary")
		showLine(cmd, "Version", internal.Version)
		showLine(cmd, "Go version", internal.GoVersion)
		showLine(cmd, "Commit", internal.Commit)
		showLine(cmd, "OS/Arch", fmt.Sprintf("%s/%s", internal.Os, internal.Arch))
		showLine(cmd, "Built", internal.Buildtime)
	},
}

func showLine(cmd *cobra.Command, title string, value string) {
	cmd.Printf(" %-16s %s\n", fmt.Sprintf("%s:", cases.Title(language.Und, cases.NoLower).String(title)), value)
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
