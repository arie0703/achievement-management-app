# Terraform Backend Configuration

This document describes the Terraform backend configuration for the Achievement Management infrastructure, including remote state storage and state locking mechanisms.

## Overview

The Terraform backend configuration uses:
- **S3 buckets** for remote state storage with versioning and encryption
- **DynamoDB tables** for state locking to prevent concurrent modifications
- **Environment-specific organization** to isolate state files per environment

## Architecture

```
Backend Infrastructure
├── S3 Buckets (State Storage)
│   ├── achievement-management-terraform-state-dev
│   ├── achievement-management-terraform-state-staging
│   └── achievement-management-terraform-state-prod
└── DynamoDB Tables (State Locking)
    ├── terraform-state-lock-dev
    ├── terraform-state-lock-staging
    └── terraform-state-lock-prod
```

## Environment-Specific Configuration

### Development Environment
- **S3 Bucket**: `achievement-management-terraform-state-dev`
- **State Key**: `dev/terraform.tfstate`
- **DynamoDB Table**: `terraform-state-lock-dev`
- **Region**: `us-east-1`

### Staging Environment
- **S3 Bucket**: `achievement-management-terraform-state-staging`
- **State Key**: `staging/terraform.tfstate`
- **DynamoDB Table**: `terraform-state-lock-staging`
- **Region**: `us-east-1`

### Production Environment
- **S3 Bucket**: `achievement-management-terraform-state-prod`
- **State Key**: `prod/terraform.tfstate`
- **DynamoDB Table**: `terraform-state-lock-prod`
- **Region**: `us-east-1`

## Security Features

### S3 Bucket Security
- **Server-side encryption**: AES256 encryption enabled
- **Versioning**: Enabled for state file history and rollback
- **Public access blocking**: All public access blocked
- **SSL-only policy**: HTTPS connections required
- **Lifecycle management**: Old versions cleaned up after 30 days

### DynamoDB Security
- **Server-side encryption**: Enabled for all tables
- **Point-in-time recovery**: Enabled for production environment
- **Pay-per-request billing**: Cost-effective for infrequent access

## Deployment Process

### 1. Deploy Backend Infrastructure

The backend infrastructure must be deployed first using local state:

```bash
cd terraform/backend-infrastructure
terraform init
terraform plan
terraform apply
```

### 2. Initialize Main Infrastructure

After backend infrastructure is deployed, initialize the main infrastructure:

```bash
cd terraform
terraform init -backend-config=backend-dev.hcl     # For development
terraform init -backend-config=backend-staging.hcl # For staging
terraform init -backend-config=backend-prod.hcl    # For production
```

## Backend Configuration Files

### backend-dev.hcl
```hcl
bucket         = "achievement-management-terraform-state-dev"
key            = "dev/terraform.tfstate"
region         = "us-east-1"
dynamodb_table = "terraform-state-lock-dev"
encrypt        = true
```

### backend-staging.hcl
```hcl
bucket         = "achievement-management-terraform-state-staging"
key            = "staging/terraform.tfstate"
region         = "us-east-1"
dynamodb_table = "terraform-state-lock-staging"
encrypt        = true
```

### backend-prod.hcl
```hcl
bucket         = "achievement-management-terraform-state-prod"
key            = "prod/terraform.tfstate"
region         = "us-east-1"
dynamodb_table = "terraform-state-lock-prod"
encrypt        = true
```

## State Management Best Practices

### State File Organization
- Each environment has its own state file
- State files are stored in environment-specific S3 keys
- No cross-environment state dependencies

### State Locking
- DynamoDB tables prevent concurrent modifications
- Lock ID is automatically managed by Terraform
- Locks are automatically released after operations complete

### State Backup and Recovery
- S3 versioning provides automatic backup
- Previous state versions can be restored if needed
- Lifecycle policies prevent unlimited version accumulation

## Troubleshooting

### Common Issues

#### Backend Initialization Errors
```bash
# Error: bucket does not exist
# Solution: Deploy backend infrastructure first
cd terraform/backend-infrastructure
terraform apply
```

#### State Lock Issues
```bash
# Error: state is locked
# Check if another operation is running
# If stuck, force unlock (use with caution)
terraform force-unlock <lock-id>
```

#### Permission Errors
Ensure AWS credentials have permissions for:
- S3 bucket read/write access
- DynamoDB table read/write access
- S3 bucket policy management

### Recovery Procedures

#### Restore Previous State Version
```bash
# List state file versions
aws s3api list-object-versions --bucket <bucket-name> --prefix <state-key>

# Download specific version
aws s3api get-object --bucket <bucket-name> --key <state-key> --version-id <version-id> terraform.tfstate.backup

# Restore if needed (backup current state first)
```

#### Migrate State Between Backends
```bash
# Initialize with new backend
terraform init -backend-config=new-backend.hcl

# Terraform will prompt to migrate state
# Answer 'yes' to migrate existing state
```

## Monitoring and Maintenance

### CloudWatch Metrics
- Monitor S3 bucket size and request metrics
- Monitor DynamoDB table read/write capacity
- Set up alarms for unusual activity

### Cost Optimization
- S3 lifecycle policies clean up old versions
- DynamoDB uses pay-per-request billing
- Monitor costs through AWS Cost Explorer

### Regular Maintenance
- Review and clean up old state versions
- Monitor access logs for security
- Update backend configuration as needed

## Security Considerations

### Access Control
- Use IAM roles with minimal required permissions
- Implement MFA for production access
- Audit access logs regularly

### Encryption
- State files contain sensitive information
- All data encrypted at rest and in transit
- Use AWS KMS for additional encryption if needed

### Network Security
- Backend resources are region-specific
- Use VPC endpoints for enhanced security if needed
- Monitor network access patterns

## Integration with CI/CD

### Pipeline Configuration
```yaml
# Example GitHub Actions workflow
- name: Initialize Terraform
  run: terraform init -backend-config=backend-${{ env.ENVIRONMENT }}.hcl

- name: Plan Terraform
  run: terraform plan -var-file=environments/${{ env.ENVIRONMENT }}.tfvars

- name: Apply Terraform
  run: terraform apply -auto-approve -var-file=environments/${{ env.ENVIRONMENT }}.tfvars
```

### Environment Promotion
- Use same backend configuration across environments
- Promote changes through dev → staging → prod
- Maintain separate state files for isolation