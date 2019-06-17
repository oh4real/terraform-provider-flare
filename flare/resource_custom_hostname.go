package flare

import (
	"fmt"

	"github.com/cloudflare/cloudflare-go"
	"github.com/hashicorp/terraform/helper/schema"
)



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
	zoneID := d.Get("zone_id").(string)
	hostName := d.Get("host_name").(string)
	customOriginServer := d.Get("custom_origin_server").(string)
	type ExtendedCustomHostname struct {
		cloudflare.CustomHostname
		CustomOriginServer string `json:"custom_origin_server,omitempty"`
	}
	customHostName := ExtendedCustomHostname{}
	customHostName.Hostname = hostName
	customHostName.CustomOriginServer = customOriginServer

	resp, err := client.CreateCustomHostname(zoneID, customHostName)
	if err != nil {
		return err
	} else if !resp.Success {
		return fmt.Errorf("CloudFlare CreateCustomHostname failed. Errors: %+v", resp.Errors)
	}

	d.SetId(resp.Result.ID)

	return resourceCustomHostnameRead(d, m)
}

func resourceCustomHostnameRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*cloudflare.API)
	zoneID := d.Get("zone_id").(string)

	resp, err := client.CustomHostname(zoneID, d.Id())

	return nil
}

func resourceCustomHostnameUpdate(d *schema.ResourceData, m interface{}) error {
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
