package cmd

import (
	"fmt"
	"strings"

	"github.com/jacobtomlinson/containercanairy/internal"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of containercanairy",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Container Canairy")
		showLine("Version", internal.Version)
		showLine("Go version", internal.GoVersion)
		showLine("Commit", internal.Commit)
		showLine("OS/Arch", fmt.Sprintf("%s/%s", internal.Os, internal.Arch))
		showLine("Built", internal.Buildtime)
	},
}

func showLine(title string, value string) {
	fmt.Printf(" %-16s %s\n", fmt.Sprintf("%s:", strings.Title(title)), value)
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
