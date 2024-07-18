/*
Copyright © 2024 Billy G. Allie <bill.allie@defiant.mug.org>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
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
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Println("status called")
		getCurrentVersion()
		fmt.Printf("Current Version: %s\n", curVersion)
		scrapeLatestVersion()
		if curVersion != newVersion {
			fmt.Printf("A new version, %s, is available.\n", newVersion)
		}
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// statusCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// statusCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func isCacheValid(cacheFile string) bool {
	fCache, err := os.Open(cacheFile)
	if err != nil {
		log.Fatal(err)
	}
	defer fCache.Close()
	cacheInfo, err := fCache.Stat()
	if err != nil {
		log.Fatal(err)
	}
	return time.Since(cacheInfo.ModTime()).Hours() <= maxCacheTime
}
