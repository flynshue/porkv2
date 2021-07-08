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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	nap "github.com/flynshuePersonal/napv2"
	"github.com/spf13/cobra"
)

type SearchResponse struct {
	Items []SearchResult `json:"items"`
}

type SearchResult struct {
	Name string `json:"full_name"`
}

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Fatalln("must supply search keywords")
		}
		if err := SearchByKeywords(args); err != nil {
			log.Fatalln("error search for repositories: ", err)
		}
	},
}

func SearchByKeywords(keywords []string) error {
	query := strings.Join(keywords, "+")
	params := map[string]string{"query": query}
	return GithubAPI().Call("search", params, nil)
}

func SearchSuccess(resp *http.Response) error {
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	response := &SearchResponse{}
	if err := json.Unmarshal(b, response); err != nil {
		return err
	}
	for _, results := range response.Items {
		fmt.Println(results.Name)
	}
	return nil
}

func SearchResource() nap.RestResource {
	router := nap.NewRouter()
	router.RegisterFunc(200, SearchSuccess)
	search := nap.NewResource("GET", "/search/repositories?q={{.query}}", router)
	return search
}

func init() {
	rootCmd.AddCommand(searchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// searchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// searchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
