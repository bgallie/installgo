/*
Copyright © 2024 Billy G. Allie <bill.allie@defiant.mug.org>
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	curVersion     string
	newVersion     string
	dlFileName     string
	dlFileCheckSum string
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "status will check for a newer version of GO.",
	Long: `status will check https://go.dev for the latest version of GO and optionally
install it if the --autoinstall option is given.`,
	Run: func(cmd *cobra.Command, args []string) {
		if !cmd.Flags().Lookup("autoupdate").Changed {
			autoupdate = igoViper.GetBool("autoupdate")
		}
		if !cmd.Flags().Lookup("maxcachetime").Changed {
			maxCacheTime = igoViper.GetFloat64("maxcachetime")
		}
		getCurrentVersion()
		fmt.Printf("Current Version: %s\n", curVersion)
		if err := scrapeLatestVersion(); err != nil {
			if err = scrapeLatestVersion(); err != nil {
				log.Fatal("Something went terribly wrong checking for the latest version:", err)
			}
		}
		if curVersion != newVersion {
			fmt.Printf("A new version, %s, is available.\n", newVersion)
			if autoupdate {
				updateGo()
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
	statusCmd.Flags().BoolVarP(&autoupdate, "autoupdate", "a", false, "install the latest version automatically.")
	statusCmd.Flags().Lookup("autoupdate").NoOptDefVal = "true"
	viper.BindPFlag("autoupdate", statusCmd.Flags().Lookup("autoupdate"))
	statusCmd.Flags().Float64VarP(&maxCacheTime, "maxcachetime", "m", 6.0, "time (in hours) that the cache is valid for.")
	statusCmd.Flags().Lookup("maxcachetime").NoOptDefVal = "0.0"
	viper.BindPFlag("maxcachetime", statusCmd.Flags().Lookup("maxcachetime"))
}
