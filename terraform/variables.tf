variable "neon_api_key" {
  description = "Neon API key from console.neon.tech"
  type        = string
  sensitive   = true
}

variable "neon_org_id" {
  description = "Neon organization ID from console.neon.tech organization settings (e.g. org-something-123456)"
  type        = string
}

variable "aws_region" {
  description = "AWS region for Lambda and API Gateway"
  type        = string
  default     = "us-west-2"
}

variable "app_name" {
  description = "Application name prefix for all resources"
  type        = string
  default     = "go-todo"
}

variable "cloudflare_api_token" {
  description = "Cloudflare API token with Zone:Read and DNS:Edit permissions for amohr.net"
  type        = string
  sensitive   = true
}

variable "domain_name" {
  description = "Full domain name for the app"
  type        = string
  default     = "gotodo.amohr.net"
}

variable "cloudflare_zone_name" {
  description = "Root domain managed in Cloudflare"
  type        = string
  default     = "amohr.net"
}
