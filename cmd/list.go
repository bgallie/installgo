/*
Copyright © 2024 Billy G. Allie <bill.allie@defiant.mug.org>
*/
package cmd

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"
)

// getCmd represents the get command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list the contents of the config file.",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			keys := igoViper.AllKeys()
			if len(keys) == 0 {
				fmt.Println("No configuration keys found.")
				return
			}
			sort.Strings(keys)
			for _, a := range keys {
				val := igoViper.Get(a)
				if val == nil {
					fmt.Printf("%s: <nil>\n", a)
				} else {
					switch val := val.(type) {
					case string:
						fmt.Printf("%s = '%s'\n", a, val)
					case int:
						fmt.Printf("%s = %d\n", a, val)
					case bool:
						fmt.Printf("%s = %t\n", a, val)
					case []interface{}:
						fmt.Printf("%s = [", a)
						for i, v := range val {
							if i > 0 {
								fmt.Print(", ")
							}
							// Handle different types in the slice
							switch v := v.(type) {
							case string:
								fmt.Printf("'%s'", v)
							case int:
								fmt.Printf("%d", v)
							case bool:
								fmt.Printf("%t", v)
							default:
								// Default case for unexpected types
								// This can be adjusted based on expected types
								fmt.Printf("%v", v)
							}
						}
						fmt.Println("]")
					case map[string]interface{}:
						fmt.Printf("%s = {\n", a)
						for k, v := range val {
							fmt.Printf("  %s: %v\n", k, v)
						}
						fmt.Println("}")
					default:
						fmt.Printf("%s: %v\n", a, val)
					}
				}
			}
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
