# IAM Module for Achievement Management Application
# This module creates IAM roles and policies for ECS tasks and CLI operations

# Data source for current AWS account
data "aws_caller_identity" "current" {}

# Data source for current AWS region
data "aws_region" "current" {}

# ECS Task Execution Role
# This role allows ECS to pull images from ECR and write logs to CloudWatch
resource "aws_iam_role" "ecs_task_execution_role" {
  name = "${var.app_name}-${var.environment}-ecs-task-execution-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "ecs-tasks.amazonaws.com"
        }
      }
    ]
  })

  tags = merge(var.tags, {
    Name = "${var.app_name}-${var.environment}-ecs-task-execution-role"
    ResourceType = "iam-role"
    RoleType = "ecs-task-execution"
    Service = "ecs-tasks"
    Purpose = "container-execution"
  })
}

# Attach AWS managed policy for ECS task execution
resource "aws_iam_role_policy_attachment" "ecs_task_execution_role_policy" {
  role       = aws_iam_role.ecs_task_execution_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}

# Custom policy for ECR access (if using private ECR)
resource "aws_iam_role_policy" "ecs_task_execution_ecr_policy" {
  name = "${var.app_name}-${var.environment}-ecs-execution-ecr-policy"
  role = aws_iam_role.ecs_task_execution_role.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "ecr:GetAuthorizationToken",
          "ecr:BatchCheckLayerAvailability",
          "ecr:GetDownloadUrlForLayer",
          "ecr:BatchGetImage"
        ]
        Resource = "*"
      }
    ]
  })
}

# ECS Task Role for API Service
# This role allows the application to access DynamoDB tables
resource "aws_iam_role" "ecs_task_role" {
  name = "${var.app_name}-${var.environment}-ecs-task-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "ecs-tasks.amazonaws.com"
        }
      }
    ]
  })

  tags = merge(var.tags, {
    Name = "${var.app_name}-${var.environment}-ecs-task-role"
    ResourceType = "iam-role"
    RoleType = "ecs-task"
    Service = "ecs-tasks"
    Purpose = "application-access"
  })
}

# DynamoDB access policy for ECS task role
resource "aws_iam_role_policy" "ecs_task_dynamodb_policy" {
  name = "${var.app_name}-${var.environment}-ecs-task-dynamodb-policy"
  role = aws_iam_role.ecs_task_role.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "dynamodb:GetItem",
          "dynamodb:PutItem",
          "dynamodb:UpdateItem",
          "dynamodb:DeleteItem",
          "dynamodb:Query",
          "dynamodb:Scan",
          "dynamodb:BatchGetItem",
          "dynamodb:BatchWriteItem"
        ]
        Resource = [
          for table_name in var.dynamodb_table_names :
          "arn:aws:dynamodb:${data.aws_region.current.id}:${data.aws_caller_identity.current.account_id}:table/${var.app_name}-${var.environment}-${table_name}"
        ]
      },
      {
        Effect = "Allow"
        Action = [
          "dynamodb:Query",
          "dynamodb:Scan"
        ]
        Resource = [
          for table_name in var.dynamodb_table_names :
          "arn:aws:dynamodb:${data.aws_region.current.id}:${data.aws_caller_identity.current.account_id}:table/${var.app_name}-${var.environment}-${table_name}/index/*"
        ]
      }
    ]
  })
}

# CLI Task Role for Achievement Management
# This role allows CLI tasks to access DynamoDB for achievement operations
resource "aws_iam_role" "cli_achievement_task_role" {
  name = "${var.app_name}-${var.environment}-cli-achievement-task-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "ecs-tasks.amazonaws.com"
        }
      }
    ]
  })

  tags = merge(var.tags, {
    Name = "${var.app_name}-${var.environment}-cli-achievement-task-role"
    ResourceType = "iam-role"
    RoleType = "cli-task"
    Service = "ecs-tasks"
    Purpose = "cli-achievement-operations"
    CLITaskType = "achievement"
  })
}

# DynamoDB access policy for CLI achievement tasks
resource "aws_iam_role_policy" "cli_achievement_dynamodb_policy" {
  name = "${var.app_name}-${var.environment}-cli-achievement-dynamodb-policy"
  role = aws_iam_role.cli_achievement_task_role.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "dynamodb:GetItem",
          "dynamodb:PutItem",
          "dynamodb:UpdateItem",
          "dynamodb:DeleteItem",
          "dynamodb:Query",
          "dynamodb:Scan",
          "dynamodb:BatchGetItem",
          "dynamodb:BatchWriteItem"
        ]
        Resource = [
          "arn:aws:dynamodb:${data.aws_region.current.id}:${data.aws_caller_identity.current.account_id}:table/${var.app_name}-${var.environment}-achievements",
          "arn:aws:dynamodb:${data.aws_region.current.id}:${data.aws_caller_identity.current.account_id}:table/${var.app_name}-${var.environment}-current_points"
        ]
      },
      {
        Effect = "Allow"
        Action = [
          "dynamodb:Query",
          "dynamodb:Scan"
        ]
        Resource = [
          "arn:aws:dynamodb:${data.aws_region.current.id}:${data.aws_caller_identity.current.account_id}:table/${var.app_name}-${var.environment}-achievements/index/*",
          "arn:aws:dynamodb:${data.aws_region.current.id}:${data.aws_caller_identity.current.account_id}:table/${var.app_name}-${var.environment}-current_points/index/*"
        ]
      }
    ]
  })
}

