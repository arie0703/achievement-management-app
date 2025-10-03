# DynamoDB Tables Module
# This module creates DynamoDB tables for the achievement management application
# with proper configuration for data access patterns, security, and backup

# Local values for resource naming and configuration
locals {
  name_prefix = "${var.app_name}-${var.environment}"

  # Separate achievements table to handle GSI
  non_achievements_tables = {
    for key, value in var.dynamodb_tables : key => value
    if key != "achievements"
  }
}

# Standard DynamoDB Tables (without GSI)
resource "aws_dynamodb_table" "tables" {
  for_each = local.non_achievements_tables

  name         = "${local.name_prefix}-${each.key}"
  billing_mode = each.value.billing_mode
  hash_key     = each.value.hash_key
  range_key    = each.value.range_key

  # Provisioned throughput (only used when billing_mode is PROVISIONED)
  read_capacity  = each.value.billing_mode == "PROVISIONED" ? each.value.read_capacity : null
  write_capacity = each.value.billing_mode == "PROVISIONED" ? each.value.write_capacity : null

  # Hash key attribute definition
  attribute {
    name = each.value.hash_key
    type = "S" # String type for all keys in this application
  }

  # Range key attribute definition (if specified)
  dynamic "attribute" {
    for_each = each.value.range_key != null ? [each.value.range_key] : []
    content {
      name = attribute.value
      type = "S" # String type for all keys in this application
    }
  }

  # Point-in-time recovery configuration
  point_in_time_recovery {
    enabled = each.value.point_in_time_recovery
  }

  # Server-side encryption configuration
  server_side_encryption {
    enabled = each.value.server_side_encryption
  }

  # Deletion protection for production environments
  deletion_protection_enabled = var.environment == "prod" ? true : false

  # Resource tags
  tags = merge(var.tags, {
    Name              = "${local.name_prefix}-${each.key}"
    ResourceType      = "dynamodb-table"
    TableType         = each.key
    BillingMode       = each.value.billing_mode
    BackupEnabled     = each.value.point_in_time_recovery ? "true" : "false"
    EncryptionEnabled = each.value.server_side_encryption ? "true" : "false"
    DataStore         = "primary"
    AccessPattern     = each.key == "current_points" ? "high-frequency" : "standard"
  })
}

# Achievements table (simplified schema)
resource "aws_dynamodb_table" "achievements" {
  count = contains(keys(var.dynamodb_tables), "achievements") ? 1 : 0

  name         = "${local.name_prefix}-achievements"
  billing_mode = var.dynamodb_tables.achievements.billing_mode
  hash_key     = var.dynamodb_tables.achievements.hash_key

  # Provisioned throughput (only used when billing_mode is PROVISIONED)
  read_capacity  = var.dynamodb_tables.achievements.billing_mode == "PROVISIONED" ? var.dynamodb_tables.achievements.read_capacity : null
  write_capacity = var.dynamodb_tables.achievements.billing_mode == "PROVISIONED" ? var.dynamodb_tables.achievements.write_capacity : null

  # Hash key attribute definition (id)
  attribute {
    name = var.dynamodb_tables.achievements.hash_key
    type = "S"
  }

  # Point-in-time recovery configuration
  point_in_time_recovery {
    enabled = var.dynamodb_tables.achievements.point_in_time_recovery
  }

  # Server-side encryption configuration
  server_side_encryption {
    enabled = var.dynamodb_tables.achievements.server_side_encryption
  }

  # Deletion protection for production environments
  deletion_protection_enabled = var.environment == "prod" ? true : false

  # Resource tags
  tags = merge(var.tags, {
    Name              = "${local.name_prefix}-achievements"
    ResourceType      = "dynamodb-table"
    TableType         = "achievements"
    BillingMode       = var.dynamodb_tables.achievements.billing_mode
    BackupEnabled     = var.dynamodb_tables.achievements.point_in_time_recovery ? "true" : "false"
    EncryptionEnabled = var.dynamodb_tables.achievements.server_side_encryption ? "true" : "false"
    DataStore         = "primary"
    AccessPattern     = "standard"
    HasGSI            = "true"
  })
}
