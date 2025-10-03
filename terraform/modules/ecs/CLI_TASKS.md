# ECS CLI Task Definitions

This document describes the CLI task definitions created for the Achievement Management application.

## Overview

The ECS module creates task definitions for running CLI operations in AWS ECS Fargate. These tasks allow you to perform management operations without direct server access.

## Task Types

### 1. General CLI Tasks

These are the main CLI task definitions that can be used with command overrides:

- `achievement-management-{env}-cli-achievement`: For achievement management operations
- `achievement-management-{env}-cli-points`: For points management operations  
- `achievement-management-{env}-cli-reward`: For reward management operations

### 2. Specific CLI Tasks (Optional)

When `enable_cli_specific_tasks` is set to `true`, additional specific task definitions are created:

#### Achievement Operations
- `achievement-management-{env}-cli-achievement-create`: Create new achievements
- `achievement-management-{env}-cli-achievement-list`: List all achievements

#### Points Operations
- `achievement-management-{env}-cli-points-current`: Show current point balance
- `achievement-management-{env}-cli-points-aggregate`: Show point aggregation summary

#### Reward Operations
- `achievement-management-{env}-cli-reward-create`: Create new rewards
- `achievement-management-{env}-cli-reward-list`: List all rewards
- `achievement-management-{env}-cli-reward-redeem`: Redeem rewards

## Running CLI Tasks

### Using AWS CLI

To run a CLI task using AWS CLI:

```bash
# Run a general CLI task with command override
aws ecs run-task \
  --cluster achievement-management-dev-cluster \
  --task-definition achievement-management-dev-cli-achievement \
  --launch-type FARGATE \
  --network-configuration "awsvpcConfiguration={subnets=[subnet-xxx,subnet-yyy],securityGroups=[sg-xxx],assignPublicIp=DISABLED}" \
  --overrides '{
    "containerOverrides": [
      {
        "name": "cli-achievement",
        "command": ["./achievement-app", "achievement", "create", "--title", "Test Achievement", "--description", "Test description", "--point", "10"]
      }
    ]
  }'

# Run a specific CLI task (no override needed)
aws ecs run-task \
  --cluster achievement-management-dev-cluster \
  --task-definition achievement-management-dev-cli-achievement-create \
  --launch-type FARGATE \
  --network-configuration "awsvpcConfiguration={subnets=[subnet-xxx,subnet-yyy],securityGroups=[sg-xxx],assignPublicIp=DISABLED}" \
  --overrides '{
    "containerOverrides": [
      {
        "name": "cli-achievement-create",
        "command": ["./achievement-app", "achievement", "create", "--title", "Test Achievement", "--description", "Test description", "--point", "10"]
      }
    ]
  }'
```

### Using Terraform Outputs

The module provides outputs that contain the necessary information for running tasks:

```hcl
# Get task execution information
output "cli_execution_info" {
  value = module.ecs.task_execution_info
}

# Example usage in another resource
resource "null_resource" "run_cli_task" {
  provisioner "local-exec" {
    command = <<-EOT
      aws ecs run-task \
        --cluster ${module.ecs.task_execution_info.cluster_name} \
        --task-definition ${module.ecs.task_execution_info.task_definition_arns.achievement} \
        --launch-type FARGATE \
        --network-configuration "awsvpcConfiguration={subnets=[${join(",", module.ecs.task_execution_info.subnet_ids)}],securityGroups=[${module.ecs.task_execution_info.security_group_id}],assignPublicIp=DISABLED}"
    EOT
  }
}
```

## Available CLI Commands

### Achievement Commands

```bash
# Create achievement
./achievement-app achievement create --title "Achievement Title" --description "Description" --point 10

# List achievements
./achievement-app achievement list

# Update achievement
./achievement-app achievement update --id "achievement-id" --title "New Title" --point 20

# Delete achievement
./achievement-app achievement delete --id "achievement-id"
```

### Points Commands

```bash
# Show current points
./achievement-app points current

# Show point aggregation
./achievement-app points aggregate

# Show reward history
./achievement-app points history
```

### Reward Commands

```bash
# Create reward
./achievement-app reward create --title "Reward Title" --description "Description" --point 50

# List rewards
./achievement-app reward list

# Update reward
./achievement-app reward update --id "reward-id" --title "New Title" --point 75

# Redeem reward
./achievement-app reward redeem --id "reward-id"

# Delete reward
./achievement-app reward delete --id "reward-id"
```

## IAM Permissions

Each CLI task type has its own IAM role with specific permissions:

- **Achievement tasks**: Access to achievements and current_points tables
- **Points tasks**: Access to current_points and achievements tables (read-only for aggregation)
- **Reward tasks**: Access to rewards, reward_history, and current_points tables

## Logging

All CLI tasks log to CloudWatch Logs:

- General tasks: `/ecs/achievement-management-{env}-cli-{type}`
- Specific tasks: Same log group with different stream prefixes

## Environment Variables

CLI tasks inherit the following environment variables:

- `ENVIRONMENT`: The deployment environment (dev, staging, prod)
- `AWS_REGION`: The AWS region
- `DYNAMODB_TABLE_PREFIX`: Prefix for DynamoDB table names

Additional environment variables can be configured using the `cli_environment_variables` variable.

## Configuration Variables

- `cli_task_cpu`: CPU units for CLI tasks (default: 256)
- `cli_task_memory`: Memory for CLI tasks in MB (default: 512)
- `cli_task_timeout`: Timeout for CLI tasks in seconds (default: 300)
- `enable_cli_specific_tasks`: Enable specific task definitions (default: true)
- `cli_task_types`: List of CLI task types (default: ["achievement", "points", "reward"])
- `cli_commands`: Commands for each CLI task type
- `cli_environment_variables`: Additional environment variables for CLI tasks

## Monitoring

Monitor CLI task execution through:

1. **CloudWatch Logs**: Task output and errors
2. **ECS Console**: Task status and history
3. **CloudWatch Metrics**: Task execution metrics

## Troubleshooting

Common issues and solutions:

1. **Task fails to start**: Check IAM permissions and network configuration
2. **Command not found**: Verify container image contains the CLI binary
3. **Database connection issues**: Check security groups and DynamoDB permissions
4. **Task timeout**: Increase `cli_task_timeout` or optimize CLI operations

## Security Considerations

- CLI tasks run in private subnets without public IP addresses
- Each task type has minimal required IAM permissions
- All communication with DynamoDB uses IAM roles (no hardcoded credentials)
- Logs may contain sensitive information - configure appropriate retention policies