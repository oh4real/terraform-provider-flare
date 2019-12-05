package flare

import (
	"fmt"
	"log"

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
				ForceNew: true,
			},
			"zone_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"custom_origin_server": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ssl_method": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "http",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					methods := []string{"http", "cname", "email"}
					for _, method := range methods {
						if method == v {
							return
						}
					}
					errs = append(errs, fmt.Errorf("%q must be a valid SSL Method: http, cname, email. Got: %q", key, v))
					return
				},
			},
		},
	}
}

func resourceCustomHostnameCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*cloudflare.API)
	hostName := d.Get("host_name").(string)
	zoneID := d.Get("zone_id").(string)

	customHostName := cloudflare.CustomHostname{
		Hostname: hostName,
		SSL: cloudflare.CustomHostnameSSL{
			Method: d.Get("ssl_method").(string),
			Type:   "dv",
		},
		CustomOriginServer: d.Get("custom_origin_server").(string),
	}

	id, err := client.CustomHostnameIDByName(zoneID, hostName)
	if err != nil {
		// could legit error or err because not found, cloudflare-go treats general or 404s the same
		if msg := err.Error(); msg != "CustomHostname could not be found" {
			return err
		}

		log.Println("CustomHostnameIDByName err: ", err)
	}

	// custom hostname is already provisioned. set ID
	if len(id) > 0 {
		d.SetId(id)

		_, err := client.CustomHostname(zoneID, id)
		if err != nil {
			return err
		}

		// let's persist state directly from cloudflare
		return resourceCustomHostnameRead(d, m)
	}

	customHost, err := client.CreateCustomHostname(zoneID, customHostName)
	if err != nil {
		return err
	}

	log.Printf("%+v", customHost.Result)

	d.SetId(customHost.Result.ID)

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
	client := m.(*cloudflare.API)
	zoneID := d.Get("zone_id").(string)

	ssl := struct {
		SSL cloudflare.CustomHostnameSSL `json:"ssl,omitempty"`
	}{
		SSL: cloudflare.CustomHostnameSSL{
			Method: d.Get("ssl_method").(string),
			Type:   "dv",
		},
	}

	log.Printf("ssl: %+v", ssl)

	customHost, err := client.CustomHostname(zoneID, d.Id())
	if err != nil {
		return err
	}

	log.Printf("customHost: %+v", customHost)

	log.Println(fmt.Sprintf("/zones/%s/custom_hostnames/%s", zoneID, customHost.ID))

	// Until cloudflare-go client.UpdateCustomHostnameSSL() gets implemented, 
	// do this `manually` with client since it's got api key, etc
	raw, err := client.Raw("PATCH", fmt.Sprintf("/zones/%s/custom_hostnames/%s", zoneID, customHost.ID), ssl)
	if err != nil {
		return err
	}

	log.Printf("%+v", string(raw))

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
