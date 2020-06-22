package main

import (
	"github.com/hashicorp/terraform/plugin"
	provider "github.com/milamice62/terraplugin/resources"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: provider.Provider,
	})
}