# CLI Task Role for Points Management
# This role allows CLI tasks to access DynamoDB for points operations
resource "aws_iam_role" "cli_points_task_role" {
  name = "${var.app_name}-${var.environment}-cli-points-task-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "ecs-tasks.amazonaws.com"
        }
      }
    ]
  })

  tags = merge(var.tags, {
    Name = "${var.app_name}-${var.environment}-cli-points-task-role"
    ResourceType = "iam-role"
    RoleType = "cli-task"
    Service = "ecs-tasks"
    Purpose = "cli-points-operations"
    CLITaskType = "points"
  })
}

# DynamoDB access policy for CLI points tasks
resource "aws_iam_role_policy" "cli_points_dynamodb_policy" {
  name = "${var.app_name}-${var.environment}-cli-points-dynamodb-policy"
  role = aws_iam_role.cli_points_task_role.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "dynamodb:GetItem",
          "dynamodb:PutItem",
          "dynamodb:UpdateItem",
          "dynamodb:DeleteItem",
          "dynamodb:Query",
          "dynamodb:Scan",
          "dynamodb:BatchGetItem",
          "dynamodb:BatchWriteItem"
        ]
        Resource = [
          "arn:aws:dynamodb:${data.aws_region.current.id}:${data.aws_caller_identity.current.account_id}:table/${var.app_name}-${var.environment}-current_points",
          "arn:aws:dynamodb:${data.aws_region.current.id}:${data.aws_caller_identity.current.account_id}:table/${var.app_name}-${var.environment}-achievements"
        ]
      },
      {
        Effect = "Allow"
        Action = [
          "dynamodb:Query",
          "dynamodb:Scan"
        ]
        Resource = [
          "arn:aws:dynamodb:${data.aws_region.current.id}:${data.aws_caller_identity.current.account_id}:table/${var.app_name}-${var.environment}-current_points/index/*",
          "arn:aws:dynamodb:${data.aws_region.current.id}:${data.aws_caller_identity.current.account_id}:table/${var.app_name}-${var.environment}-achievements/index/*"
        ]
      }
    ]
  })
}

# CLI Task Role for Reward Management
# This role allows CLI tasks to access DynamoDB for reward operations
resource "aws_iam_role" "cli_reward_task_role" {
  name = "${var.app_name}-${var.environment}-cli-reward-task-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "ecs-tasks.amazonaws.com"
        }
      }
    ]
  })

  tags = merge(var.tags, {
    Name = "${var.app_name}-${var.environment}-cli-reward-task-role"
    ResourceType = "iam-role"
    RoleType = "cli-task"
    Service = "ecs-tasks"
    Purpose = "cli-reward-operations"
    CLITaskType = "reward"
  })
}

# DynamoDB access policy for CLI reward tasks
resource "aws_iam_role_policy" "cli_reward_dynamodb_policy" {
  name = "${var.app_name}-${var.environment}-cli-reward-dynamodb-policy"
  role = aws_iam_role.cli_reward_task_role.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "dynamodb:GetItem",
          "dynamodb:PutItem",
          "dynamodb:UpdateItem",
          "dynamodb:DeleteItem",
          "dynamodb:Query",
          "dynamodb:Scan",
          "dynamodb:BatchGetItem",
          "dynamodb:BatchWriteItem"
        ]
        Resource = [
          "arn:aws:dynamodb:${data.aws_region.current.id}:${data.aws_caller_identity.current.account_id}:table/${var.app_name}-${var.environment}-rewards",
          "arn:aws:dynamodb:${data.aws_region.current.id}:${data.aws_caller_identity.current.account_id}:table/${var.app_name}-${var.environment}-reward_history",
          "arn:aws:dynamodb:${data.aws_region.current.id}:${data.aws_caller_identity.current.account_id}:table/${var.app_name}-${var.environment}-current_points"
        ]
      },
      {
        Effect = "Allow"
        Action = [
          "dynamodb:Query",
          "dynamodb:Scan"
        ]
        Resource = [
          "arn:aws:dynamodb:${data.aws_region.current.id}:${data.aws_caller_identity.current.account_id}:table/${var.app_name}-${var.environment}-rewards/index/*",
          "arn:aws:dynamodb:${data.aws_region.current.id}:${data.aws_caller_identity.current.account_id}:table/${var.app_name}-${var.environment}-reward_history/index/*",
          "arn:aws:dynamodb:${data.aws_region.current.id}:${data.aws_caller_identity.current.account_id}:table/${var.app_name}-${var.environment}-current_points/index/*"
        ]
      }
    ]
  })
}

# CloudWatch Logs access for all CLI task roles
resource "aws_iam_role_policy" "cli_cloudwatch_logs_policy" {
  count = 3
  name  = "${var.app_name}-${var.environment}-cli-${local.cli_role_names[count.index]}-logs-policy"
  role  = local.cli_roles[count.index].id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "logs:CreateLogGroup",
          "logs:CreateLogStream",
          "logs:PutLogEvents",
          "logs:DescribeLogStreams"
        ]
        Resource = "arn:aws:logs:${data.aws_region.current.id}:${data.aws_caller_identity.current.account_id}:log-group:/ecs/${var.app_name}-${var.environment}-cli-*"
      }
    ]
  })
}

# Local values for easier management
locals {
  cli_roles = [
    aws_iam_role.cli_achievement_task_role,
    aws_iam_role.cli_points_task_role,
    aws_iam_role.cli_reward_task_role
  ]

  cli_role_names = [
    "achievement",
    "points",
    "reward"
  ]
}