package cmd

import (
	nap "github.com/flynshuePersonal/napv2"
	"github.com/spf13/viper"
)

var api *nap.API

func GithubAPI() *nap.API {
	if api == nil {
		api = nap.NewAPI("https://api.github.com")
		auth := nap.AuthToken{Token: viper.GetString("token")}
		api.SetAuth(auth)
		api.AddResource("docs", DocsResource())
		api.AddResource("fork", ForkResource())
	}
	return api
}
