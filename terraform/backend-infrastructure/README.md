# Terraform Backend Infrastructure

This directory contains the Terraform configuration for creating the backend infrastructure needed for remote state storage and locking. This includes S3 buckets for state storage and DynamoDB tables for state locking.

## Overview

The backend infrastructure creates:
- S3 buckets for each environment (dev, staging, prod) with versioning and encryption
- DynamoDB tables for each environment for state locking
- Proper security configurations including public access blocking and SSL-only policies

## Prerequisites

- AWS CLI configured with appropriate permissions
- Terraform >= 1.0 installed

## Required AWS Permissions

The AWS credentials used must have permissions to:
- Create and manage S3 buckets
- Create and manage DynamoDB tables
- Apply bucket policies and lifecycle configurations

## Deployment Instructions

### 1. Initialize and Deploy Backend Infrastructure

```bash
# Navigate to the backend infrastructure directory
cd terraform/backend-infrastructure

# Initialize Terraform (uses local state)
terraform init

# Review the planned changes
terraform plan

# Apply the configuration
terraform apply
```

### 2. Note the Outputs

After deployment, note the output values which will confirm the created resources match the backend configuration files.

### 3. Initialize Main Infrastructure with Remote Backend

After the backend infrastructure is deployed, you can initialize the main infrastructure with remote state:

```bash
# Navigate back to the main terraform directory
cd ..

# Initialize with the appropriate backend configuration
terraform init -backend-config=backend-dev.hcl     # For development
terraform init -backend-config=backend-staging.hcl # For staging  
terraform init -backend-config=backend-prod.hcl    # For production
```

## Environment-Specific Resources

### Development Environment
- S3 Bucket: `achievement-management-terraform-state-dev`
- DynamoDB Table: `terraform-state-lock-dev`
- State Key: `dev/terraform.tfstate`

### Staging Environment
- S3 Bucket: `achievement-management-terraform-state-staging`
- DynamoDB Table: `terraform-state-lock-staging`
- State Key: `staging/terraform.tfstate`

### Production Environment
- S3 Bucket: `achievement-management-terraform-state-prod`
- DynamoDB Table: `terraform-state-lock-prod`
- State Key: `prod/terraform.tfstate`

## Security Features

- **Encryption**: All S3 buckets use server-side encryption (AES256)
- **Versioning**: Enabled on all state buckets for rollback capability
- **Public Access**: Blocked on all buckets
- **SSL Only**: Bucket policies enforce HTTPS connections only
- **Lifecycle Management**: Old versions are automatically cleaned up after 30 days

## State Management

- **Local State**: This backend infrastructure uses local state to avoid circular dependency
- **Remote State**: The main infrastructure will use the created S3 buckets for remote state
- **Locking**: DynamoDB tables prevent concurrent modifications
- **Backup**: S3 versioning provides automatic backup of state files

## Troubleshooting

### Bucket Already Exists Error
If you get a bucket already exists error, either:
1. Choose a different `app_name` variable value
2. Delete the existing bucket if it's safe to do so
3. Import the existing bucket into Terraform state

### Permission Errors
Ensure your AWS credentials have the necessary permissions listed in the prerequisites section.

### State Lock Issues
If you encounter state lock issues:
1. Check the DynamoDB table exists and is accessible
2. Verify the table name matches the backend configuration
3. Use `terraform force-unlock <lock-id>` if necessary (use with caution)

## Cleanup

To destroy the backend infrastructure (only do this if you're sure):

```bash
# Make sure no infrastructure is using these backends
terraform destroy
```

**Warning**: Only destroy backend infrastructure after ensuring no other Terraform configurations are using it for remote state storage.