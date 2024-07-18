/*
Copyright © 2024 Billy G. Allie <bill.allie@defiant.mug.org>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display version information",
	Long:  `Display version and detailed build information for installgo.`,
	Run: func(cmd *cobra.Command, args []string) {
		if Version == "" {
			Version = "(devel)"
		}
		fmt.Println("    Version:", Version)
		if len(GitDate) > 1 {
			fmt.Println("Commit Date:", GitDate)
			fmt.Println("     Commit:", GitCommit)
			fmt.Println("      State:", GitState)
			fmt.Println("    Summary:", GitSummary)
		}
		if BuildDate != "not set" {
			fmt.Println(" Build Date:", BuildDate)
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
