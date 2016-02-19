// Copyright © 2016 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"

	"github.com/kildevaeld/lang"
	"github.com/spf13/cobra"
)

var maxPrintFlag int32
var installedFlag bool

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) > 0 {
			la := service.GetLanguage(args[0])
			if la == nil {
				printError(fmt.Errorf("version :%s", args[0]))
			}

			def := la.Definition()
			var found lang.StrSlice

			for _, s := range def.Stable {
				if int32(len(found)) == maxPrintFlag && maxPrintFlag != 0 {
					break
				}
				if !found.Contains(s.Version) {
					found = append(found, s.Version)
				}

			}

			fmt.Printf("%s", found.Join(" "))
		} else {
			for _, la := range service.Languages() {
				fmt.Printf("%s ", la)
			}
		}

		fmt.Println()
	},
}

func init() {
	RootCmd.AddCommand(listCmd)

	listCmd.Flags().Int32VarP(&maxPrintFlag, "max", "m", 0, "")
	listCmd.Flags().BoolVarP(&installedFlag, "installed", "i", false, "Print installed versions")

}
