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
			if reinstall {
				fmt.Printf("Re-installing the current version (%s).\n", curVersion)
			} else {
				fmt.Printf("The latest version (%s) is already installed.\n", curVersion)
				return
			}
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
	updateCmd.Flags().BoolVarP(&reinstall, "reinstall", "r", false, "reinstall the latest version")
	updateCmd.Flags().Lookup("reinstall").NoOptDefVal = "true"
	updateCmd.Flags().BoolVarP(&autoupdate, "autoupdate", "a", false, "install the latest version without asking")
	updateCmd.Flags().Lookup("autoupdate").NoOptDefVal = "true"
}
