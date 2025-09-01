/*
Copyright © 2025 Billy G. Allie <bill.allie@defiant.mug.org>
*/
package cmd

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

type templateArgs struct {
	Editor     string
	ConfigFile string
}

var commandArgs []string

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Start an editor to edit the installgo configuration file",
	Run: func(cmd *cobra.Command, args []string) {
		commandArgs = getCommandArgs()
		cmdLine := exec.Command(commandArgs[0], commandArgs[1:]...)
		cmdErr := cmdLine.Run()
		if cmdErr != nil {
			log.Fatal(cmdErr)
		}
	},
}

func init() {
	rootCmd.AddCommand(editCmd)
}

func getCommandArgs() []string {
	var err error
	parms := new(templateArgs)
	parms.Editor = igoViper.GetString(fmt.Sprintf("editor.%s.editor", osCpuType))
	parms.ConfigFile = igoViper.ConfigFileUsed()
	command := template.Must(template.New("command").Parse(igoViper.GetString(fmt.Sprintf("editor.%s.command", osCpuType))))
	var b strings.Builder
	err = command.Execute(&b, parms)
	if err != nil {
		log.Fatal(err)
	}

	return strings.Split(b.String(), igoViper.GetString("separator"))
}
