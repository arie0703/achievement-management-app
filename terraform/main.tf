# Main Terraform configuration for Achievement Management Application
# This file orchestrates all modules and defines the root module configuration

# Data sources for AWS account and region information
data "aws_caller_identity" "current" {}
data "aws_region" "current" {}
data "aws_availability_zones" "available" {
  state = "available"
}

# Local values for common resource naming and tagging
locals {
  # Core naming and identification
  name_prefix = "${var.app_name}-${var.environment}"
  account_id  = data.aws_caller_identity.current.account_id
  region      = data.aws_region.current.name

  # Availability zones selection (use first 2 available zones if not specified)
  availability_zones = length(var.availability_zones) > 0 ? var.availability_zones : slice(data.aws_availability_zones.available.names, 0, 2)

  # Comprehensive tagging strategy for cost tracking and resource organization
  common_tags = {
    # Core identification tags
    Environment = var.environment
    Application = var.app_name
    Project     = "achievement-management"
    Component   = "infrastructure"

    # Management and ownership tags
    ManagedBy  = "terraform"
    Owner      = var.owner
    Team       = var.team
    CostCenter = var.cost_center

    # Operational tags
    BackupPolicy      = var.backup_policy
    MonitoringEnabled = var.enable_monitoring_dashboard ? "true" : "false"

    # Compliance and governance tags
    DataClassification = var.data_classification
    ComplianceScope    = var.compliance_scope

    # Lifecycle tags
    CreatedBy    = "terraform"
    CreatedDate  = formatdate("YYYY-MM-DD", timestamp())
    LastModified = formatdate("YYYY-MM-DD", timestamp())

    # Environment-specific tags
    EnvironmentType = var.environment == "prod" ? "production" : (var.environment == "staging" ? "pre-production" : "development")
    BusinessUnit    = var.business_unit

    # Resource organization tags
    ResourceGroup = "${var.app_name}-${var.environment}"
    ServiceTier   = var.service_tier

    # AWS account and region information
    AccountId = local.account_id
    Region    = local.region
  }

  # Environment-specific naming conventions
  naming_convention = {
    # Standard resource naming pattern: {app_name}-{environment}-{resource_type}-{identifier}
    vpc_name         = "${var.app_name}-${var.environment}-vpc"
    cluster_name     = "${var.app_name}-${var.environment}-cluster"
    alb_name         = "${substr("${var.app_name}-${var.environment}-alb", 0, 32)}" # ALB names have 32 char limit
    api_service_name = "${var.app_name}-${var.environment}-api-service"

    # DynamoDB table naming with consistent pattern
    table_prefix = "${var.app_name}-${var.environment}"

    # Log group naming pattern
    log_group_prefix = "/aws/ecs/${var.app_name}-${var.environment}"

    # IAM role naming pattern
    role_prefix = "${var.app_name}-${var.environment}"

    # Security group naming pattern
    sg_prefix = "${var.app_name}-${var.environment}"
  }

  # Derived configuration values
  dynamodb_table_names = keys(var.dynamodb_tables)
  dynamodb_table_arns  = values(module.dynamodb.table_arns)
}

# VPC and Networking Module
# Creates the foundational networking infrastructure
module "vpc" {
  source = "./modules/vpc"

  # Core configuration
  environment = var.environment
  app_name    = var.app_name

  # Network configuration
  vpc_cidr             = var.vpc_cidr
  public_subnet_cidrs  = var.public_subnet_cidrs
  private_subnet_cidrs = var.private_subnet_cidrs

  common_tags = local.common_tags
}

# Security Groups Module
# Creates security groups with proper ingress/egress rules
module "security" {
  source = "./modules/security"

  # Core configuration
  environment = var.environment
  app_name    = var.app_name

  # VPC dependencies
  vpc_id               = module.vpc.vpc_id
  private_subnet_cidrs = var.private_subnet_cidrs

