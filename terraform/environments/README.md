# Environment-Specific Configuration

This directory contains environment-specific variable files for the Achievement Management application infrastructure. Each environment is configured with appropriate resource sizing, security settings, and operational parameters.

## Environment Overview

### Development Environment (`dev.tfvars`)
- **Purpose**: Cost-optimized configuration for development and testing
- **Resource Sizing**: Minimal resources (0.25 vCPU, 512 MB RAM)
- **Scaling**: 1-2 instances maximum
- **Features**: 
  - Debug logging enabled
  - No HTTPS/SSL
  - Point-in-time recovery disabled for cost savings
  - Relaxed monitoring thresholds
  - Local secret management

### Staging Environment (`staging.tfvars`)
- **Purpose**: Production-like configuration for testing and validation
- **Resource Sizing**: Moderate resources (0.5 vCPU, 1 GB RAM)
- **Scaling**: 2-5 instances
- **Features**:
  - HTTPS enabled with SSL redirect
  - Point-in-time recovery enabled
  - AWS Secrets Manager integration
  - X-Ray tracing and Container Insights enabled
  - Production-like monitoring thresholds

### Production Environment (`prod.tfvars`)
- **Resource Sizing**: High-performance resources (1 vCPU, 2 GB RAM)
- **Scaling**: 3-20 instances with high availability
- **Features**:
  - HTTPS with deletion protection
  - Provisioned DynamoDB capacity for predictable performance
  - AWS Secrets Manager with automatic rotation
  - Strict monitoring thresholds and alerting
  - Extended log retention (30 days)
  - Debug logging disabled for security

## Usage

### Deploying to an Environment

```bash
# Development environment
terraform plan -var-file="environments/dev.tfvars"
terraform apply -var-file="environments/dev.tfvars"

# Staging environment
terraform plan -var-file="environments/staging.tfvars"
terraform apply -var-file="environments/staging.tfvars"

# Production environment
terraform plan -var-file="environments/prod.tfvars"
terraform apply -var-file="environments/prod.tfvars"
```

### Using Backend Configuration

Each environment has its own backend configuration file:

```bash
# Initialize with environment-specific backend
terraform init -backend-config="backend-dev.hcl"     # For development
terraform init -backend-config="backend-staging.hcl" # For staging
terraform init -backend-config="backend-prod.hcl"    # For production
```

## Resource Sizing Guidelines

### ECS Configuration Comparison

| Environment | CPU (vCPU) | Memory (MB) | Min Instances | Max Instances | Use Case |
|-------------|------------|-------------|---------------|---------------|----------|
| Development | 0.25       | 512         | 1             | 2             | Testing, development |
| Staging     | 0.5        | 1024        | 2             | 5             | Pre-production testing |
| Production  | 1.0        | 2048        | 3             | 20            | Production workloads |

### DynamoDB Configuration

| Environment | Billing Mode | Read Capacity | Write Capacity | Point-in-Time Recovery |
|-------------|--------------|---------------|----------------|------------------------|
| Development | PAY_PER_REQUEST | N/A | N/A | Disabled |
| Staging     | PAY_PER_REQUEST | N/A | N/A | Enabled |
| Production  | PROVISIONED | 10-15 | 5-10 | Enabled |

## Security Configuration

### Secret Management

- **Development**: Uses environment variables or local configuration
- **Staging**: AWS Secrets Manager for production-like testing
- **Production**: AWS Secrets Manager with automatic rotation

### Network Security

- **Development**: HTTP only, relaxed security groups
- **Staging**: HTTPS with SSL redirect, production-like security
- **Production**: HTTPS with deletion protection, strict security groups

## Monitoring and Alerting

### Log Retention

- **Development**: 7 days
- **Staging**: 14 days  
- **Production**: 30 days

### Monitoring Thresholds

| Metric | Development | Staging | Production |
|--------|-------------|---------|------------|
| CPU Alarm | 90% | 80% | 70% |
| Memory Alarm | 90% | 80% | 75% |
| Response Time | 5.0s | 3.0s | 2.0s |
| HTTP 5xx Errors | 20 | 10 | 5 |
| DynamoDB Throttling | 5 | 1 | 0 |

## Environment-Specific Features

### Development
- Debug logging enabled
- Container Insights disabled (cost optimization)
- X-Ray tracing disabled
- Relaxed alarm thresholds

### Staging
- Debug logging enabled for troubleshooting
- Container Insights enabled
- X-Ray tracing enabled
- Production-like alarm thresholds

### Production
- Debug logging disabled for security
- Container Insights enabled for monitoring
- X-Ray tracing enabled for performance analysis
- Strict alarm thresholds with notifications

## Customization

To customize an environment configuration:

1. Copy the appropriate `.tfvars` file
2. Modify the values according to your requirements
3. Ensure resource naming follows the pattern: `{app_name}-{environment}-{resource_type}`
4. Update monitoring thresholds based on your SLA requirements
5. Configure appropriate alarm actions (SNS topics) for notifications

## Best Practices

1. **Resource Naming**: All resources include environment prefix to avoid conflicts
2. **Cost Optimization**: Development uses minimal resources and pay-per-request billing
3. **Security**: Production uses HTTPS, deletion protection, and strict monitoring
4. **Scalability**: Each environment configured for appropriate scaling patterns
5. **Monitoring**: Environment-appropriate log retention and alerting thresholds
6. **Secret Management**: Environment-specific secret management strategies