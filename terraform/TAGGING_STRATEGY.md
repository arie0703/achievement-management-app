# Resource Tagging and Naming Strategy

## Overview

This document outlines the comprehensive tagging and naming strategy implemented for the Achievement Management Application infrastructure. The strategy ensures consistent resource organization, cost tracking, compliance, and operational efficiency across all environments.

## Naming Conventions

### Standard Naming Pattern

All resources follow the consistent naming pattern:
```
{app_name}-{environment}-{resource_type}-{identifier}
```

**Examples:**
- VPC: `achievement-management-prod-vpc`
- ECS Cluster: `achievement-management-prod-cluster`
- DynamoDB Table: `achievement-management-prod-achievements`
- Load Balancer: `achievement-management-prod-alb`

### Resource-Specific Naming Rules

#### Load Balancers
- ALB names are truncated to 32 characters due to AWS limits
- Pattern: `{app_name}-{environment}-alb` (truncated if needed)

#### Target Groups
- Pattern: `{app_name}-{environment}-{service}-tg`
- Truncated to 24 characters for AWS limits

#### DynamoDB Tables
- Pattern: `{app_name}-{environment}-{table_type}`
- Examples: `achievement-management-prod-achievements`

#### ECS Resources
- Cluster: `{app_name}-{environment}-cluster`
- Service: `{app_name}-{environment}-{service_type}-service`
- Task Definition: `{app_name}-{environment}-{service_type}`

#### Network Resources
- VPC: `{app_name}-{environment}-vpc`
- Subnets: `{app_name}-{environment}-{subnet_type}-subnet-{number}`
- Security Groups: `{app_name}-{environment}-{purpose}-sg`

## Tagging Strategy

### Core Tags (Applied to All Resources)

| Tag Key | Description | Example Values |
|---------|-------------|----------------|
| `Environment` | Environment name | `dev`, `staging`, `prod` |
| `Application` | Application name | `achievement-management` |
| `Project` | Project identifier | `achievement-management` |
| `Component` | Infrastructure component | `infrastructure` |
| `ManagedBy` | Management tool | `terraform` |
| `Owner` | Resource owner | `platform-team`, `development-team` |
| `Team` | Responsible team | `engineering` |
| `CostCenter` | Cost center for billing | `engineering-prod`, `engineering-dev` |

### Operational Tags

| Tag Key | Description | Example Values |
|---------|-------------|----------------|
| `BackupPolicy` | Backup requirements | `daily`, `weekly`, `none` |
| `MonitoringEnabled` | Monitoring status | `true`, `false` |
| `ServiceTier` | Service criticality | `critical`, `important`, `standard` |
| `BusinessUnit` | Business unit | `technology` |

### Compliance and Governance Tags

| Tag Key | Description | Example Values |
|---------|-------------|----------------|
| `DataClassification` | Data sensitivity | `public`, `internal`, `confidential`, `restricted` |
| `ComplianceScope` | Compliance requirements | `production-soc2`, `pre-production-testing` |

### Lifecycle Tags

| Tag Key | Description | Example Values |
|---------|-------------|----------------|
| `CreatedBy` | Creation method | `terraform` |
| `CreatedDate` | Creation date | `2024-01-15` |
| `LastModified` | Last modification date | `2024-01-15` |

### Resource-Specific Tags

#### Network Resources
- `ResourceType`: `vpc`, `subnet`, `internet-gateway`, `nat-gateway`, `route-table`
- `NetworkTier`: `public`, `private`, `core`
- `SubnetType`: `public`, `private`
- `AvailabilityZone`: `us-east-1a`, `us-east-1b`

#### Compute Resources (ECS)
- `ResourceType`: `ecs-cluster`, `ecs-service`, `ecs-task-definition`
- `ServiceType`: `api`, `cli`
- `LaunchType`: `fargate`
- `ContainerInsights`: `enabled`, `disabled`
- `CPU`: `256`, `512`, `1024`
- `Memory`: `512`, `1024`, `2048`

#### Database Resources (DynamoDB)
- `ResourceType`: `dynamodb-table`
- `TableType`: `achievements`, `rewards`, `current_points`, `reward_history`
- `BillingMode`: `PAY_PER_REQUEST`, `PROVISIONED`
- `BackupEnabled`: `true`, `false`
- `EncryptionEnabled`: `true`, `false`
- `AccessPattern`: `high-frequency`, `standard`
- `HasGSI`: `true`, `false`

