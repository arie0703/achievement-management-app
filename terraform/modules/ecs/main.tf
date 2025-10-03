# ECS Module for Achievement Management Application
# This module creates ECS cluster, services, task definitions, and load balancer

# Data sources
data "aws_caller_identity" "current" {}
data "aws_region" "current" {}

# Local values for common configurations
locals {
  name_prefix = "${var.app_name}-${var.environment}"
  
  # Common container environment variables
  common_env_vars = [
    {
      name  = "ENVIRONMENT"
      value = var.environment
    },
    {
      name  = "AWS_REGION"
      value = data.aws_region.current.name
    },
    {
      name  = "DYNAMODB_TABLE_PREFIX"
      value = "${var.app_name}-${var.environment}"
    }
  ]
}

# ECS Cluster
resource "aws_ecs_cluster" "main" {
  name = "${local.name_prefix}-cluster"

  setting {
    name  = "containerInsights"
    value = var.enable_container_insights ? "enabled" : "disabled"
  }

  tags = merge(var.tags, {
    Name = "${local.name_prefix}-cluster"
    ResourceType = "ecs-cluster"
    ContainerInsights = var.enable_container_insights ? "enabled" : "disabled"
    ComputeType = "fargate"
  })
}

# ECS Cluster Capacity Providers
resource "aws_ecs_cluster_capacity_providers" "main" {
  cluster_name = aws_ecs_cluster.main.name

  capacity_providers = ["FARGATE", "FARGATE_SPOT"]

  default_capacity_provider_strategy {
    base              = var.fargate_base_capacity
    weight            = var.fargate_weight
    capacity_provider = "FARGATE"
  }

  dynamic "default_capacity_provider_strategy" {
    for_each = var.enable_fargate_spot ? [1] : []
    content {
      base              = var.fargate_spot_base_capacity
      weight            = var.fargate_spot_weight
      capacity_provider = "FARGATE_SPOT"
    }
  }
}

# CloudWatch Log Groups are managed by the monitoring module
# Log group names are referenced in task definitions

# Application Load Balancer
resource "aws_lb" "main" {
  name               = "${substr(local.name_prefix, 0, 28)}-alb"
  internal           = false
  load_balancer_type = "application"
  security_groups    = [var.alb_security_group_id]
  subnets            = var.public_subnet_ids

  enable_deletion_protection = var.enable_deletion_protection

  tags = merge(var.tags, {
    Name = "${substr(local.name_prefix, 0, 28)}-alb"
    ResourceType = "application-load-balancer"
    LoadBalancerType = "application"
    Scheme = "internet-facing"
    DeletionProtection = var.enable_deletion_protection ? "enabled" : "disabled"
  })
}

# Target Group for API service
resource "aws_lb_target_group" "api" {
  name        = "${substr(local.name_prefix, 0, 24)}-api-tg"
  port        = var.container_port
  protocol    = "HTTP"
  vpc_id      = var.vpc_id
  target_type = "ip"

  health_check {
    enabled             = true
    healthy_threshold   = var.health_check_healthy_threshold
    interval            = var.health_check_interval
    matcher             = var.health_check_matcher
    path                = var.health_check_path
    port                = "traffic-port"
    protocol            = "HTTP"
    timeout             = var.health_check_timeout
    unhealthy_threshold = var.health_check_unhealthy_threshold
  }

  tags = merge(var.tags, {
    Name = "${substr(local.name_prefix, 0, 24)}-api-tg"
    ResourceType = "target-group"
    TargetType = "ip"
    Protocol = "HTTP"
    HealthCheckPath = var.health_check_path
    ServiceType = "api"
  })
}

