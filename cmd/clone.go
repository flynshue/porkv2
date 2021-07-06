/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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
	"log"

	"github.com/spf13/cobra"
)

var (
	ref    string
	create bool
)

// cloneCmd represents the clone command
var cloneCmd = &cobra.Command{
	Use:   "clone",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Fatalln("must supply repository")
		}
		if err := CloneRepo(args[0]); err != nil {
			log.Fatalln("error cloning repository:", err)
		}
	},
}

func CloneRepo(repo string) error {
	ghRepo, err := NewGHRepo(repo)
	if err != nil {
		return err
	}
	if err := ghRepo.Clone(); err != nil {
		return err
	}
	log.Printf("successfully cloned repo %s into %s", repo, ghRepo.Dir)
	return ghRepo.Checkout(ref, create)
}

func init() {
	rootCmd.AddCommand(cloneCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cloneCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cloneCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	cloneCmd.PersistentFlags().StringVar(&ref, "ref", "main", "remote reference to checkout i.e branch")
	cloneCmd.PersistentFlags().BoolVar(&create, "create", false, "create remote reference if it doesn't exist")
}
