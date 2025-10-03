# ECS Module

This module creates an ECS Fargate cluster with an API service and CLI task definitions for the Achievement Management application.

## Features

- **ECS Fargate Cluster**: Serverless container execution environment
- **API Service**: Long-running HTTP API service with auto-scaling
- **CLI Task Definitions**: On-demand task definitions for command execution
- **Application Load Balancer**: HTTP load balancer with health checks
- **Auto Scaling**: CPU and memory-based auto scaling policies
- **CloudWatch Logging**: Centralized logging for all tasks
- **Security**: Network isolation using security groups

## Architecture

```
Internet → ALB → ECS Service (API) → DynamoDB
                ↓
            CLI Tasks → DynamoDB
```

## Resources Created

### Core ECS Resources
- ECS Fargate Cluster with Container Insights
- ECS Service for API with auto-scaling
- Task definitions for API and CLI operations
- CloudWatch log groups for all tasks

### Load Balancing
- Application Load Balancer (ALB)
- Target group with health checks
- HTTP listener (port 80)

### Auto Scaling
- Application Auto Scaling target
- CPU utilization scaling policy
- Memory utilization scaling policy

## Usage

```hcl
module "ecs" {
  source = "./modules/ecs"

  # Basic Configuration
  environment = var.environment
  app_name    = var.app_name
  tags        = local.common_tags

  # Network Configuration
  vpc_id                = module.vpc.vpc_id
  public_subnet_ids     = module.vpc.public_subnet_ids
  private_subnet_ids    = module.vpc.private_subnet_ids
  alb_security_group_id = module.security.alb_security_group_id
  ecs_security_group_id = module.security.ecs_security_group_id

  # IAM Configuration
  ecs_task_execution_role_arn = module.iam.ecs_task_execution_role_arn
  ecs_task_role_arn          = module.iam.ecs_task_role_arn
  cli_task_role_arns         = module.iam.cli_task_role_arns

  # Container Configuration
  api_container_image = var.api_container_image
  cli_container_image = var.cli_container_image
  container_port      = 8080

  # Scaling Configuration
  api_desired_count = 2
  api_min_capacity  = 1
  api_max_capacity  = 10
}
```

## Task Definitions

### API Service Task
- **Purpose**: Long-running HTTP API server
- **CPU**: 512 units (0.5 vCPU)
- **Memory**: 1024 MB
- **Port**: 8080
- **Health Check**: `/health` endpoint
- **Auto Scaling**: Based on CPU and memory utilization

### CLI Task Definitions
Three separate task definitions for different CLI operations:

1. **Achievement Tasks** (`achievement`)
   - Command: `./achievement-cli achievement`
   - IAM Role: CLI Achievement Task Role

2. **Points Tasks** (`points`)
   - Command: `./achievement-cli points`
   - IAM Role: CLI Points Task Role

3. **Reward Tasks** (`reward`)
   - Command: `./achievement-cli reward`
   - IAM Role: CLI Reward Task Role

## Auto Scaling Configuration

The API service includes auto scaling based on:
- **CPU Utilization**: Target 70%
- **Memory Utilization**: Target 80%
- **Min Capacity**: 1 task
- **Max Capacity**: 10 tasks

## Health Checks

### ALB Health Check
- **Path**: `/health`
- **Interval**: 30 seconds
- **Timeout**: 5 seconds
- **Healthy Threshold**: 2 consecutive successes
- **Unhealthy Threshold**: 3 consecutive failures

### Container Health Check
- **Command**: `curl -f http://localhost:8080/health`
- **Interval**: 30 seconds
- **Timeout**: 5 seconds
- **Retries**: 3
- **Start Period**: 60 seconds

## Logging

All tasks log to CloudWatch with the following log groups:
- API Service: `/ecs/{app_name}-{environment}-api`
- CLI Tasks: `/ecs/{app_name}-{environment}-cli-{task_type}`

Log retention is configurable (default: 7 days).

## Running CLI Tasks

To run a CLI task, use the AWS CLI or SDK:

```bash
aws ecs run-task \
  --cluster {cluster_name} \
  --task-definition {task_definition_arn} \
  --network-configuration "awsvpcConfiguration={subnets=[subnet-xxx],securityGroups=[sg-xxx],assignPublicIp=DISABLED}" \
  --launch-type FARGATE
```

## Variables

| Name | Description | Type | Default |
|------|-------------|------|---------|
| `environment` | Environment name | `string` | - |
| `app_name` | Application name | `string` | - |
| `vpc_id` | VPC ID | `string` | - |
| `public_subnet_ids` | Public subnet IDs for ALB | `list(string)` | - |
| `private_subnet_ids` | Private subnet IDs for ECS | `list(string)` | - |
| `api_container_image` | API container image | `string` | - |
| `cli_container_image` | CLI container image | `string` | - |
| `api_desired_count` | Desired API task count | `number` | `2` |
| `api_min_capacity` | Min API task count | `number` | `1` |
| `api_max_capacity` | Max API task count | `number` | `10` |

## Outputs

| Name | Description |
|------|-------------|
| `cluster_name` | ECS cluster name |
| `load_balancer_dns_name` | ALB DNS name |
| `load_balancer_url` | ALB URL |
| `api_service_name` | API service name |
| `cli_task_definition_arns` | CLI task definition ARNs |

## Security Considerations

- All tasks run in private subnets
- Security groups restrict network access
- IAM roles follow least privilege principle
- Container images should be scanned for vulnerabilities
- ALB only accepts HTTP traffic (HTTPS should be configured for production)

## Cost Optimization

- Uses Fargate for serverless execution
- Auto scaling prevents over-provisioning
- CLI tasks run on-demand only
- Log retention limits storage costs
- Optional Fargate Spot for cost savings