# Resource Naming Conventions

## Overview

This document defines the standardized naming conventions used across all AWS resources in the Achievement Management Application infrastructure. Consistent naming enables better resource organization, cost tracking, and operational management.

## General Naming Pattern

All resources follow the pattern:
```
{app_name}-{environment}-{resource_type}-{identifier}
```

Where:
- `app_name`: `achievement-management`
- `environment`: `dev`, `staging`, `prod`
- `resource_type`: Descriptive resource type
- `identifier`: Optional unique identifier or sequence number

## Resource-Specific Naming Rules

### VPC and Networking Resources

| Resource Type | Naming Pattern | Example |
|---------------|----------------|---------|
| VPC | `{app_name}-{environment}-vpc` | `achievement-management-prod-vpc` |
| Internet Gateway | `{app_name}-{environment}-igw` | `achievement-management-prod-igw` |
| Public Subnet | `{app_name}-{environment}-public-subnet-{number}` | `achievement-management-prod-public-subnet-1` |
| Private Subnet | `{app_name}-{environment}-private-subnet-{number}` | `achievement-management-prod-private-subnet-1` |
| NAT Gateway | `{app_name}-{environment}-nat-gateway-{number}` | `achievement-management-prod-nat-gateway-1` |
| Elastic IP | `{app_name}-{environment}-nat-eip-{number}` | `achievement-management-prod-nat-eip-1` |
| Route Table (Public) | `{app_name}-{environment}-public-rt` | `achievement-management-prod-public-rt` |
| Route Table (Private) | `{app_name}-{environment}-private-rt-{number}` | `achievement-management-prod-private-rt-1` |

### Security Groups

| Resource Type | Naming Pattern | Example |
|---------------|----------------|---------|
| ALB Security Group | `{app_name}-{environment}-alb-sg` | `achievement-management-prod-alb-sg` |
| ECS Security Group | `{app_name}-{environment}-ecs-sg` | `achievement-management-prod-ecs-sg` |
| VPC Endpoints Security Group | `{app_name}-{environment}-vpc-endpoints-sg` | `achievement-management-prod-vpc-endpoints-sg` |

### ECS Resources

| Resource Type | Naming Pattern | Example |
|---------------|----------------|---------|
| ECS Cluster | `{app_name}-{environment}-cluster` | `achievement-management-prod-cluster` |
| API Service | `{app_name}-{environment}-api-service` | `achievement-management-prod-api-service` |
| API Task Definition | `{app_name}-{environment}-api` | `achievement-management-prod-api` |
| CLI Task Definition | `{app_name}-{environment}-cli-{task_type}` | `achievement-management-prod-cli-achievement` |
| Specific CLI Task | `{app_name}-{environment}-cli-{task_type}-{operation}` | `achievement-management-prod-cli-achievement-create` |

### Load Balancer Resources

| Resource Type | Naming Pattern | Example | Notes |
|---------------|----------------|---------|-------|
| Application Load Balancer | `{app_name}-{environment}-alb` | `achievement-management-prod-alb` | Truncated to 32 chars if needed |
| Target Group | `{app_name}-{environment}-{service}-tg` | `achievement-management-prod-api-tg` | Truncated to 24 chars if needed |
| HTTP Listener | `{app_name}-{environment}-http-listener` | `achievement-management-prod-http-listener` |
| HTTPS Listener | `{app_name}-{environment}-https-listener` | `achievement-management-prod-https-listener` |

### DynamoDB Resources

| Resource Type | Naming Pattern | Example |
|---------------|----------------|---------|
| Achievements Table | `{app_name}-{environment}-achievements` | `achievement-management-prod-achievements` |
| Rewards Table | `{app_name}-{environment}-rewards` | `achievement-management-prod-rewards` |
| Current Points Table | `{app_name}-{environment}-current_points` | `achievement-management-prod-current_points` |
| Reward History Table | `{app_name}-{environment}-reward_history` | `achievement-management-prod-reward_history` |

### IAM Resources

| Resource Type | Naming Pattern | Example |
|---------------|----------------|---------|
| ECS Task Execution Role | `{app_name}-{environment}-ecs-task-execution-role` | `achievement-management-prod-ecs-task-execution-role` |
| ECS Task Role | `{app_name}-{environment}-ecs-task-role` | `achievement-management-prod-ecs-task-role` |
| CLI Task Role | `{app_name}-{environment}-cli-{task_type}-task-role` | `achievement-management-prod-cli-achievement-task-role` |
| IAM Policy | `{app_name}-{environment}-{service}-{purpose}-policy` | `achievement-management-prod-ecs-task-dynamodb-policy` |

### CloudWatch Resources

