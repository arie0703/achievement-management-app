# Staging environment configuration
# Production-like configuration for testing and validation
# Resource sizing: Moderate resources to simulate production workloads

environment = "staging"
aws_region  = "us-east-1"
cost_center = "engineering-staging"
app_name    = "achievement-management"

# VPC Configuration - same as production
vpc_cidr             = "10.1.0.0/16"
availability_zones   = ["us-east-1a", "us-east-1b"]
public_subnet_cidrs  = ["10.1.1.0/24", "10.1.2.0/24"]
private_subnet_cidrs = ["10.1.10.0/24", "10.1.20.0/24"]

# ECS Configuration - moderate resources for staging
# CPU: 0.5 vCPU, Memory: 1 GB - production-like sizing for testing
ecs_config = {
  cpu                = 512
  memory             = 1024
  desired_count      = 2
  max_capacity       = 5
  min_capacity       = 2
  container_port     = 8080
  health_check_path  = "/health"
  container_image    = "achievement-management:staging"
}

# DynamoDB Configuration - on-demand with backup enabled
dynamodb_tables = {
  achievements = {
    hash_key               = "user_id"
    range_key              = "achievement_id"
    billing_mode           = "PAY_PER_REQUEST"
    point_in_time_recovery = true
    server_side_encryption = true
  }
  rewards = {
    hash_key               = "reward_id"
    billing_mode           = "PAY_PER_REQUEST"
    point_in_time_recovery = true
    server_side_encryption = true
  }
  current_points = {
    hash_key               = "user_id"
    billing_mode           = "PAY_PER_REQUEST"
    point_in_time_recovery = true
    server_side_encryption = true
  }
  reward_history = {
    hash_key               = "user_id"
    range_key              = "timestamp"
    billing_mode           = "PAY_PER_REQUEST"
    point_in_time_recovery = true
    server_side_encryption = true
  }
}

# Load Balancer Configuration - HTTPS with redirect for staging
alb_config = {
  enable_https               = true
  enable_https_redirect      = true
  certificate_arn           = ""  # To be provided via environment variable or updated manually
  ssl_policy                = "ELBSecurityPolicy-TLS-1-2-2017-01"
  enable_deletion_protection = false
}

# Monitoring Configuration - moderate retention for staging
log_retention_days = 14
cli_task_types     = ["achievement", "points", "reward"]
enable_monitoring_dashboard = true

# Monitoring thresholds - production-like for staging
monitoring_thresholds = {
  cpu_alarm_threshold         = 80  # Production-like threshold
  memory_alarm_threshold      = 80  # Production-like threshold
  response_time_threshold     = 3.0 # Moderate response time
  http_5xx_threshold          = 10  # Production-like error threshold
  dynamodb_throttle_threshold = 1   # Low throttling tolerance
}

# Alarm actions for staging - can be configured with SNS topics
alarm_actions = []

# Backend Configuration
backend_bucket         = "achievement-management-terraform-state-staging"
backend_key_prefix     = "staging"
backend_dynamodb_table = "terraform-state-lock-staging"

# Secret Management Configuration
# Staging environment uses AWS Secrets Manager for production-like testing
secrets_manager_enabled = true
database_password_secret_name = "achievement-management/staging/db-password"
api_key_secret_name = "achievement-management/staging/api-keys"

# Additional Staging-specific Configuration
enable_debug_logging = true
enable_xray_tracing = true
enable_container_insights = true

# Enhanced Tagging Configuration for Staging Environment
owner = "platform-team"
team = "engineering"
business_unit = "technology"
service_tier = "important"
backup_policy = "daily"
data_classification = "internal"
compliance_scope = "pre-production-testing"

# Resource Naming Configuration
resource_naming_convention = {
  include_region     = false
  include_account_id = false
  separator         = "-"
  max_length        = 64
}