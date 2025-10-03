# Development environment configuration
# Optimized for cost and development workflows
# Resource sizing: Minimal resources for development and testing

environment = "dev"
aws_region  = "us-east-1"
cost_center = "engineering-dev"
app_name    = "achievement-management"

# VPC Configuration - smaller CIDR blocks for dev
vpc_cidr             = "10.0.0.0/16"
availability_zones   = ["us-east-1a", "us-east-1b"]
public_subnet_cidrs  = ["10.0.1.0/24", "10.0.2.0/24"]
private_subnet_cidrs = ["10.0.10.0/24", "10.0.20.0/24"]

# ECS Configuration - minimal resources for development
# CPU: 0.25 vCPU, Memory: 512 MB - cost-optimized for development
ecs_config = {
  cpu                = 256
  memory             = 512
  desired_count      = 1
  max_capacity       = 2
  min_capacity       = 1
  container_port     = 8080
  health_check_path  = "/health"
  container_image    = "achievement-management:dev"
}

# DynamoDB Configuration - on-demand billing for variable dev workloads
dynamodb_tables = {
  achievements = {
    hash_key               = "user_id"
    range_key              = "achievement_id"
    billing_mode           = "PAY_PER_REQUEST"
    point_in_time_recovery = false  # Disabled for cost savings in dev
    server_side_encryption = true
  }
  rewards = {
    hash_key               = "reward_id"
    billing_mode           = "PAY_PER_REQUEST"
    point_in_time_recovery = false
    server_side_encryption = true
  }
  current_points = {
    hash_key               = "user_id"
    billing_mode           = "PAY_PER_REQUEST"
    point_in_time_recovery = false
    server_side_encryption = true
  }
  reward_history = {
    hash_key               = "user_id"
    range_key              = "timestamp"
    billing_mode           = "PAY_PER_REQUEST"
    point_in_time_recovery = false
    server_side_encryption = true
  }
}

# Load Balancer Configuration - HTTP only for dev
alb_config = {
  enable_https               = false
  enable_https_redirect      = false
  certificate_arn           = ""
  ssl_policy                = "ELBSecurityPolicy-TLS-1-2-2017-01"
  enable_deletion_protection = false
}

# Monitoring Configuration - shorter retention for dev
log_retention_days = 7
cli_task_types     = ["achievement", "points", "reward"]
enable_monitoring_dashboard = true

# Monitoring thresholds - relaxed for development
monitoring_thresholds = {
  cpu_alarm_threshold         = 90  # Higher threshold for dev
  memory_alarm_threshold      = 90  # Higher threshold for dev
  response_time_threshold     = 5.0 # More lenient response time
  http_5xx_threshold          = 20  # Allow more errors in dev
  dynamodb_throttle_threshold = 5   # Allow some throttling in dev
}

# No alarm actions in dev environment
alarm_actions = []

# Backend Configuration
backend_bucket         = "achievement-management-terraform-state-dev"
backend_key_prefix     = "dev"
backend_dynamodb_table = "terraform-state-lock-dev"

# Secret Management Configuration
# Development environment uses local secrets or environment variables
secrets_manager_enabled = false
database_password_secret_name = ""
api_key_secret_name = ""

# Additional Development-specific Configuration
enable_debug_logging = true
enable_xray_tracing = false
enable_container_insights = false

# Enhanced Tagging Configuration for Development Environment
owner = "development-team"
team = "engineering"
business_unit = "technology"
service_tier = "standard"
backup_policy = "none"  # No backup needed for dev environment
data_classification = "internal"
compliance_scope = "development-only"

# Resource Naming Configuration
resource_naming_convention = {
  include_region     = false
  include_account_id = false
  separator         = "-"
  max_length        = 64
}