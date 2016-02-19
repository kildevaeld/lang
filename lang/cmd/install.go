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
	"errors"
	"fmt"
	"strings"

	"github.com/cheggaaa/pb"
	"github.com/kildevaeld/go-ascii2"
	"github.com/kildevaeld/lang"
	"github.com/spf13/cobra"
)

var versionFlag string
var binaryFlag bool

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "",
	Long:  ``,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("Usage: lang install <language>...")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here
		var ver string
		for _, l := range args {
			ver = versionFlag
			if index := strings.Index(l, "@"); index > -1 {
				split := strings.Split(l, "@")
				ver = split[1]
				l = split[0]
			}

			if ver == "" {
				ver = "latest"
			}

			if err := install(l, ver); err != nil {
				printError(err)
			}
		}

	},
}

func init() {
	RootCmd.AddCommand(installCmd)
	installCmd.Aliases = []string{"i"}

	installCmd.Flags().StringVarP(&versionFlag, "version", "v", "", "")
	installCmd.Flags().BoolVarP(&binaryFlag, "binary", "b", true, "")

}

func install(l, version string) error {
	var currentStep lang.Step
	var bar *pb.ProgressBar
	var process *Process
	fmt.Printf("Installing %s@%s\n", l, version)
	err := service.Install(l, version, binaryFlag, func(step lang.Step, progress, total int64) {

		if currentStep != step {
			if bar != nil {
				bar.NotPrint = true
				bar.Finish()
				fmt.Printf(ascii2.EraseLine)
				bar = nil
			}

			if process != nil {
				process.Done("")
				process = nil
			}

			if total > 0 {
				bar = pb.New64(total).Prefix("  " + stepToMsg(step) + "\t\t")
				bar.SetWidth(40)
				bar.ShowCounters = false
				//fmt.Printf("%s\n", step)
				//bar.NotPrint = true
				bar.Start()
				currentStep = step

			} else {
				process := &Process{Msg: stepToMsg(step) + "\t\t"}
				process.Start()
			}

		}
		if bar != nil {
			bar.Set64(progress)
		}

	})

	if bar != nil {
		bar.NotPrint = true
		bar.Finish()

		fmt.Printf(ascii2.EraseLines(2) + ascii2.EraseLine + fmt.Sprintf("  %s installed", l))
	}

	if process != nil {
		process.Done("\n")
	}
	//fmt.Printf(ascii2.EraseLine + ascii2.CursorUp(1) + ascii2.EraseLine)

	if err != nil {
		fmt.Printf("Could not install %s@%s: \n  %s\n", l, version, err.Error())
	} else {
		fmt.Printf("  %s@%s installed!\n\n", l, version)
	}

	return err
}

func stepToMsg(step lang.Step) string {
	switch step {
	case lang.Download:
		return "Dowloading"
	case lang.Install:
		return "Installing"
	case lang.Unpack:
		return "Unpacking"
	default:
		return step.String()
	}
}