# ALB Listener for HTTP
resource "aws_lb_listener" "api_http" {
  load_balancer_arn = aws_lb.main.arn
  port              = "80"
  protocol          = "HTTP"

  default_action {
    type = var.enable_https_redirect ? "redirect" : "forward"
    
    dynamic "redirect" {
      for_each = var.enable_https_redirect ? [1] : []
      content {
        port        = "443"
        protocol    = "HTTPS"
        status_code = "HTTP_301"
      }
    }
    
    dynamic "forward" {
      for_each = var.enable_https_redirect ? [] : [1]
      content {
        target_group {
          arn = aws_lb_target_group.api.arn
        }
      }
    }
  }

  tags = merge(var.tags, {
    Name = "${local.name_prefix}-http-listener"
    ResourceType = "lb-listener"
    Protocol = "HTTP"
    Port = "80"
    RedirectEnabled = var.enable_https_redirect ? "true" : "false"
  })
}

# ALB Listener for HTTPS (optional)
resource "aws_lb_listener" "api_https" {
  count = var.enable_https ? 1 : 0
  
  load_balancer_arn = aws_lb.main.arn
  port              = "443"
  protocol          = "HTTPS"
  ssl_policy        = var.ssl_policy
  certificate_arn   = var.certificate_arn

  default_action {
    type = "forward"
    forward {
      target_group {
        arn = aws_lb_target_group.api.arn
      }
    }
  }

  tags = merge(var.tags, {
    Name = "${local.name_prefix}-https-listener"
    ResourceType = "lb-listener"
    Protocol = "HTTPS"
    Port = "443"
    SSLPolicy = var.ssl_policy
  })
}

# Task Definition for API Service
resource "aws_ecs_task_definition" "api" {
  family                   = "${local.name_prefix}-api"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = var.api_task_cpu
  memory                   = var.api_task_memory
  execution_role_arn       = var.ecs_task_execution_role_arn
  task_role_arn           = var.ecs_task_role_arn

  container_definitions = jsonencode([
    {
      name  = "api"
      image = var.api_container_image
      
      portMappings = [
        {
          containerPort = var.container_port
          protocol      = "tcp"
        }
      ]

      environment = concat(local.common_env_vars, var.api_environment_variables)

      logConfiguration = {
        logDriver = "awslogs"
        options = {
          "awslogs-group"         = "/ecs/${local.name_prefix}-api"
          "awslogs-region"        = data.aws_region.current.name
          "awslogs-stream-prefix" = "ecs"
        }
      }

      healthCheck = {
        command = [
          "CMD-SHELL",
          "curl -f http://localhost:${var.container_port}${var.health_check_path} || exit 1"
        ]
        interval    = 30
        timeout     = 5
        retries     = 3
        startPeriod = 60
      }

      essential = true
    }
  ])

  tags = merge(var.tags, {
    Name = "${local.name_prefix}-api-task-definition"
    ResourceType = "ecs-task-definition"
    ServiceType = "api"
    LaunchType = "fargate"
    CPU = tostring(var.api_task_cpu)
    Memory = tostring(var.api_task_memory)
    NetworkMode = "awsvpc"
  })
}

# ECS Service for API
resource "aws_ecs_service" "api" {
  name            = "${local.name_prefix}-api-service"
  cluster         = aws_ecs_cluster.main.id
  task_definition = aws_ecs_task_definition.api.arn
  desired_count   = var.api_desired_count
  launch_type     = "FARGATE"

  network_configuration {
    security_groups  = [var.ecs_security_group_id]
    subnets          = var.private_subnet_ids
    assign_public_ip = false
  }

  load_balancer {
    target_group_arn = aws_lb_target_group.api.arn
    container_name   = "api"
    container_port   = var.container_port
  }

  # Deployment configuration - will be added in a future update
  # deployment_configuration {
  #   maximum_percent         = var.deployment_maximum_percent
  #   minimum_healthy_percent = var.deployment_minimum_healthy_percent
  # }

  depends_on = [
    aws_lb_listener.api_http,
    aws_lb_listener.api_https
  ]

  tags = merge(var.tags, {
    Name = "${local.name_prefix}-api-service"
    ResourceType = "ecs-service"
    ServiceType = "api"
    LaunchType = "fargate"
    DesiredCount = tostring(var.api_desired_count)
    LoadBalancerEnabled = "true"
  })
}