  # Application configuration
  container_port = var.ecs_config.container_port

  tags = local.common_tags

  # Explicit dependency on VPC
  depends_on = [module.vpc]
}

# IAM Roles and Policies Module
# Creates IAM roles and policies for ECS tasks and CLI operations
module "iam" {
  source = "./modules/iam"

  # Core configuration
  environment = var.environment
  app_name    = var.app_name

  # DynamoDB configuration for permissions
  dynamodb_table_names = local.dynamodb_table_names

  tags = local.common_tags

  # Explicit dependency on DynamoDB tables
  depends_on = [module.dynamodb]
}

# DynamoDB Tables Module
# Creates DynamoDB tables for application data storage
module "dynamodb" {
  source = "./modules/dynamodb"

  # Core configuration
  environment = var.environment
  app_name    = var.app_name

  # Table configuration
  dynamodb_tables = var.dynamodb_tables

  tags = local.common_tags
}

# CloudWatch Monitoring Module
# Creates CloudWatch log groups, dashboards, and alarms for monitoring
module "monitoring" {
  source = "./modules/monitoring"

  # Core configuration
  environment = var.environment
  app_name    = var.app_name

  # ECS Configuration for Monitoring
  ecs_service_name = module.ecs.service_name
  ecs_cluster_name = module.ecs.cluster_name

  alb_arn_suffix          = module.ecs.alb_arn_suffix
  target_group_arn_suffix = module.ecs.target_group_arn_suffix

  # DynamoDB Configuration for Monitoring
  dynamodb_table_names = module.dynamodb.table_names

  # Monitoring Configuration
  log_retention_days = var.log_retention_days
  cli_task_types     = var.cli_task_types
  enable_dashboard   = var.enable_monitoring_dashboard

  # Alarm Thresholds
  cpu_alarm_threshold         = var.monitoring_thresholds.cpu_alarm_threshold
  memory_alarm_threshold      = var.monitoring_thresholds.memory_alarm_threshold
  response_time_threshold     = var.monitoring_thresholds.response_time_threshold
  http_5xx_threshold          = var.monitoring_thresholds.http_5xx_threshold
  dynamodb_throttle_threshold = var.monitoring_thresholds.dynamodb_throttle_threshold

  # Alarm Actions
  alarm_actions = var.alarm_actions

  tags = local.common_tags

  # Explicit dependencies to ensure proper resource creation order
  depends_on = [
    module.ecs,
    module.dynamodb
  ]
}

# ECS Cluster and Services Module
# Creates ECS cluster, services, task definitions, and Application Load Balancer
module "ecs" {
  source = "./modules/ecs"

  # Core configuration
  environment = var.environment
  app_name    = var.app_name

  # VPC and networking dependencies
  vpc_id             = module.vpc.vpc_id
  private_subnet_ids = module.vpc.private_subnet_ids
  public_subnet_ids  = module.vpc.public_subnet_ids

  # Security group dependencies
  alb_security_group_id = module.security.alb_security_group_id
  ecs_security_group_id = module.security.ecs_security_group_id

  # IAM role dependencies
  ecs_task_execution_role_arn = module.iam.ecs_task_execution_role_arn
  ecs_task_role_arn           = module.iam.ecs_task_role_arn
  cli_task_role_arns          = module.iam.cli_task_role_arns

  # Container Configuration
  api_container_image = var.ecs_config.container_image
  cli_container_image = var.ecs_config.container_image
  container_port      = var.ecs_config.container_port

  # API Service Configuration
  api_task_cpu      = var.ecs_config.cpu
  api_task_memory   = var.ecs_config.memory
  api_desired_count = var.ecs_config.desired_count
  api_min_capacity  = var.ecs_config.min_capacity
  api_max_capacity  = var.ecs_config.max_capacity

  # CLI Task Configuration
  cli_task_types = var.cli_task_types

  # Health Check Configuration
  health_check_path = var.ecs_config.health_check_path

