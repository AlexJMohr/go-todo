locals {
  # Use the pooler endpoint — Lambda creates a new connection per cold start,
  # pooler prevents connection exhaustion on Neon's free tier
  database_url = "postgresql://${neon_role.app.name}:${neon_role.app.password}@${neon_project.main.database_host_pooler}/${neon_database.main.name}?sslmode=require"
}

resource "aws_iam_role" "lambda" {
  name = "${var.app_name}-lambda"
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect    = "Allow"
      Principal = { Service = "lambda.amazonaws.com" }
      Action    = "sts:AssumeRole"
    }]
  })
}

resource "aws_iam_role_policy" "lambda_logs" {
  name = "${var.app_name}-lambda-logs"
  role = aws_iam_role.lambda.id
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect = "Allow"
      Action = [
        "logs:CreateLogGroup",
        "logs:CreateLogStream",
        "logs:PutLogEvents",
      ]
      Resource = "arn:aws:logs:*:*:*"
    }]
  })
}

resource "aws_cloudwatch_log_group" "lambda" {
  name              = "/aws/lambda/${var.app_name}"
  retention_in_days = 14
}

data "archive_file" "placeholder" {
  type        = "zip"
  output_path = "${path.module}/placeholder.zip"
  source {
    content  = "placeholder"
    filename = "bootstrap"
  }
}

resource "aws_lambda_function" "app" {
  function_name = var.app_name
  role          = aws_iam_role.lambda.arn
  runtime       = "provided.al2023"
  handler       = "bootstrap"
  architectures = ["arm64"]
  filename      = data.archive_file.placeholder.output_path

  # Prevent Terraform from overwriting code deployed by GitHub Actions
  lifecycle {
    ignore_changes = [filename, source_code_hash]
  }

  environment {
    variables = {
      DATABASE_URL = local.database_url
      GIN_MODE     = "release"
    }
  }

  depends_on = [aws_cloudwatch_log_group.lambda]
}
