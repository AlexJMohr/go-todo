data "cloudflare_zone" "main" {
  name = var.cloudflare_zone_name
}

# ACM certificate for the custom domain (must be in same region as API Gateway)
resource "aws_acm_certificate" "main" {
  domain_name       = var.domain_name
  validation_method = "DNS"

  lifecycle {
    create_before_destroy = true
  }
}

# Add ACM validation CNAME records to Cloudflare
resource "cloudflare_record" "acm_validation" {
  for_each = {
    for dvo in aws_acm_certificate.main.domain_validation_options : dvo.domain_name => {
      name  = dvo.resource_record_name
      value = dvo.resource_record_value
      type  = dvo.resource_record_type
    }
  }

  zone_id = data.cloudflare_zone.main.id
  name    = each.value.name
  content = each.value.value
  type    = each.value.type
  ttl     = 60
  proxied = false # Must be false — ACM validates the record directly
}

# Wait for ACM to confirm the certificate is valid
resource "aws_acm_certificate_validation" "main" {
  certificate_arn         = aws_acm_certificate.main.arn
  validation_record_fqdns = [for record in cloudflare_record.acm_validation : record.hostname]
}

# API Gateway custom domain
resource "aws_apigatewayv2_domain_name" "main" {
  domain_name = var.domain_name

  domain_name_configuration {
    certificate_arn = aws_acm_certificate_validation.main.certificate_arn
    endpoint_type   = "REGIONAL"
    security_policy = "TLS_1_2"
  }
}

# Map the API Gateway stage to the custom domain
resource "aws_apigatewayv2_api_mapping" "main" {
  api_id      = aws_apigatewayv2_api.main.id
  domain_name = aws_apigatewayv2_domain_name.main.id
  stage       = aws_apigatewayv2_stage.default.id
}

# CNAME: gotodo.amohr.net → API Gateway regional domain
resource "cloudflare_record" "app" {
  zone_id = data.cloudflare_zone.main.id
  name    = "gotodo"
  content = aws_apigatewayv2_domain_name.main.domain_name_configuration[0].target_domain_name
  type    = "CNAME"
  proxied = true # Routes through Cloudflare's network; set false for simple passthrough
}
