resource "neon_project" "main" {
  name                       = var.app_name
  region_id                  = "aws-us-west-2"
  org_id                     = var.neon_org_id
  history_retention_seconds  = 21600 # 6 hours — free tier maximum
}

resource "neon_role" "app" {
  project_id = neon_project.main.id
  branch_id  = neon_project.main.default_branch_id
  name       = "${var.app_name}-app"
}

resource "neon_database" "main" {
  project_id = neon_project.main.id
  branch_id  = neon_project.main.default_branch_id
  name       = "go_todo"
  owner_name = neon_role.app.name
}
