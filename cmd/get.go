/*
Copyright © 2024 Billy G. Allie <bill.allie@defiant.mug.org>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "get the value associated with the given keys from the config file.",
	Long: `get the values associated with the given keys from the config file.  If no
keys are given, display all the key/value pairs in the config file.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			for _, a := range viper.AllKeys() {
				fmt.Printf("%s = %s\n", a, viper.GetString(a))
			}
			return
		}
		for _, a := range args {
			fmt.Printf("%s = %s\n", a, viper.GetString(a))
		}
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}