# Auto Scaling Target for API Service
resource "aws_appautoscaling_target" "api" {
  max_capacity       = var.api_max_capacity
  min_capacity       = var.api_min_capacity
  resource_id        = "service/${aws_ecs_cluster.main.name}/${aws_ecs_service.api.name}"
  scalable_dimension = "ecs:service:DesiredCount"
  service_namespace  = "ecs"

  tags = merge(var.tags, {
    Name = "${local.name_prefix}-api-autoscaling-target"
    ResourceType = "autoscaling-target"
    ServiceType = "api"
    MinCapacity = tostring(var.api_min_capacity)
    MaxCapacity = tostring(var.api_max_capacity)
  })
}

# Auto Scaling Policy - CPU
resource "aws_appautoscaling_policy" "api_cpu" {
  name               = "${local.name_prefix}-api-cpu-scaling"
  policy_type        = "TargetTrackingScaling"
  resource_id        = aws_appautoscaling_target.api.resource_id
  scalable_dimension = aws_appautoscaling_target.api.scalable_dimension
  service_namespace  = aws_appautoscaling_target.api.service_namespace

  target_tracking_scaling_policy_configuration {
    predefined_metric_specification {
      predefined_metric_type = "ECSServiceAverageCPUUtilization"
    }
    target_value = var.cpu_target_value
  }
}

# Auto Scaling Policy - Memory
resource "aws_appautoscaling_policy" "api_memory" {
  name               = "${local.name_prefix}-api-memory-scaling"
  policy_type        = "TargetTrackingScaling"
  resource_id        = aws_appautoscaling_target.api.resource_id
  scalable_dimension = aws_appautoscaling_target.api.scalable_dimension
  service_namespace  = aws_appautoscaling_target.api.service_namespace

  target_tracking_scaling_policy_configuration {
    predefined_metric_specification {
      predefined_metric_type = "ECSServiceAverageMemoryUtilization"
    }
    target_value = var.memory_target_value
  }
}

# Task Definitions for CLI operations
resource "aws_ecs_task_definition" "cli" {
  for_each = toset(var.cli_task_types)
  
  family                   = "${local.name_prefix}-cli-${each.key}"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = var.cli_task_cpu
  memory                   = var.cli_task_memory
  execution_role_arn       = var.ecs_task_execution_role_arn
  task_role_arn           = var.cli_task_role_arns[each.key]

  container_definitions = jsonencode([
    {
      name  = "cli-${each.key}"
      image = var.cli_container_image
      
      # Override the command to run CLI operations
      # This allows the task to be run with specific subcommands
      command = var.cli_commands[each.key]

      environment = concat(local.common_env_vars, var.cli_environment_variables)

      logConfiguration = {
        logDriver = "awslogs"
        options = {
          "awslogs-group"         = "/ecs/${local.name_prefix}-cli-${each.key}"
          "awslogs-region"        = data.aws_region.current.name
          "awslogs-stream-prefix" = "ecs"
        }
      }

      # CLI tasks don't need health checks as they are one-time operations
      essential = true
      
      # Set working directory if needed
      workingDirectory = "/app"
    }
  ])

  tags = merge(var.tags, {
    Name = "${local.name_prefix}-cli-${each.key}-task-definition"
    ResourceType = "ecs-task-definition"
    ServiceType = "cli"
    CLITaskType = each.key
    LaunchType = "fargate"
    CPU = tostring(var.cli_task_cpu)
    Memory = tostring(var.cli_task_memory)
    NetworkMode = "awsvpc"
  })
}

