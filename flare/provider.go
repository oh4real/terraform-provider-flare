package flare

import (
	"github.com/cloudflare/cloudflare-go"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"flare_purge":           resourcePurge(),
			"flare_custom_hostname": resourceCustomHostname(),
		},

		Schema: map[string]*schema.Schema{
			"email": &schema.Schema{
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
			},
			"api_token": &schema.Schema{
				Type:      schema.TypeString,
				Required:  false,
				Optional:  true,
				Sensitive: true,
			},
			"api_key": &schema.Schema{
				Type:      schema.TypeString,
				Required:  false,
				Sensitive: true,
				Optional:  true,
			},
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		Email:    d.Get("email").(string),
		APIToken: d.Get("api_token").(string),
		APIKey:   d.Get("api_key").(string),
		Options:  []cloudflare.Option{},
	}

	client, err := config.Client()
	if err != nil {
		return nil, err
	}

	return client, err
}
