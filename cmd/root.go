// Copyright © 2018 Stefan Bueringer
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

	"github.com/spf13/cobra"
	"github.com/sbueringer/kubernetes-rbacq/query"
)

var RootCmd = &cobra.Command{
	Use:   "rbacq",
	Short: "rbacq simplifies querying the Kubernetes RBAC API",
	Long:  `rbacq simplifies querying the Kubernetes RBAC API`,
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(query.InitKubeCfg)

	RootCmd.PersistentFlags().StringVarP(&query.KubeCfgFile, "kubeconfig", "k", fmt.Sprintf("%s%s.kube%sconfig",os.Getenv("HOME"),string(os.PathSeparator), string(os.PathSeparator)), "Path to the kubeconfig file to use for CLI requests")
	RootCmd.PersistentFlags().StringVarP(&query.Namespace, "namespace", "n", "default", "Specifies the Namespace in which to query")
	RootCmd.PersistentFlags().BoolVarP(&query.AllNamespaces, "all-namespaces", "a", false, "Specifies that all Namespaces should be queried (default \"false\")")
	RootCmd.PersistentFlags().BoolVarP(&query.System, "system", "s", false, "Show also System Objects (default \"false\")")
	RootCmd.PersistentFlags().BoolVarP(&query.ClusterWide, "cluster-wide", "c", false, "Search cluster-wide (which includes ClusterRoles & ClusterRolebindings)")
}
