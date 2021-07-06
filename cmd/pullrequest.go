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

var (
	title string
	src   string
	dst   string
	msg   string
)

type PullRequestPayload struct {
	Title   string `json:"title"`
	Src     string `json:"head"`
	Dst     string `json:"base"`
	Message string `json:"body"`
	Modify  bool   `json:"maintainer_can_modify"`
}

type PullRequestResponse struct {
	URL string `json:"html_url"`
}

// pullrequestCmd represents the pullrequest command
var pullrequestCmd = &cobra.Command{
	Use:   "pullrequest",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := PullRequest(); err != nil {
			log.Fatalln("error creating pull request:", err)
		}
	},
}

func PullRequest() error {
	dstValues := strings.Split(dst, ":")
	if len(dstValues) != 2 {
		return fmt.Errorf("must supply dst branch in owner/project:branch format")
	}
	repo := strings.Split(dstValues[0], "/")
	if len(repo) != 2 {
		return fmt.Errorf("must supply dst branch in owner/project:branch format")
	}
	params := map[string]string{"owner": repo[0], "project": repo[1]}
	body := &PullRequestPayload{
		Title:   title,
		Src:     src,
		Dst:     dstValues[1],
		Message: msg,
		Modify:  true,
	}
	return GithubAPI().Call("pullrequest", params, body)
}

func PullRequestSuccess(resp *http.Response) error {
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	response := &PullRequestResponse{}
	if err := json.Unmarshal(b, response); err != nil {
		return err
	}
	fmt.Printf("create pull request %s\n", response.URL)
	return nil
}

func PullRequestFailed(resp *http.Response) error {
	return fmt.Errorf(`%s from %s
Verify that pull request for branch doesn't already exist and that the src/head branch exists on the remote with pending changes`,
		resp.Status, resp.Request.URL.Path)
}

func PullRequestResource() nap.RestResource {
	router := nap.NewRouter()
	router.RegisterFunc(201, PullRequestSuccess)
	router.RegisterFunc(422, PullRequestFailed)
	pullrequest := nap.NewResource("POST", "/repos/{{.owner}}/{{.project}}/pulls", router)
	return pullrequest
}

func init() {
	rootCmd.AddCommand(pullrequestCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pullrequestCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pullrequestCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	pullrequestCmd.PersistentFlags().StringVar(&title, "title", "Basic Pull Request", "Title for the pull request")
	pullrequestCmd.PersistentFlags().StringVar(&src, "src", "", "source branch for the pull request")
	pullrequestCmd.PersistentFlags().StringVar(&dst, "dst", "", "destination branch where you want the changes pulled into")
	pullrequestCmd.PersistentFlags().StringVar(&msg, "msg", "", "Message body contents for the pull request")
}
