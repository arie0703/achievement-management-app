# Security Groups Module
# This module creates security groups for ALB, ECS, and VPC endpoints
# Following least privilege principles with restrictive ingress/egress rules

# Local values for common configurations
locals {
  name_prefix = "${var.app_name}-${var.environment}"
}

# Security Group for Application Load Balancer
# Allows HTTP and HTTPS traffic from the internet
resource "aws_security_group" "alb" {
  name_prefix = "${local.name_prefix}-alb-"
  description = "Security group for Application Load Balancer"
  vpc_id      = var.vpc_id

  # Ingress rules - Allow HTTP and HTTPS from internet
  ingress {
    description = "HTTP from internet"
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    description = "HTTPS from internet"
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  # Egress rules will be added via separate security group rules to avoid circular dependency

  # Allow HTTPS outbound for health checks and other services
  egress {
    description = "HTTPS outbound"
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = merge(var.tags, {
    Name = "${local.name_prefix}-alb-sg"
    ResourceType = "security-group"
    SecurityGroupType = "alb"
    Purpose = "load-balancer-access"
    NetworkTier = "public"
  })

  lifecycle {
    create_before_destroy = true
  }
}

# Security Group for ECS Services
# Allows traffic only from ALB and outbound to required services
resource "aws_security_group" "ecs" {
  name_prefix = "${local.name_prefix}-ecs-"
  description = "Security group for ECS services"
  vpc_id      = var.vpc_id

  # Ingress rules will be added via separate security group rules to avoid circular dependency

  # Egress rules - Allow outbound traffic to required services
  egress {
    description = "HTTPS for AWS services"
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  # Allow DNS resolution
  egress {
    description = "DNS resolution"
    from_port   = 53
    to_port     = 53
    protocol    = "udp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    description = "DNS resolution TCP"
    from_port   = 53
    to_port     = 53
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = merge(var.tags, {
    Name = "${local.name_prefix}-ecs-sg"
    ResourceType = "security-group"
    SecurityGroupType = "ecs"
    Purpose = "container-access"
    NetworkTier = "private"
  })

  lifecycle {
    create_before_destroy = true
  }
}

# Security Group Rules to handle ALB <-> ECS communication
# These are created separately to avoid circular dependencies

resource "aws_security_group_rule" "alb_to_ecs" {
  type                     = "ingress"
  from_port                = var.container_port
  to_port                  = var.container_port
  protocol                 = "tcp"
  source_security_group_id = aws_security_group.alb.id
  security_group_id        = aws_security_group.ecs.id
  description              = "Allow ALB to communicate with ECS services"
}

resource "aws_security_group_rule" "alb_egress_to_ecs" {
  type                     = "egress"
  from_port                = var.container_port
  to_port                  = var.container_port
  protocol                 = "tcp"
  source_security_group_id = aws_security_group.ecs.id
  security_group_id        = aws_security_group.alb.id
  description              = "Allow ALB to send traffic to ECS services"
}

# Security Group for VPC Endpoints
# Allows HTTPS traffic from private subnets for AWS service access
resource "aws_security_group" "vpc_endpoints" {
  name_prefix = "${local.name_prefix}-vpc-endpoints-"
  description = "Security group for VPC endpoints"
  vpc_id      = var.vpc_id

  # Ingress rules - Allow HTTPS from private subnets
  ingress {
    description = "HTTPS from private subnets"
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = var.private_subnet_cidrs
  }

  # Egress rules - Allow outbound HTTPS for AWS services
  egress {
    description = "HTTPS to AWS services"
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = merge(var.tags, {
    Name = "${local.name_prefix}-vpc-endpoints-sg"
    ResourceType = "security-group"
    SecurityGroupType = "vpc-endpoints"
    Purpose = "aws-service-access"
    NetworkTier = "private"
  })

  lifecycle {
    create_before_destroy = true
  }
}

# Security Group for RDS (if needed in future)
# Currently commented out as not required by current design
# resource "aws_security_group" "rds" {
#   name_prefix = "${local.name_prefix}-rds-"
#   description = "Security group for RDS database"
#   vpc_id      = var.vpc_id
#
#   ingress {
#     description     = "MySQL/Aurora from ECS"
#     from_port       = 3306
#     to_port         = 3306
#     protocol        = "tcp"
#     security_groups = [aws_security_group.ecs.id]
#   }
#
#   tags = merge(var.tags, {
#     Name = "${local.name_prefix}-rds-sg"
#     Type = "rds"
#   })
# }