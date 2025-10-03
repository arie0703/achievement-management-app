# DynamoDB Module

This module creates DynamoDB tables for the achievement management application with proper configuration for data access patterns, security, and backup.

## Features

- Creates DynamoDB tables based on configuration
- Supports both PAY_PER_REQUEST and PROVISIONED billing modes
- Enables point-in-time recovery for data backup
- Configures server-side encryption for data security
- Implements deletion protection for production environments
- Provides comprehensive resource tagging
- Creates Global Secondary Index for achievements table

## Tables Created

### 1. Achievements Table
- **Partition Key**: `user_id` (String)
- **Sort Key**: `achievement_id` (String)
- **Purpose**: Store user achievements with user-specific access patterns
- **GSI**: Achievement ID index for querying by achievement

### 2. Rewards Table
- **Partition Key**: `reward_id` (String)
- **Purpose**: Store available rewards and their configurations

### 3. Current Points Table
- **Partition Key**: `user_id` (String)
- **Purpose**: Store current point balances for each user

### 4. Reward History Table
- **Partition Key**: `user_id` (String)
- **Sort Key**: `timestamp` (String)
- **Purpose**: Store historical record of reward redemptions

## Usage

```hcl
module "dynamodb" {
  source = "./modules/dynamodb"

  environment = "dev"
  app_name    = "achievement-management"

  dynamodb_tables = {
    achievements = {
      hash_key               = "user_id"
      range_key              = "achievement_id"
      billing_mode           = "PAY_PER_REQUEST"
      point_in_time_recovery = true
      server_side_encryption = true
    }
    rewards = {
      hash_key               = "reward_id"
      billing_mode           = "PAY_PER_REQUEST"
      point_in_time_recovery = true
      server_side_encryption = true
    }
    current_points = {
      hash_key               = "user_id"
      billing_mode           = "PAY_PER_REQUEST"
      point_in_time_recovery = true
      server_side_encryption = true
    }
    reward_history = {
      hash_key               = "user_id"
      range_key              = "timestamp"
      billing_mode           = "PAY_PER_REQUEST"
      point_in_time_recovery = true
      server_side_encryption = true
    }
  }

  tags = {
    Environment = "dev"
    Application = "achievement-management"
    ManagedBy   = "terraform"
  }
}
```

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| environment | Environment name (dev, staging, prod) | `string` | n/a | yes |
| app_name | Application name used for resource naming | `string` | n/a | yes |
| dynamodb_tables | DynamoDB table configurations | `map(object)` | n/a | yes |
| tags | Common tags to apply to all resources | `map(string)` | `{}` | no |

## Outputs

| Name | Description |
|------|-------------|
| table_names | Map of table logical names to actual table names |
| table_arns | Map of table logical names to table ARNs |
| achievements_table_name | Name of the achievements table |
| achievements_table_arn | ARN of the achievements table |
| rewards_table_name | Name of the rewards table |
| rewards_table_arn | ARN of the rewards table |
| current_points_table_name | Name of the current points table |
| current_points_table_arn | ARN of the current points table |
| reward_history_table_name | Name of the reward history table |
| reward_history_table_arn | ARN of the reward history table |
| achievements_gsi_table_name | Name of the achievements GSI table |
| achievements_gsi_table_arn | ARN of the achievements GSI table |

## Security Features

- **Server-side encryption**: All tables are encrypted at rest using AWS managed keys
- **Point-in-time recovery**: Enabled for all tables to allow data recovery
- **Deletion protection**: Enabled for production environments to prevent accidental deletion
- **IAM integration**: Tables are designed to work with least-privilege IAM policies

## Data Access Patterns

The table design supports the following access patterns:

1. **User-specific achievements**: Query achievements by user_id
2. **Achievement lookup**: Query specific achievement by achievement_id (via GSI)
3. **Reward catalog**: Query rewards by reward_id
4. **User points**: Get current points for a user
5. **Reward history**: Query user's reward history with time-based sorting

## Cost Optimization

- Uses PAY_PER_REQUEST billing mode by default for cost efficiency
- Supports PROVISIONED mode for predictable workloads
- Implements appropriate retention and backup policies