// Copyright (c) 2022, NVIDIA CORPORATION.

package cmd

import (
	"fmt"
	"strings"

	"github.com/jacobtomlinson/containercanary/internal"
	"github.com/spf13/cobra"
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
	cmd.Printf(" %-16s %s\n", fmt.Sprintf("%s:", strings.Title(title)), value)
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
