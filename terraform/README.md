# Achievement Management Infrastructure

This directory contains Terraform configuration for deploying the Achievement Management application to AWS.

## Project Structure

```
terraform/
├── main.tf                 # Root module configuration
├── variables.tf            # Input variables
├── outputs.tf             # Output values
├── terraform.tf           # Provider and backend configuration
├── README.md              # This file
├── BACKEND_CONFIGURATION.md # Backend setup documentation
├── MODULE_INTEGRATION.md   # Module integration guide
├── environments/          # Environment-specific variable files
│   ├── dev.tfvars         # Development environment
│   ├── staging.tfvars     # Staging environment
│   └── prod.tfvars        # Production environment
├── backend-*.hcl          # Backend configuration files
│   ├── backend-dev.hcl    # Development backend config
│   ├── backend-staging.hcl # Staging backend config
│   └── backend-prod.hcl   # Production backend config
├── backend-infrastructure/ # Backend infrastructure setup
│   ├── main.tf            # S3 buckets and DynamoDB tables
│   ├── variables.tf       # Backend variables
│   ├── outputs.tf         # Backend outputs
│   └── README.md          # Backend setup instructions
├── scripts/               # Utility scripts
│   └── run-cli-task.sh    # Run ECS CLI tasks
└── modules/               # Terraform modules
    ├── vpc/               # VPC and networking resources
    ├── ecs/               # ECS cluster, services, and tasks
    ├── dynamodb/          # DynamoDB tables
    ├── iam/               # IAM roles and policies
    ├── security/          # Security groups
    └── monitoring/        # CloudWatch and logging
```

## Prerequisites

1. **AWS CLI** configured with appropriate credentials
2. **Terraform** >= 1.0 installed
3. **S3 buckets** for state storage (one per environment)
4. **DynamoDB tables** for state locking (one per environment)

## Backend Setup

Before using this Terraform configuration, you need to create the backend infrastructure for state management. This includes S3 buckets for state storage and DynamoDB tables for state locking.

### Backend Infrastructure Deployment

1. **Deploy Backend Infrastructure First**:
   ```bash
   cd backend-infrastructure
   terraform init
   terraform plan
   terraform apply
   ```

2. **Initialize Main Infrastructure with Remote Backend**:
   ```bash
   cd ..
   terraform init -backend-config=backend-dev.hcl     # For development
   terraform init -backend-config=backend-staging.hcl # For staging
   terraform init -backend-config=backend-prod.hcl    # For production
   ```

The backend infrastructure creates:
- S3 buckets with versioning and encryption for each environment
- DynamoDB tables for state locking
- Proper security configurations and lifecycle policies

For detailed setup instructions, see [BACKEND_CONFIGURATION.md](BACKEND_CONFIGURATION.md).

## Module Integration

This Terraform configuration uses a modular architecture with six main modules that work together to create the complete infrastructure. For detailed information about module integration, see [MODULE_INTEGRATION.md](MODULE_INTEGRATION.md).

### Module Architecture

- **VPC Module**: Foundational networking infrastructure
- **Security Groups Module**: Network security rules and policies
- **IAM Module**: Roles and policies for ECS tasks and CLI operations
- **DynamoDB Module**: Application data storage tables
- **ECS Module**: Container orchestration with ALB
- **Monitoring Module**: CloudWatch logging, dashboards, and alarms

### Integration Validation

After deployment, you can verify the module integration status:

```bash
# Check overall integration status
terraform output module_integration_status

# Verify deployment readiness
terraform output deployment_readiness

# Review cross-module references
```

## Usage

### Initialize Terraform

Choose the appropriate backend configuration for your target environment:

```bash
# Development
terraform init -backend-config=backend-dev.hcl

# Staging
terraform init -backend-config=backend-staging.hcl

# Production
terraform init -backend-config=backend-prod.hcl
```

### Plan and Apply

```bash
# Development
terraform plan -var-file=environments/dev.tfvars
terraform apply -var-file=environments/dev.tfvars

# Staging
terraform plan -var-file=environments/staging.tfvars
terraform apply -var-file=environments/staging.tfvars

# Production
terraform plan -var-file=environments/prod.tfvars
terraform apply -var-file=environments/prod.tfvars
```

### Destroy Infrastructure

```bash
# Development
terraform destroy -var-file=environments/dev.tfvars

# Staging
terraform destroy -var-file=environments/staging.tfvars

# Production
terraform destroy -var-file=environments/prod.tfvars
```

## Environment Configurations

### Development
- Minimal resource allocation for cost optimization
- Single AZ deployment where possible
- Shorter log retention periods
- Point-in-time recovery disabled for DynamoDB

### Staging
- Production-like configuration for testing
- Multi-AZ deployment
- Moderate resource allocation
- Point-in-time recovery enabled

### Production
- High availability configuration
- Multi-AZ deployment with redundancy
- Optimized resource allocation
- Extended log retention
- Point-in-time recovery enabled
- Provisioned capacity for DynamoDB

## Security Considerations

- All resources are deployed in private subnets where possible
- Security groups follow least privilege principle
- IAM roles use minimal required permissions
- DynamoDB tables have server-side encryption enabled
- CloudWatch logs are encrypted

## Cost Optimization

- Development environment uses minimal resources
- On-demand billing for variable workloads
- Appropriate instance sizing per environment
- Log retention policies to manage storage costs

## Monitoring and Logging

- CloudWatch log groups for all ECS tasks
- Environment-specific log retention policies
- Resource tagging for cost tracking
- CloudWatch metrics and alarms (configured in monitoring module)

## Next Steps

After setting up the basic structure, implement the individual modules:

1. VPC module for networking infrastructure
2. Security module for security groups
3. IAM module for roles and policies
4. DynamoDB module for data storage
5. Monitoring module for logging and metrics
6. ECS module for container orchestration