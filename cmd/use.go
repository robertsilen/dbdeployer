// DBDeployer - The MySQL Sandbox
// Copyright © 2006-2020 Giuseppe Maxia
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
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/spf13/cobra"

	"github.com/datacharmer/dbdeployer/common"
	"github.com/datacharmer/dbdeployer/globals"
)

func runInteractiveCmd(s string) error {
	// #nosec G204
	cmd := exec.Command(s)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func useSandbox(cmd *cobra.Command, args []string) error {
	sandboxHome, _ := cmd.Flags().GetString(globals.SandboxHomeLabel)
	sandbox := ""
	executable := ""
	sandboxList, err := common.GetSandboxesByDate(sandboxHome)
	if len(args) > 0 {
		sandbox = args[0]
	} else {
		if err != nil {
			return err
		}
		if len(sandboxList) == 0 {
			return fmt.Errorf("nothing to use. No sandboxes were found")
		}
		sandbox = sandboxList[len(sandboxList)-1].SandboxName
	}
	if len(args) > 1 {
		executable = args[1]
	} else {
		executable = "use"
	}
	for _, sb := range sandboxList {
		if sb.SandboxName == sandbox {
			sandboxDir := path.Join(sandboxHome, sandbox)
			fmt.Printf("running %s/ %s\n", sandboxDir, executable)
			useSingle := path.Join(sandboxDir, executable)
			startSingle := path.Join(sandboxDir, "start")
			useMulti := path.Join(sandboxDir, "n1")
			startMulti := path.Join(sandboxDir, "start_all")
			if common.ExecExists(useSingle) {
				fmt.Printf("%s\n", useSingle)
				out, _ := common.RunCmdCtrl(startSingle, true)
				if !strings.Contains(out, "already") {
					// The server was not already started
					fmt.Printf("%s\n", out)
				}
				return runInteractiveCmd(useSingle)
			} else if common.ExecExists(useMulti) {
				out, _ := common.RunCmdCtrl(startMulti, true)
				if !strings.Contains(out, "already") {
					// The server was not already started
					fmt.Printf("%s\n", out)
				}
				return runInteractiveCmd(useMulti)
			} else {
				return fmt.Errorf("no executable found for %s", sandbox)
			}
		}
	}
	return fmt.Errorf("sandbox %s not found", sandbox)
}

var useCmd = &cobra.Command{
	Use:   "use [sandbox_name [executable]]",
	Short: "uses a sandbox",
	Long: `Uses a given sandbox.
If a sandbox is indicated, it will be used.
Otherwise, it will use the latest deployed sandbox.
Optionally, an executable can be set as second argument.`,
	Example: `
$ dbdeployer use                    # runs "use" on the latest deployed sandbox
$ dbdeployer use rsandbox_8_0_22    # runs "m" on replication sandbox rsandbox_8_0_22
$ dbdeployer use rsandbox_8_0_22 s1 # runs "s1" on replication sandbox rsandbox_8_0_22
$ echo 'SELECT @@SERVER_ID' | dbdeployer use # pipes an SQL query to latest deployed sandbox
`,
	RunE: useSandbox,
}

func init() {
	rootCmd.AddCommand(useCmd)
}
