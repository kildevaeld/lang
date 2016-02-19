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
	"time"

	"github.com/cheggaaa/pb"
	"github.com/kildevaeld/go-ascii2"
	"github.com/kildevaeld/lang"
	"github.com/spf13/cobra"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here
		fmt.Println("test called")

		var currentStep lang.Step
		var bar *pb.ProgressBar
		var process *Process

		testSteps(func(step lang.Step, progress, total int64) {
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
					process := &Process{Msg: "  " + stepToMsg(step) + "\t\t"}
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
			fmt.Printf(ascii2.EraseLines(2) + ascii2.EraseLine + "\r  installed")
		}
		/*if bar != nil {
			bar.NotPrint = true
			bar.FinishPrint(ascii2.EraseLine + ascii2.CursorUp(1) + ascii2.EraseLine + "  " + stepToMsg(currentStep) + "\t\tDone")

		}
		if process != nil {
			process.Done("\n")
		}*/
	},
}

func init() {
	RootCmd.AddCommand(testCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// testCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// testCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

func testSteps(fn func(step lang.Step, p, t int64)) error {
	steps := []lang.Step{lang.Compile, lang.Download, lang.Install}
	for _, s := range steps {
		i := int64(0)
		for i < 100 {
			fn(s, i, int64(100))
			time.Sleep(20 * time.Millisecond)
			i++
		}
	}
	i := 0
	for i < 100 {
		fn(lang.Compile, 0, 0)
		time.Sleep(20 * time.Millisecond)
		i++
	}
	return nil
}
