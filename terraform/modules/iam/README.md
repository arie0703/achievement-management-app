# IAM Module

This module creates IAM roles and policies for the Achievement Management application running on AWS ECS.

## Resources Created

### ECS Task Execution Role
- **Purpose**: Allows ECS to pull container images from ECR and write logs to CloudWatch
- **Permissions**: 
  - ECR image pull permissions
  - CloudWatch Logs write permissions
  - ECS task execution permissions

### ECS Task Role (API Service)
- **Purpose**: Allows the API service to access DynamoDB tables
- **Permissions**:
  - Full DynamoDB access to all application tables
  - Access to table indexes

### CLI Task Roles
Three separate roles for different CLI operations:

#### Achievement Management Role
- **Purpose**: CLI tasks for achievement operations
- **Permissions**: 
  - DynamoDB access to achievements and current-points tables
  - CloudWatch Logs write permissions

#### Points Management Role
- **Purpose**: CLI tasks for points operations
- **Permissions**:
  - DynamoDB access to current-points and achievements tables
  - CloudWatch Logs write permissions

#### Reward Management Role
- **Purpose**: CLI tasks for reward operations
- **Permissions**:
  - DynamoDB access to rewards, reward-history, and current-points tables
  - CloudWatch Logs write permissions

## Usage

```hcl
module "iam" {
  source = "./modules/iam"

  app_name    = var.app_name
  environment = var.environment
  
  dynamodb_table_names = [
    "achievements",
    "rewards", 
    "current-points",
    "reward-history"
  ]

  tags = {
    Environment = var.environment
    Application = var.app_name
  }
}
```

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| app_name | Name of the application | `string` | n/a | yes |
| environment | Environment name (dev, staging, prod) | `string` | n/a | yes |
| dynamodb_table_names | List of DynamoDB table names | `list(string)` | `["achievements", "rewards", "current-points", "reward-history"]` | no |
| tags | Tags to apply to all resources | `map(string)` | `{}` | no |

## Outputs

| Name | Description |
|------|-------------|
| ecs_task_execution_role_arn | ARN of the ECS task execution role |
| ecs_task_execution_role_name | Name of the ECS task execution role |
| ecs_task_role_arn | ARN of the ECS task role for API service |
| ecs_task_role_name | Name of the ECS task role for API service |
| cli_achievement_task_role_arn | ARN of the CLI achievement task role |
| cli_points_task_role_arn | ARN of the CLI points task role |
| cli_reward_task_role_arn | ARN of the CLI reward task role |
| cli_task_role_arns | Map of all CLI task role ARNs |

## Security Considerations

- All roles follow the principle of least privilege
- DynamoDB permissions are scoped to specific tables and indexes
- CloudWatch Logs permissions are scoped to application-specific log groups
- Each CLI operation has its own dedicated role with minimal required permissions
- Roles can only be assumed by ECS tasks