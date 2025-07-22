terraform {
  required_providers {
    digitalocean = {
      source = "opentofu/digitalocean"
      version = "2.60.0"
    }

    github = {
      source  = "opentofu/github"
      version = "6.3.0"
    }

    cloudflare = {
      source  = "opentofu/cloudflare"
      version = "5.7.1"
    }
  }
}

variable "github_token" {
  description = "GitHub personal access token"
  type        = string
}

variable "cloudflare_token" {
  description = "Cloudflare API token"
  type        = string
}

variable "cloudflare_zone_id" {
  description = "Cloudflare zone ID"
  type        = string
}

variable "digitalocean_token" {
  description = "DigitalOcean API token"
  type        = string
}

provider "github" {
  token = var.github_token
}

provider "cloudflare" {
  api_token = var.cloudflare_token
}

provider "digitalocean" {
  token = var.digitalocean_token
}
