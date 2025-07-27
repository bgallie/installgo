/*
Copyright © 2025 Billy G. Allie <bill.allie@defiant.mug.org>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
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

	return strings.Split(b.String(), igoViper.GetString("seperator"))
}
