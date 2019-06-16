package flare

import (
	"github.com/cloudflare/cloudflare-go"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"flare_purge": resourceServer(),
		},

		Schema: map[string]*schema.Schema{
			"email": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"token": &schema.Schema{
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	options := []cloudflare.Option{}
	config := Config{
		Email:   d.Get("email").(string),
		Token:   d.Get("token").(string),
		Options: options,
	}

	client, err := config.Client()
	if err != nil {
		return nil, err
	}

	return client, err
}
