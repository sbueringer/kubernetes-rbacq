// Copyright © 2017 Stefan Bueringer <sbueringer@gmail.com>
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
	"github.com/spf13/cobra"
	"github.com/sbueringer/kubernetes-rbacq/query"
)

var getCmd = &cobra.Command{
	Use:        "get [RESOURCE-TYPE]",
	Short:      "Displays one or many resources",
	Long:       `Displays one or many resources`,
	ValidArgs:  []string{"subject", "right"},
	ArgAliases: []string{"s", "subjects", "r", "rights"},
	Run: func(cmd *cobra.Command, args []string) {
		if args != nil && len(args) > 0 {
			switch
			args[0] {
			case "s": fallthrough
			case "subject": fallthrough
			case "subjects":
				query.GetSubjects(args)
			case "r": fallthrough
			case "right": fallthrough
			case "rights":
				query.GetRights(args)
			default:
				printUsage(cmd)
			}
		} else {
			printUsage(cmd)
		}
	},
}

func printUsage(cmd *cobra.Command) {
	cmd.Printf("You must specify the type of resource to get. Valid resource types are:\n")
	cmd.Println("\t* subjects (aka 'sub')")
	cmd.Println("\t* rights (aka 'r')")
}

func init() {
	RootCmd.AddCommand(getCmd)
	getCmd.PersistentFlags().StringVarP(&query.Output, "output", "o", "", "Set jsonpath e.g. with -o jsonpath='{.kind}:{.Name}'")
}
