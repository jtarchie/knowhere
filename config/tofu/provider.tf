terraform {
  required_providers {
    digitalocean = {
      source  = "digitalocean/digitalocean"
      version = "~> 2.0"
    }
  }
}

variable "digitalocean_token" {
  description = "DigitalOcean API token"
  type        = string
}

variable "ip_allowlist" {
  description = "IP address to whitelist for SSH access"
  type        = string
}

provider "digitalocean" {
  token = var.digitalocean_token
}