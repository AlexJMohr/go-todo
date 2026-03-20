output "url" {
  description = "Public URL for the deployed application"
  value       = "https://${var.domain_name}"
}

output "api_gateway_url" {
  description = "Direct API Gateway URL (bypasses custom domain)"
  value       = aws_apigatewayv2_stage.default.invoke_url
}
