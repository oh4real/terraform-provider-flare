# terraform_provider_flare
A thin Terraform Provider designed to purge cache for Cloudflare custom_hostnames or DNS entries.

Sample usage:

```
provider "flare" {
    api_token = var.api_token
}

resource "flare_custom_hostname" "your_custom_host" {
    host_name = var.custom_hostname
    custom_origin_server = var.dnsrecord
    zone_id = var.zone_id
    ssl_method = "http"
}

resource "flare_purge" "your_custom_host" {
    host_names = var.host_names
    timestamp = timestamp()
    zone_id = var.zone_id
}

```

### Using travis and vendor
1. Update `go.mod` for versions and requirements
2. Execute `go mod vendor -v` to get latest into `/vendor` for travis
3. Execute `go install -mod=vendor` to verify