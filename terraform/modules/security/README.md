# Security Groups Module

This module creates security groups for the Achievement Management application infrastructure, following AWS security best practices and the principle of least privilege.

## Overview

The module creates three main security groups:

1. **ALB Security Group** - Controls traffic to the Application Load Balancer
2. **ECS Security Group** - Controls traffic to ECS services
3. **VPC Endpoints Security Group** - Controls traffic to VPC endpoints for AWS services

## Security Groups Created

### ALB Security Group
- **Purpose**: Controls inbound traffic to the Application Load Balancer
- **Ingress Rules**:
  - HTTP (port 80) from internet (0.0.0.0/0)
  - HTTPS (port 443) from internet (0.0.0.0/0)
- **Egress Rules**:
  - HTTP (port 8080) to ECS security group
  - HTTPS (port 443) to internet for health checks

### ECS Security Group
- **Purpose**: Controls traffic to ECS services running in private subnets
- **Ingress Rules**:
  - HTTP (port 8080) from ALB security group only
- **Egress Rules**:
  - HTTPS (port 443) to internet for AWS service access
  - DNS (port 53) UDP and TCP for name resolution

### VPC Endpoints Security Group
- **Purpose**: Controls traffic to VPC endpoints for AWS services
- **Ingress Rules**:
  - HTTPS (port 443) from private subnet CIDR blocks
- **Egress Rules**:
  - HTTPS (port 443) to internet for AWS service communication

## Security Principles

### Least Privilege Access
- Each security group only allows the minimum required access
- No unnecessary ports or protocols are opened
- Source and destination restrictions are applied where possible

### Defense in Depth
- Multiple layers of security controls
- Network segmentation between public and private resources
- Explicit allow rules with implicit deny

### Zero Trust Network
- No implicit trust between components
- All traffic is explicitly allowed or denied
- Security groups reference each other to create secure communication paths

## Usage

```hcl
module "security" {
  source = "./modules/security"
  
  environment = var.environment
  app_name    = var.app_name
  vpc_id      = module.vpc.vpc_id
  
  private_subnet_cidrs = var.private_subnet_cidrs
  container_port       = var.ecs_config.container_port
  
  tags = local.common_tags
}
```

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| environment | Environment name (dev, staging, prod) | `string` | n/a | yes |
| app_name | Application name used for resource naming | `string` | n/a | yes |
| vpc_id | VPC ID where security groups will be created | `string` | n/a | yes |
| private_subnet_cidrs | CIDR blocks for private subnets | `list(string)` | `["10.0.10.0/24", "10.0.20.0/24"]` | no |
| allowed_cidr_blocks | CIDR blocks allowed to access the ALB | `list(string)` | `["0.0.0.0/0"]` | no |
| container_port | Port on which the container application runs | `number` | `8080` | no |
| tags | Tags to apply to all security group resources | `map(string)` | `{}` | no |

## Outputs

| Name | Description |
|------|-------------|
| alb_security_group_id | ID of the Application Load Balancer security group |
| ecs_security_group_id | ID of the ECS services security group |
| vpc_endpoints_security_group_id | ID of the VPC endpoints security group |
| alb_security_group_arn | ARN of the Application Load Balancer security group |
| ecs_security_group_arn | ARN of the ECS services security group |
| vpc_endpoints_security_group_arn | ARN of the VPC endpoints security group |
| security_group_names | Map of security group names |

## Dependencies

This module depends on:
- VPC module (for vpc_id)
- The security groups have circular dependencies that are handled by Terraform's dependency resolution

## Security Considerations

1. **Internet Access**: Only the ALB security group allows inbound traffic from the internet
2. **Service Communication**: ECS services can only receive traffic from the ALB
3. **AWS Service Access**: ECS services can access AWS services through VPC endpoints or internet gateway
4. **DNS Resolution**: ECS services can perform DNS lookups for service discovery
5. **Monitoring**: All security group changes should be monitored and logged

## Compliance

This module helps meet the following compliance requirements:
- **Requirement 4.4**: Network security with minimum privilege access
- **Requirement 5.4**: Service-to-service authentication and authorization