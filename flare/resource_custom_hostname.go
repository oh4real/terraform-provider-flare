package flare

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/cloudflare/cloudflare-go"
	"github.com/hashicorp/terraform/helper/schema"
)

type ExtendedCustomHostname struct {
	cloudflare.CustomHostname
	CustomOriginServer string `json:"custom_origin_server,omitempty"`
}

func resourceCustomHostname() *schema.Resource {
	return &schema.Resource{
		Create: resourceCustomHostnameCreate,
		Read:   resourceCustomHostnameRead,
		Update: resourceCustomHostnameUpdate,
		Delete: resourceCustomHostnameDelete,

		Schema: map[string]*schema.Schema{
			"host_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"zone_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"custom_origin_server": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceCustomHostnameCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*cloudflare.API)

	customHostName := ExtendedCustomHostname{}
	customHostName.Hostname = d.Get("host_name").(string)
	customHostName.CustomOriginServer = d.Get("custom_origin_server").(string)
	customHostName.SSL = cloudflare.CustomHostnameSSL{
		Method: "http",
		Type:   "dv",
	}

	// Until cloudflare-go gets "CustomOriginServer", do this `manually`
	raw, err := client.Raw("POST", fmt.Sprintf("/zones/%s/custom_hostnames", d.Get("zone_id").(string)), customHostName)
	if err != nil {
		return err
	}

	log.Printf("%+v", string(raw))

	resp := ExtendedCustomHostname{}

	json.Unmarshal(raw, &resp)

	log.Printf("%+v", resp)

	d.SetId(resp.ID)

	return resourceCustomHostnameRead(d, m)
}

func resourceCustomHostnameRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*cloudflare.API)
	zoneID := d.Get("zone_id").(string)

	_, err := client.CustomHostname(zoneID, d.Id())
	if err != nil {
		return err
	}

	return nil
}

func resourceCustomHostnameUpdate(d *schema.ResourceData, m interface{}) error {
	// only changes allowed are SSL
	// not implemented, simple SSL http/dv for now
	return nil
}

func resourceCustomHostnameDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*cloudflare.API)
	zoneID := d.Get("zone_id").(string)

	err := client.DeleteCustomHostname(zoneID, d.Id())
	if err != nil {
		return err
	}

	// d.SetId("") is automatically called assuming delete returns no errors, but
	// it is added here for explicitness.
	d.SetId("")

	return nil
}
