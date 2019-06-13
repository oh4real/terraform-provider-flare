package flare

import (
	"fmt"

	"crypto/sha256"

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
			"host_name": &schema.Schema{
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
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
				StateFunc: hashSum,
			},
			"timestamp": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceServerCreate(d *schema.ResourceData, m interface{}) error {
	hostName := d.Get("host_name").(string)
	d.SetId(hostName)

	zoneID := d.Get("zone_id").(string)
	email := d.Get("email").(string)
	token := d.Get("token").(string)

	purgeRequest(hostName, email, token, zoneID)

	return resourceServerRead(d, m)
}

func resourceServerRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceServerUpdate(d *schema.ResourceData, m interface{}) error {
	// Enable partial state mode
	d.Partial(true)

	hostName := d.Get("host_name").(string)
	zoneID := d.Get("zone_id").(string)
	email := d.Get("email").(string)
	token := d.Get("token").(string)

	d.SetPartial("timestamp")

	purgeRequest(hostName, email, token, zoneID)

	d.Partial(false)

	return resourceServerRead(d, m)
}

func purgeRequest(hostName string, email string, token string, zoneID string) error {

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

	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/purge_cache", zoneID)

	body := fmt.Sprintf(`{"hosts":["%s"]}`, hostName)

	req.Post(url, body, authHeader)

	return nil
}

func resourceServerDelete(d *schema.ResourceData, m interface{}) error {
	// d.SetId("") is automatically called assuming delete returns no errors, but
	// it is added here for explicitness.
	d.SetId("")
	return nil
}

func hashSum(contents interface{}) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(contents.(string))))
}