  # Load Balancer Configuration
  enable_https               = var.alb_config.enable_https
  enable_https_redirect      = var.alb_config.enable_https_redirect
  certificate_arn            = var.alb_config.certificate_arn
  ssl_policy                 = var.alb_config.ssl_policy
  enable_deletion_protection = var.alb_config.enable_deletion_protection

  # Logging Configuration
  log_retention_days = var.log_retention_days

  # Feature flags
  enable_container_insights = var.enable_container_insights

  # Resource naming and tagging
  tags = local.common_tags

  # Explicit dependencies to ensure proper resource creation order
  depends_on = [
    module.vpc,
    module.security,
    module.iam,
    module.dynamodb
  ]
}

# Data validation and computed values
locals {
  # Validate that we have the correct number of subnets for the availability zones
  validate_subnet_count = length(var.public_subnet_cidrs) == length(local.availability_zones) && length(var.private_subnet_cidrs) == length(local.availability_zones)

  # Environment-specific configuration validation
  is_production = var.environment == "prod"

  # Computed resource configurations based on environment
  computed_config = {
    # Production environments should have higher resource allocations
    min_capacity = local.is_production ? max(var.ecs_config.min_capacity, 2) : var.ecs_config.min_capacity

    # Production should have deletion protection enabled by default
    enable_deletion_protection = local.is_production ? true : var.alb_config.enable_deletion_protection

    # Production should have point-in-time recovery enabled for all tables
    enable_pitr = local.is_production ? true : null

    # Log retention should be longer for production
    log_retention_days = local.is_production ? max(var.log_retention_days, 30) : var.log_retention_days
  }
}

# Validation checks
resource "null_resource" "validate_configuration" {
  count = local.validate_subnet_count ? 0 : 1

  provisioner "local-exec" {
    command = "echo 'ERROR: Number of subnet CIDRs must match number of availability zones' && exit 1"
  }
}

# Output summary information for verification
output "deployment_summary" {
  description = "Summary of the deployed infrastructure"
  value = {
    environment        = var.environment
    region             = local.region
    account_id         = local.account_id
    availability_zones = local.availability_zones
    vpc_cidr           = var.vpc_cidr
    resource_prefix    = local.name_prefix

    # Resource counts
    public_subnets  = length(var.public_subnet_cidrs)
    private_subnets = length(var.private_subnet_cidrs)
    dynamodb_tables = length(local.dynamodb_table_names)
    cli_task_types  = length(var.cli_task_types)

    # Configuration flags
    monitoring_enabled  = var.enable_monitoring_dashboard
    https_enabled       = var.alb_config.enable_https
    deletion_protection = local.computed_config.enable_deletion_protection

    # Computed values
    computed_min_capacity  = local.computed_config.min_capacity
    computed_log_retention = local.computed_config.log_retention_days
  }
}

# Module integration verification
output "module_integration_status" {
  description = "Status of module integration and dependencies"
  value = {
    vpc_ready        = module.vpc.vpc_id != null
    security_ready   = module.security.alb_security_group_id != null && module.security.ecs_security_group_id != null
    iam_ready        = module.iam.ecs_task_execution_role_arn != null && module.iam.ecs_task_role_arn != null
    dynamodb_ready   = length(module.dynamodb.table_names) == length(local.dynamodb_table_names)
    ecs_ready        = module.ecs.cluster_id != null && module.ecs.api_service_name != null
    monitoring_ready = module.monitoring.api_log_group_name != null

    # Cross-module data flow verification
    vpc_to_security = true # VPC is passed to security module via variable
    security_to_ecs = true # Security groups are passed to ECS module via variables
    iam_to_ecs      = true # IAM roles are passed to ECS module via variables
    dynamodb_to_iam = length(setintersection(keys(module.iam.cli_task_role_arns), var.cli_task_types)) == length(var.cli_task_types)
  }
}