# Task Definitions for specific CLI operations with subcommands
# Achievement Management Task Definitions
resource "aws_ecs_task_definition" "cli_achievement_create" {
  count = var.enable_cli_specific_tasks ? 1 : 0
  
  family                   = "${local.name_prefix}-cli-achievement-create"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = var.cli_task_cpu
  memory                   = var.cli_task_memory
  execution_role_arn       = var.ecs_task_execution_role_arn
  task_role_arn           = var.cli_task_role_arns["achievement"]

  container_definitions = jsonencode([
    {
      name  = "cli-achievement-create"
      image = var.cli_container_image
      
      # Command for creating achievements
      command = ["./achievement-app", "achievement", "create"]

      environment = concat(local.common_env_vars, var.cli_environment_variables)

      logConfiguration = {
        logDriver = "awslogs"
        options = {
          "awslogs-group"         = "/ecs/${local.name_prefix}-cli-achievement"
          "awslogs-region"        = data.aws_region.current.name
          "awslogs-stream-prefix" = "ecs-create"
        }
      }

      essential = true
      workingDirectory = "/app"
    }
  ])

  tags = merge(var.tags, {
    Name = "${local.name_prefix}-cli-achievement-create-task-definition"
    ResourceType = "ecs-task-definition"
    ServiceType = "cli"
    CLITaskType = "achievement"
    CLIOperation = "create"
    LaunchType = "fargate"
    CPU = tostring(var.cli_task_cpu)
    Memory = tostring(var.cli_task_memory)
    NetworkMode = "awsvpc"
  })
}

resource "aws_ecs_task_definition" "cli_achievement_list" {
  count = var.enable_cli_specific_tasks ? 1 : 0
  
  family                   = "${local.name_prefix}-cli-achievement-list"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = var.cli_task_cpu
  memory                   = var.cli_task_memory
  execution_role_arn       = var.ecs_task_execution_role_arn
  task_role_arn           = var.cli_task_role_arns["achievement"]

  container_definitions = jsonencode([
    {
      name  = "cli-achievement-list"
      image = var.cli_container_image
      
      # Command for listing achievements
      command = ["./achievement-app", "achievement", "list"]

      environment = concat(local.common_env_vars, var.cli_environment_variables)

      logConfiguration = {
        logDriver = "awslogs"
        options = {
          "awslogs-group"         = "/ecs/${local.name_prefix}-cli-achievement"
          "awslogs-region"        = data.aws_region.current.name
          "awslogs-stream-prefix" = "ecs-list"
        }
      }

      essential = true
      workingDirectory = "/app"
    }
  ])

  tags = merge(var.tags, {
    Name = "${local.name_prefix}-cli-achievement-list-task-definition"
    ResourceType = "ecs-task-definition"
    ServiceType = "cli"
    CLITaskType = "achievement"
    CLIOperation = "list"
    LaunchType = "fargate"
    CPU = tostring(var.cli_task_cpu)
    Memory = tostring(var.cli_task_memory)
    NetworkMode = "awsvpc"
  })
}

# Points Management Task Definitions
resource "aws_ecs_task_definition" "cli_points_current" {
  count = var.enable_cli_specific_tasks ? 1 : 0
  
  family                   = "${local.name_prefix}-cli-points-current"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = var.cli_task_cpu
  memory                   = var.cli_task_memory
  execution_role_arn       = var.ecs_task_execution_role_arn
  task_role_arn           = var.cli_task_role_arns["points"]

  container_definitions = jsonencode([
    {
      name  = "cli-points-current"
      image = var.cli_container_image
      
      # Command for checking current points
      command = ["./achievement-app", "points", "current"]

      environment = concat(local.common_env_vars, var.cli_environment_variables)

      logConfiguration = {
        logDriver = "awslogs"
        options = {
          "awslogs-group"         = "/ecs/${local.name_prefix}-cli-points"
          "awslogs-region"        = data.aws_region.current.name
          "awslogs-stream-prefix" = "ecs-current"
        }
      }

      essential = true
      workingDirectory = "/app"
    }
  ])

  tags = var.tags
}

