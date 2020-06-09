package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/milamice62/terraplugin/provider"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: provider.Provider,
	})
}
