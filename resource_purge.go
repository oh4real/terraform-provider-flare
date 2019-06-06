package main

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/imroc/req"
)

func resourceServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceServerCreate,
		Read:   resourceServerRead,
		Update: resourceServerUpdate,
		Delete: resourceServerDelete,

		Schema: map[string]*schema.Schema{
			"address": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"zone_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"email": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"token": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"timestamp": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceServerCreate(d *schema.ResourceData, m interface{}) error {
	address := d.Get("address").(string)
	d.SetId(address)

	zone_id := d.Get("zone_id").(string)
	email := d.Get("email").(string)
	token := d.Get("token").(string)

	purgeRequest(address, email, token, zone_id)

	return resourceServerRead(d, m)
}

func resourceServerRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceServerUpdate(d *schema.ResourceData, m interface{}) error {
	// Enable partial state mode
	d.Partial(true)

	address := d.Get("address").(string)
	zone_id := d.Get("zone_id").(string)
	email := d.Get("email").(string)
	token := d.Get("token").(string)

	d.SetPartial("timestamp")

	purgeRequest(address, email, token, zone_id)

	d.Partial(false)

	return resourceServerRead(d, m)
}

func purgeRequest(address string, email string, token string, zone_id string) error {

	// curl -X POST \
	//   https://api.cloudflare.com/client/v4/zones/ZONE_ID/purge_cache \
	//   -H 'Content-Type: application/json' \
	//   -H 'X-Auth-Email: API_PASSWORD' \
	//   -H 'X-Auth-Key: API_TOKEN' \
	//   -d '{"hosts":["CUSTOM_HOSTNAME"]}'

	authHeader := req.Header{
		"Content-Type": "application/json",
		"X-Auth-Email": email,
		"X-Auth-Key":   token,
	}

	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/purge_cache", zone_id)

	body := fmt.Sprintf(`{"hosts":["%s"]}`, address)

	req.Post(url, body, authHeader)

	return nil
}

func resourceServerDelete(d *schema.ResourceData, m interface{}) error {
	// d.SetId("") is automatically called assuming delete returns no errors, but
	// it is added here for explicitness.
	d.SetId("")
	return nil
}
