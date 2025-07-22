resource "digitalocean_droplet" "web" {
  name     = "kamal-web"
  image    = "docker-20-04"
  region   = "nyc1"
  size     = "s-1vcpu-1gb"
  ssh_keys = [46054806]
}

resource "digitalocean_firewall" "web" {
  name = "kamal-web-firewall"

  droplet_ids = [digitalocean_droplet.web.id]

  inbound_rule {
    protocol         = "tcp"
    port_range       = "22"
    source_addresses = ["0.0.0.0/0", "::/0"]
  }

  inbound_rule {
    protocol         = "tcp"
    port_range       = "80"
    source_addresses = ["0.0.0.0/0", "::/0"]
  }

  outbound_rule {
    protocol              = "tcp"
    port_range            = "1-65535"
    destination_addresses = ["0.0.0.0/0", "::/0"]
  }

  outbound_rule {
    protocol              = "udp"
    port_range            = "1-65535"
    destination_addresses = ["0.0.0.0/0", "::/0"]
  }

  outbound_rule {
    protocol              = "icmp"
    destination_addresses = ["0.0.0.0/0", "::/0"]
  }
}

resource "github_actions_secret" "ip_address" {
  repository      = "knowhere"
  secret_name     = "DROPLET_IP_ADDRESS"
  plaintext_value = digitalocean_droplet.web.ipv4_address
}

resource "cloudflare_dns_record" "kamal_web" {
  zone_id = var.cloudflare_zone_id
  name    = "api"
  content = digitalocean_droplet.web.ipv4_address
  type    = "A"
  ttl     = 3600
  proxied = true
}

output "droplet_ip_address" {
  value = digitalocean_droplet.web.ipv4_address
}
