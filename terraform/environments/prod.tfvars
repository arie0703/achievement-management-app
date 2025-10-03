# Production environment configuration
# High availability and performance optimized configuration
# Resource sizing: Production-grade resources for high availability and performance

environment = "prod"
aws_region  = "us-east-1"
cost_center = "engineering-prod"
app_name    = "achievement-management"

# VPC Configuration - production network layout
vpc_cidr             = "10.2.0.0/16"
availability_zones   = ["us-east-1a", "us-east-1b"]
public_subnet_cidrs  = ["10.2.1.0/24", "10.2.2.0/24"]
private_subnet_cidrs = ["10.2.10.0/24", "10.2.20.0/24"]

# ECS Configuration - production-grade resources
# CPU: 1 vCPU, Memory: 2 GB - high-performance configuration for production workloads
ecs_config = {
  cpu                = 1024
  memory             = 2048
  desired_count      = 3
  max_capacity       = 20
  min_capacity       = 3
  container_port     = 8080
  health_check_path  = "/health"
  container_image    = "achievement-management:latest"
}

# DynamoDB Configuration - provisioned capacity for predictable performance
dynamodb_tables = {
  achievements = {
    hash_key               = "user_id"
    range_key              = "achievement_id"
    billing_mode           = "PROVISIONED"
    read_capacity          = 10
    write_capacity         = 5
    point_in_time_recovery = true
    server_side_encryption = true
  }
  rewards = {
    hash_key               = "reward_id"
    billing_mode           = "PROVISIONED"
    read_capacity          = 5
    write_capacity         = 2
    point_in_time_recovery = true
    server_side_encryption = true
  }
  current_points = {
    hash_key               = "user_id"
    billing_mode           = "PROVISIONED"
    read_capacity          = 15
    write_capacity         = 10
    point_in_time_recovery = true
    server_side_encryption = true
  }
  reward_history = {
    hash_key               = "user_id"
    range_key              = "timestamp"
    billing_mode           = "PROVISIONED"
    read_capacity          = 8
    write_capacity         = 5
    point_in_time_recovery = true
    server_side_encryption = true
  }
}

# Load Balancer Configuration - HTTPS with redirect and deletion protection for production
alb_config = {
  enable_https               = true
  enable_https_redirect      = true
  certificate_arn           = ""  # To be provided via environment variable or updated manually
  ssl_policy                = "ELBSecurityPolicy-TLS-1-2-2017-01"
  enable_deletion_protection = true
}

# Monitoring Configuration - extended retention for production
log_retention_days = 30
cli_task_types     = ["achievement", "points", "reward"]
enable_monitoring_dashboard = true

# Monitoring thresholds - strict for production
monitoring_thresholds = {
  cpu_alarm_threshold         = 70  # Conservative threshold for production
  memory_alarm_threshold      = 75  # Conservative threshold for production
  response_time_threshold     = 2.0 # Strict response time requirement
  http_5xx_threshold          = 5   # Low error tolerance
  dynamodb_throttle_threshold = 0   # No throttling allowed in production
}

# Alarm actions for production - should be configured with SNS topics for notifications
alarm_actions = [
  # "arn:aws:sns:us-east-1:123456789012:production-alerts"
]

# Backend Configuration
backend_bucket         = "achievement-management-terraform-state-prod"
backend_key_prefix     = "prod"
backend_dynamodb_table = "terraform-state-lock-prod"

# Secret Management Configuration
# Production environment uses AWS Secrets Manager with automatic rotation
secrets_manager_enabled = true
database_password_secret_name = "achievement-management/prod/db-password"
api_key_secret_name = "achievement-management/prod/api-keys"

# Additional Production-specific Configuration
enable_debug_logging = false
enable_xray_tracing = true
enable_container_insights = true

# Enhanced Tagging Configuration for Production Environment
owner = "platform-team"
team = "engineering"
business_unit = "technology"
service_tier = "critical"
backup_policy = "daily"
data_classification = "confidential"
compliance_scope = "production-soc2"

# Resource Naming Configuration
resource_naming_convention = {
  include_region     = false
  include_account_id = false
  separator         = "-"
  max_length        = 64
}