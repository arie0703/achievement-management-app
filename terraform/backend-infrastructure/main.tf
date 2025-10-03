# Backend Infrastructure for Terraform State Management
# This configuration creates the S3 buckets and DynamoDB tables needed
# for Terraform remote state storage and locking.
#
# This should be deployed BEFORE the main infrastructure and uses local state.
# After deployment, the main infrastructure can use these resources for remote state.

terraform {
  required_version = ">= 1.0"
  
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
  
  # This configuration uses local state since it creates the backend infrastructure
  # Do not configure a remote backend here to avoid circular dependency
}

provider "aws" {
  region = var.aws_region
  
  default_tags {
    tags = {
      Application = var.app_name
      ManagedBy   = "terraform"
      Project     = "achievement-management"
      Purpose     = "terraform-backend"
    }
  }
}

# Local values for resource naming
locals {
  environments = ["dev", "staging", "prod"]
  
  # Common tags for all backend resources
  common_tags = {
    Application = var.app_name
    ManagedBy   = "terraform"
    Project     = "achievement-management"
    Purpose     = "terraform-backend"
  }
}
# S3 B
# buckets for Terraform State Storage
# Creates one bucket per environment with versioning and encryption enabled
resource "aws_s3_bucket" "terraform_state" {
  for_each = toset(local.environments)
  
  bucket = "${var.app_name}-terraform-state-${each.value}"
  
  tags = merge(local.common_tags, {
    Environment = each.value
    Name        = "${var.app_name}-terraform-state-${each.value}"
  })
}

# Enable versioning on state buckets
resource "aws_s3_bucket_versioning" "terraform_state" {
  for_each = aws_s3_bucket.terraform_state
  
  bucket = each.value.id
  versioning_configuration {
    status = "Enabled"
  }
}

# Enable server-side encryption on state buckets
resource "aws_s3_bucket_server_side_encryption_configuration" "terraform_state" {
  for_each = aws_s3_bucket.terraform_state
  
  bucket = each.value.id
  
  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm = "AES256"
    }
    bucket_key_enabled = true
  }
}

# Block public access to state buckets
resource "aws_s3_bucket_public_access_block" "terraform_state" {
  for_each = aws_s3_bucket.terraform_state
  
  bucket = each.value.id
  
  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

# Enable lifecycle configuration to manage old versions
resource "aws_s3_bucket_lifecycle_configuration" "terraform_state" {
  for_each = aws_s3_bucket.terraform_state
  
  bucket = each.value.id
  
  rule {
    id     = "terraform_state_lifecycle"
    status = "Enabled"
    
    # Keep non-current versions for 30 days
    noncurrent_version_expiration {
      noncurrent_days = 30
    }
    
    # Delete incomplete multipart uploads after 7 days
    abort_incomplete_multipart_upload {
      days_after_initiation = 7
    }
  }
}
# DynamoDB Tables for Terraform State Locking
# Creates one table per environment for state locking
resource "aws_dynamodb_table" "terraform_state_lock" {
  for_each = toset(local.environments)
  
  name           = "terraform-state-lock-${each.value}"
  billing_mode   = "PAY_PER_REQUEST"
  hash_key       = "LockID"
  
  attribute {
    name = "LockID"
    type = "S"
  }
  
  # Enable point-in-time recovery for production
  point_in_time_recovery {
    enabled = each.value == "prod" ? true : false
  }
  
  # Enable server-side encryption
  server_side_encryption {
    enabled = true
  }
  
  tags = merge(local.common_tags, {
    Environment = each.value
    Name        = "terraform-state-lock-${each.value}"
  })
}

# Optional: Create a bucket policy to restrict access to the state buckets
resource "aws_s3_bucket_policy" "terraform_state" {
  for_each = aws_s3_bucket.terraform_state
  
  bucket = each.value.id
  
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid       = "DenyInsecureConnections"
        Effect    = "Deny"
        Principal = "*"
        Action    = "s3:*"
        Resource = [
          each.value.arn,
          "${each.value.arn}/*"
        ]
        Condition = {
          Bool = {
            "aws:SecureTransport" = "false"
          }
        }
      }
    ]
  })
}