resource "aws_ecs_task_definition" "cli_points_aggregate" {
  count = var.enable_cli_specific_tasks ? 1 : 0
  
  family                   = "${local.name_prefix}-cli-points-aggregate"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = var.cli_task_cpu
  memory                   = var.cli_task_memory
  execution_role_arn       = var.ecs_task_execution_role_arn
  task_role_arn           = var.cli_task_role_arns["points"]

  container_definitions = jsonencode([
    {
      name  = "cli-points-aggregate"
      image = var.cli_container_image
      
      # Command for aggregating points
      command = ["./achievement-app", "points", "aggregate"]

      environment = concat(local.common_env_vars, var.cli_environment_variables)

      logConfiguration = {
        logDriver = "awslogs"
        options = {
          "awslogs-group"         = "/ecs/${local.name_prefix}-cli-points"
          "awslogs-region"        = data.aws_region.current.name
          "awslogs-stream-prefix" = "ecs-aggregate"
        }
      }

      essential = true
      workingDirectory = "/app"
    }
  ])

  tags = var.tags
}

# Reward Management Task Definitions
resource "aws_ecs_task_definition" "cli_reward_create" {
  count = var.enable_cli_specific_tasks ? 1 : 0
  
  family                   = "${local.name_prefix}-cli-reward-create"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = var.cli_task_cpu
  memory                   = var.cli_task_memory
  execution_role_arn       = var.ecs_task_execution_role_arn
  task_role_arn           = var.cli_task_role_arns["reward"]

  container_definitions = jsonencode([
    {
      name  = "cli-reward-create"
      image = var.cli_container_image
      
      # Command for creating rewards
      command = ["./achievement-app", "reward", "create"]

      environment = concat(local.common_env_vars, var.cli_environment_variables)

      logConfiguration = {
        logDriver = "awslogs"
        options = {
          "awslogs-group"         = "/ecs/${local.name_prefix}-cli-reward"
          "awslogs-region"        = data.aws_region.current.name
          "awslogs-stream-prefix" = "ecs-create"
        }
      }

      essential = true
      workingDirectory = "/app"
    }
  ])

  tags = var.tags
}

resource "aws_ecs_task_definition" "cli_reward_list" {
  count = var.enable_cli_specific_tasks ? 1 : 0
  
  family                   = "${local.name_prefix}-cli-reward-list"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = var.cli_task_cpu
  memory                   = var.cli_task_memory
  execution_role_arn       = var.ecs_task_execution_role_arn
  task_role_arn           = var.cli_task_role_arns["reward"]

  container_definitions = jsonencode([
    {
      name  = "cli-reward-list"
      image = var.cli_container_image
      
      # Command for listing rewards
      command = ["./achievement-app", "reward", "list"]

      environment = concat(local.common_env_vars, var.cli_environment_variables)

      logConfiguration = {
        logDriver = "awslogs"
        options = {
          "awslogs-group"         = "/ecs/${local.name_prefix}-cli-reward"
          "awslogs-region"        = data.aws_region.current.name
          "awslogs-stream-prefix" = "ecs-list"
        }
      }

      essential = true
      workingDirectory = "/app"
    }
  ])

  tags = var.tags
}

resource "aws_ecs_task_definition" "cli_reward_redeem" {
  count = var.enable_cli_specific_tasks ? 1 : 0
  
  family                   = "${local.name_prefix}-cli-reward-redeem"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = var.cli_task_cpu
  memory                   = var.cli_task_memory
  execution_role_arn       = var.ecs_task_execution_role_arn
  task_role_arn           = var.cli_task_role_arns["reward"]

  container_definitions = jsonencode([
    {
      name  = "cli-reward-redeem"
      image = var.cli_container_image
      
      # Command for redeeming rewards
      command = ["./achievement-app", "reward", "redeem"]

      environment = concat(local.common_env_vars, var.cli_environment_variables)

      logConfiguration = {
        logDriver = "awslogs"
        options = {
          "awslogs-group"         = "/ecs/${local.name_prefix}-cli-reward"
          "awslogs-region"        = data.aws_region.current.name
          "awslogs-stream-prefix" = "ecs-redeem"
        }
      }

      essential = true
      workingDirectory = "/app"
    }
  ])

  tags = var.tags
}