| Resource Type | Naming Pattern | Example |
|---------------|----------------|---------|
| API Log Group | `/ecs/{app_name}-{environment}-api` | `/ecs/achievement-management-prod-api` |
| CLI Log Group | `/ecs/{app_name}-{environment}-cli-{task_type}` | `/ecs/achievement-management-prod-cli-achievement` |
| ALB Log Group | `/aws/applicationloadbalancer/{app_name}-{environment}-alb` | `/aws/applicationloadbalancer/achievement-management-prod-alb` |
| CloudWatch Alarm | `{app_name}-{environment}-{service}-{metric}-{condition}` | `achievement-management-prod-ecs-cpu-high` |
| CloudWatch Dashboard | `{app_name}-{environment}-monitoring` | `achievement-management-prod-monitoring` |

### Auto Scaling Resources

| Resource Type | Naming Pattern | Example |
|---------------|----------------|---------|
| Auto Scaling Target | `{app_name}-{environment}-{service}-autoscaling-target` | `achievement-management-prod-api-autoscaling-target` |
| Auto Scaling Policy | `{app_name}-{environment}-{service}-{metric}-scaling` | `achievement-management-prod-api-cpu-scaling` |

## AWS Resource Naming Constraints

### Character Limits
- **ALB Names**: 32 characters maximum
- **Target Group Names**: 24 characters maximum
- **Security Group Names**: 255 characters maximum
- **IAM Role Names**: 64 characters maximum
- **DynamoDB Table Names**: 255 characters maximum

### Handling Long Names
When resource names exceed AWS limits, use truncation with the following priority:
1. Keep environment and resource type
2. Truncate application name if necessary
3. Use abbreviations for common terms:
   - `achievement-management` → `achv-mgmt`
   - `application` → `app`
   - `load-balancer` → `lb`

### Example Truncation
```
Original: achievement-management-production-application-load-balancer
Truncated: achv-mgmt-prod-alb
```

## Environment-Specific Considerations

### Development Environment
- Use shorter names where possible to save costs
- Include `dev` clearly in all resource names
- Consider using shared resources where appropriate

### Staging Environment
- Mirror production naming as closely as possible
- Use `staging` or `stg` in resource names
- Maintain consistency with production for testing

### Production Environment
- Use full descriptive names
- Include `prod` clearly in all resource names
- Prioritize clarity over brevity

## Naming Validation Rules

### Required Elements
1. All resources MUST include environment identifier
2. All resources MUST include application name or abbreviation
3. All resources MUST use consistent separator (hyphen `-`)

### Prohibited Elements
1. No spaces in resource names
2. No special characters except hyphens
3. No uppercase letters (use lowercase)
4. No sequential numbers without context (use descriptive identifiers)

### Recommended Practices
1. Use descriptive suffixes (`-sg`, `-role`, `-policy`)
2. Include purpose or function in name when helpful
3. Use consistent abbreviations across all resources
4. Include availability zone information for multi-AZ resources

## Terraform Implementation

### Local Values
```hcl
locals {
  name_prefix = "${var.app_name}-${var.environment}"
  
  naming_convention = {
    vpc_name         = "${var.app_name}-${var.environment}-vpc"
    cluster_name     = "${var.app_name}-${var.environment}-cluster"
    alb_name         = "${substr("${var.app_name}-${var.environment}-alb", 0, 32)}"
    table_prefix     = "${var.app_name}-${var.environment}"
    log_group_prefix = "/aws/ecs/${var.app_name}-${var.environment}"
    role_prefix      = "${var.app_name}-${var.environment}"
    sg_prefix        = "${var.app_name}-${var.environment}"
  }
}
```

### Name Generation Functions
```hcl
# Function to generate resource names with length constraints
locals {
  generate_name = {
    alb = substr("${var.app_name}-${var.environment}-alb", 0, 32)
    target_group = substr("${var.app_name}-${var.environment}-api-tg", 0, 24)
  }
}
```

## Monitoring and Compliance

### Name Validation
- Implement Terraform validation rules for name formats
- Use consistent patterns across all modules
- Validate name lengths against AWS limits

### Documentation
- Maintain this naming convention document
- Update when new resource types are added
- Include examples for all resource types

### Automation
- Use Terraform locals for consistent name generation
- Implement name validation in CI/CD pipelines
- Generate resource inventories based on naming patterns

## Future Considerations

### Scalability
- Plan for additional environments (e.g., `test`, `demo`)
- Consider regional deployments in naming
- Plan for multi-account scenarios

### Evolution
- Allow for naming convention updates
- Maintain backward compatibility where possible
- Document any breaking changes to naming patterns