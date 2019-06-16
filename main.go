package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/oh4real/terraform-provider-flare/flare"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: flare.Provider})
}
