/*
Copyright © 2024 Billy G. Allie <bill.allie@defiant.mug.org>
*/
package cmd

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "update will install the latest version of GO if not already installed.",
	Long: `update will check https://go.dev for the latest version of GO and optionally
install it.  You can also have update reinstall the latest version if it is
already installed on your system.`,
	Run: func(cmd *cobra.Command, args []string) {
		if !cmd.Flags().Lookup("autoupdate").Changed {
			autoupdate = igoViper.GetBool("autoupdate")
		}
		if !cmd.Flags().Lookup("reinstall").Changed {
			reinstall = igoViper.GetBool("reinstall")
		}
		if !cmd.Flags().Lookup("maxcachetime").Changed {
			maxCacheTime = igoViper.GetFloat64("maxcachetime")
		}
		if curVersion == "" {
			getCurrentVersion()
			if err := scrapeLatestVersion(); err != nil {
				if err = scrapeLatestVersion(); err != nil {
					log.Fatal("Something went terribly wrong checking for the latest version:", err)
				}
			}
		}
		if curVersion == newVersion {
			if reinstall {
				fmt.Printf("Re-installing the latest version (%s).\n", curVersion)
			} else {
				fmt.Printf("The latest version (%s) is already installed.\n", curVersion)
				return
			}
		} else {
			reinstall = false // ignore reinstall if the current version is not the new version.
		}

		if !autoupdate {
			fmt.Print("Do you want to ")
			if reinstall {
				fmt.Print("re-")
			}
			fmt.Printf("install version %s [yN]? ", newVersion)
			inLine, err := bufio.NewReader(os.Stdin).ReadString('\n')
			if err != nil {
				if err != io.EOF {
					panic(err)
				}
				return
			}
			if strings.TrimLeft(strings.ToLower(inLine), " \t")[0] != 'y' {
				return
			}
		}
		updateGo()
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
	updateCmd.Flags().BoolVarP(&reinstall, "reinstall", "r", false, "reinstall the latest version if already installed.")
	updateCmd.Flags().Lookup("reinstall").NoOptDefVal = "true"
	igoViper.BindPFlag("reinstall", updateCmd.Flags().Lookup("reinstall"))
	updateCmd.Flags().BoolVarP(&autoupdate, "autoupdate", "a", false, "install the latest version without asking.")
	updateCmd.Flags().Lookup("autoupdate").NoOptDefVal = "true"
	igoViper.BindPFlag("autoupdate", updateCmd.Flags().Lookup("autoupdate"))
	updateCmd.Flags().Float64VarP(&maxCacheTime, "maxcachetime", "m", 0.0, "time (in hours) that the cache is valid for.")
	statusCmd.Flags().Lookup("maxcachetime").NoOptDefVal = "0.0"
	igoViper.BindPFlag("maxCacheTime", updateCmd.Flags().Lookup("maxcachetime"))
}
