package flare

import (
	"crypto/sha256"
	"fmt"
	"log"
	"strings"

	"github.com/cloudflare/cloudflare-go"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceServerCreate,
		Read:   resourceServerRead,
		Update: resourceServerUpdate,
		Delete: resourceServerDelete,

		Schema: map[string]*schema.Schema{
			"host_names": &schema.Schema{
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
	client := m.(*cloudflare.API)

	zoneID := d.Get("zone_id").(string)
	hostNames := d.Get("host_names").(string)

	d.SetId(uuid.New().String())

	purgeCacheRequest(client, zoneID, hostNames)

	return resourceServerRead(d, m)
}

func resourceServerRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceServerUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*cloudflare.API)

	zoneID := d.Get("zone_id").(string)
	hostNames := d.Get("host_names").(string)

	// Enable partial state mode
	d.Partial(true)

	d.SetPartial("timestamp")

	err := purgeCacheRequest(client, zoneID, hostNames)
	if err != nil {
		return err
	}

	d.Partial(false)

	return resourceServerRead(d, m)
}

func resourceServerDelete(d *schema.ResourceData, m interface{}) error {
	// d.SetId("") is automatically called assuming delete returns no errors, but
	// it is added here for explicitness.
	d.SetId("")
	return nil
}

func purgeCacheRequest(client *cloudflare.API, zoneID string, hostNames string) error {

	hosts := strings.Split(hostNames, ",")
	pReq := cloudflare.PurgeCacheRequest{Hosts: hosts}

	log.Printf("%+v", pReq)

	resp, _ := client.PurgeCache(zoneID, pReq)

	log.Printf("%+v", resp)

	if !resp.Success {
		return fmt.Errorf("CloudFlare PurgeCache failed. Errors: %+v", resp.Errors)
	}

	return nil
}

func hashSum(contents interface{}) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(contents.(string))))
}
