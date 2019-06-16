package flare

import (
	"log"

	"github.com/cloudflare/cloudflare-go"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceServerCreate,
		Read:   resourceServerRead,
		Update: resourceServerUpdate,
		Delete: resourceServerDelete,

		Schema: map[string]*schema.Schema{
			"host_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"zone_id": &schema.Schema{
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
	// config := m.(Config)
	client := m.(*cloudflare.API)

	zoneID := d.Get("zone_id").(string)
	hostName := d.Get("host_name").(string)

	d.SetId(hostName)

	purgeCacheRequest(client, zoneID, hostName)

	return resourceServerRead(d, m)
}

func resourceServerRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceServerUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*cloudflare.API)

	zoneID := d.Get("zone_id").(string)
	hostName := d.Get("host_name").(string)

	// Enable partial state mode
	d.Partial(true)

	d.SetPartial("timestamp")

	purgeCacheRequest(client, zoneID, hostName)

	d.Partial(false)

	return resourceServerRead(d, m)
}

func resourceServerDelete(d *schema.ResourceData, m interface{}) error {
	// d.SetId("") is automatically called assuming delete returns no errors, but
	// it is added here for explicitness.
	d.SetId("")
	return nil
}

func purgeCacheRequest(client *cloudflare.API, zoneID string, hostName string) error {

	pReq := cloudflare.PurgeCacheRequest{Hosts: []string{hostName}}

	log.Printf("%+v", pReq)

	resp, _ := client.PurgeCache(zoneID, pReq)

	log.Printf("%+v", resp)

	return nil
}