#### Load Balancer Resources
- `ResourceType`: `application-load-balancer`, `target-group`, `lb-listener`
- `LoadBalancerType`: `application`
- `Scheme`: `internet-facing`, `internal`
- `Protocol`: `HTTP`, `HTTPS`
- `DeletionProtection`: `enabled`, `disabled`

## Environment-Specific Configuration

### Development Environment
```hcl
owner = "development-team"
service_tier = "standard"
backup_policy = "none"
data_classification = "internal"
compliance_scope = "development-only"
```

### Staging Environment
```hcl
owner = "platform-team"
service_tier = "important"
backup_policy = "daily"
data_classification = "internal"
compliance_scope = "pre-production-testing"
```

### Production Environment
```hcl
owner = "platform-team"
service_tier = "critical"
backup_policy = "daily"
data_classification = "confidential"
compliance_scope = "production-soc2"
```

## Cost Tracking and Organization

### Cost Allocation Tags
The following tags are specifically used for cost tracking and allocation:

- `CostCenter`: Maps to departmental budgets
- `Environment`: Separates costs by environment
- `Project`: Groups costs by project
- `Team`: Assigns costs to teams
- `BusinessUnit`: High-level cost allocation

### Resource Grouping
Resources are grouped using the `ResourceGroup` tag:
```
ResourceGroup = "{app_name}-{environment}"
```

This enables:
- AWS Resource Groups for operational management
- Cost and usage reports by resource group
- Automated resource discovery and management

## Compliance and Governance

### Data Classification
Resources are tagged with data classification levels:
- `public`: No restrictions
- `internal`: Internal use only
- `confidential`: Restricted access
- `restricted`: Highly sensitive data

### Compliance Scope
Resources are tagged with applicable compliance requirements:
- `development-only`: No compliance requirements
- `pre-production-testing`: Testing compliance controls
- `production-soc2`: Full SOC 2 compliance

## Automation and Tooling

### Terraform Implementation
Tags are implemented using:
- Local values for common tags
- Merge functions to combine common and resource-specific tags
- Environment-specific variable files for tag values

### Tag Validation
The following validations are implemented:
- Service tier must be one of: `critical`, `important`, `standard`
- Backup policy must be one of: `daily`, `weekly`, `monthly`, `none`
- Data classification must be one of: `public`, `internal`, `confidential`, `restricted`

### Resource Discovery
Tags enable automated resource discovery for:
- Backup automation based on `BackupPolicy` tag
- Monitoring setup based on `MonitoringEnabled` tag
- Security scanning based on `DataClassification` tag
- Cost optimization based on `ServiceTier` tag

## Best Practices

### Tag Consistency
1. Use consistent tag keys across all resources
2. Standardize tag values (avoid variations like "prod" vs "production")
3. Use lowercase for tag values where possible
4. Avoid special characters in tag values

### Tag Management
1. Define tags in Terraform variables for consistency
2. Use local values to compute dynamic tag values
3. Implement tag validation rules
4. Document tag meanings and usage

### Cost Optimization
1. Use cost allocation tags consistently
2. Review and update tags regularly
3. Use tags for automated cost reporting
4. Implement tag-based cost alerts

### Security and Compliance
1. Tag resources with data classification
2. Use tags for compliance scope identification
3. Implement tag-based access controls
4. Audit tag compliance regularly

## Monitoring and Alerting

### Tag-Based Monitoring
Resources can be monitored based on tags:
- `ServiceTier=critical`: Enhanced monitoring and alerting
- `MonitoringEnabled=true`: Include in monitoring dashboards
- `Environment=prod`: Production-level alerting thresholds

### Cost Alerts
Cost alerts are configured based on:
- `CostCenter`: Department-level budget alerts
- `Environment`: Environment-specific cost thresholds
- `Project`: Project budget monitoring

## Future Enhancements

### Planned Improvements
1. Automated tag compliance checking
2. Tag-based resource lifecycle management
3. Enhanced cost allocation reporting
4. Integration with AWS Config for tag governance

### Tag Evolution
As the infrastructure grows, additional tags may be added:
- `Version`: Application version tags
- `Schedule`: Automated start/stop scheduling
- `Backup`: Backup schedule specifications
- `Patch`: Patching schedule and requirements