resource "digitalocean_droplet" "web" {
  name      = "kamal-web"
  image     = "docker-20-04"
  region    = "nyc1"
  size      = "s-1vcpu-1gb"
  ssh_keys  = [46054806]
  user_data = <<-USER_DATA
    #!/bin/bash
    sudo ufw insert 1 allow from ${var.ip_allowlist} to any port 22
    sudo ufw insert 2 allow port 8080
    sudo ufw reload
  USER_DATA
}

output "droplet_ip_address" {
  value = digitalocean_droplet.web.*.ipv4_address
}