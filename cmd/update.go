/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/cavaliergopher/grab/v3"
	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("update called")
		if curVersion == "" {
			getCurrentVersion()
			scrapeLatestVersion()
		}
		if curVersion == newVersion {
			fmt.Printf("The latest version (%s) is already installed.\n", curVersion)
			return
		}
		client := grab.NewClient()
		getFile := fmt.Sprintf("https://go.dev/dl/%s", dlFileName)
		client.UserAgent = "Mozilla/5.0"
		req, err := grab.NewRequest(os.TempDir(), getFile)
		if err != nil {
			panic(err)
		}

		resp := client.Do(req)
		if err := resp.Err(); err != nil {
			panic(err)
		}

		defer os.Remove(resp.Filename)
		sha256Chksum := calculateSHA256(resp.Filename)
		if dlFileCheckSum != sha256Chksum {
			log.Fatalf("File validation failed!\nOriginal checksum.: %s\nCalculate checksum: %s\n", dlFileCheckSum, sha256Chksum)
		}
		fmt.Printf("File validation successful.\nRemoving go version %s\n", curVersion)
		cmdToRun := fmt.Sprintf("rm -rf %s/go", installDir)
		cmdErr := exec.Command("sudo", strings.Split(cmdToRun, " ")...).Run()
		if cmdErr != nil {
			log.Fatal(cmdErr)
		}
		fmt.Printf("Installing version %s\n", newVersion)
		cmdToRun = fmt.Sprintf("tar -C %s -xf %s", installDir, resp.Filename)
		cmdErr = exec.Command("sudo", strings.Split(cmdToRun, " ")...).Run()
		if cmdErr != nil {
			log.Fatal(cmdErr)
		}
		getCurrentVersion()
		fmt.Printf("Installed version is now %s\n", curVersion)
		fmt.Println("Done")
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// updateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// updateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func calculateSHA256(fileName string) string {
	f, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf("%x", h.Sum(nil))

}
