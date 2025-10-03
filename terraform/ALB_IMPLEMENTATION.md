# Application Load Balancer Implementation

## Overview

This document describes the Application Load Balancer (ALB) implementation for the Achievement Management application. The ALB provides HTTP/HTTPS traffic routing, health checks, and SSL termination for the ECS-hosted Go application.

## Architecture

The ALB is implemented as part of the ECS module and includes:

- **Application Load Balancer**: Internet-facing ALB in public subnets
- **Target Group**: Routes traffic to ECS tasks with health checks
- **HTTP Listener**: Handles HTTP traffic on port 80
- **HTTPS Listener**: Optional HTTPS traffic on port 443 with SSL termination
- **Health Checks**: Monitors application health at `/health` endpoint

## Configuration

### Basic ALB Configuration

```hcl
# Application Load Balancer
resource "aws_lb" "main" {
  name               = "${substr(local.name_prefix, 0, 28)}-alb"
  internal           = false
  load_balancer_type = "application"
  security_groups    = [var.alb_security_group_id]
  subnets            = var.public_subnet_ids
  
  enable_deletion_protection = var.enable_deletion_protection
}
```

### Target Group with Health Checks

```hcl
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
}
```

### HTTP/HTTPS Listeners

The implementation supports both HTTP and HTTPS listeners:

- **HTTP Listener**: Always created, can redirect to HTTPS if configured
- **HTTPS Listener**: Optional, requires SSL certificate

## Environment-Specific Configuration

### Development Environment
- HTTP only (no HTTPS)
- No deletion protection
- Minimal health check thresholds

```hcl
alb_config = {
  enable_https               = false
  enable_https_redirect      = false
  certificate_arn           = ""
  ssl_policy                = "ELBSecurityPolicy-TLS-1-2-2017-01"
  enable_deletion_protection = false
}
```

### Staging Environment
- HTTPS enabled with HTTP redirect
- SSL certificate required
- Standard health check configuration

```hcl
alb_config = {
  enable_https               = true
  enable_https_redirect      = true
  certificate_arn           = ""  # Provide via environment variable
  ssl_policy                = "ELBSecurityPolicy-TLS-1-2-2017-01"
  enable_deletion_protection = false
}
```

### Production Environment
- HTTPS enabled with HTTP redirect
- SSL certificate required
- Deletion protection enabled
- Optimized health check settings

```hcl
alb_config = {
  enable_https               = true
  enable_https_redirect      = true
  certificate_arn           = ""  # Provide via environment variable
  ssl_policy                = "ELBSecurityPolicy-TLS-1-2-2017-01"
  enable_deletion_protection = true
}
```

## Health Check Configuration

The ALB performs health checks against the Go application's `/health` endpoint:

- **Path**: `/health`
- **Protocol**: HTTP
- **Port**: Same as container port (8080)
- **Healthy Threshold**: 2 consecutive successes
- **Unhealthy Threshold**: 3 consecutive failures
- **Timeout**: 5 seconds
- **Interval**: 30 seconds
- **Expected Response**: HTTP 200

### Go Application Health Endpoint

The Go application provides a health check endpoint:

```go
// healthCheck ヘルスチェックハンドラー
func (s *Server) healthCheck(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "status":  "ok",
        "message": "Achievement Management API is running",
    })
}
```

## Security Configuration

### Security Groups

The ALB uses dedicated security groups:

- **ALB Security Group**: Allows HTTP (80) and HTTPS (443) from internet
- **ECS Security Group**: Allows traffic from ALB security group only

### SSL/TLS Configuration

For HTTPS-enabled environments:

- **SSL Policy**: `ELBSecurityPolicy-TLS-1-2-2017-01` (default)
- **Certificate**: AWS Certificate Manager (ACM) certificate required
- **HTTP Redirect**: Optional automatic redirect from HTTP to HTTPS

## Outputs

The ALB implementation provides the following outputs:

```hcl
# Load Balancer Information
output "load_balancer_dns_name" {
  description = "DNS name of the Application Load Balancer"
  value       = module.ecs.load_balancer_dns_name
}

output "load_balancer_url" {
  description = "URL of the Application Load Balancer"
  value       = module.ecs.load_balancer_url
}

output "load_balancer_http_url" {
  description = "HTTP URL of the Application Load Balancer"
  value       = module.ecs.load_balancer_http_url
}

output "load_balancer_https_url" {
  description = "HTTPS URL of the Application Load Balancer (if enabled)"
  value       = module.ecs.load_balancer_https_url
}
```

## Deployment

### Prerequisites

1. **VPC and Subnets**: Public subnets for ALB placement
2. **Security Groups**: Configured ALB and ECS security groups
3. **SSL Certificate**: For HTTPS-enabled environments (ACM certificate)

### Deployment Steps

1. **Configure Environment Variables**:
   ```bash
   # For HTTPS environments, set certificate ARN
   export TF_VAR_certificate_arn="arn:aws:acm:region:account:certificate/cert-id"
   ```

2. **Deploy Infrastructure**:
   ```bash
   terraform init
   terraform plan -var-file="environments/dev.tfvars"
   terraform apply -var-file="environments/dev.tfvars"
   ```

3. **Verify Deployment**:
   ```bash
   # Get ALB DNS name
   terraform output load_balancer_dns_name
   
   # Test health check
   curl http://$(terraform output -raw load_balancer_dns_name)/health
   ```

## Monitoring and Troubleshooting

### CloudWatch Metrics

The ALB automatically publishes metrics to CloudWatch:

- **RequestCount**: Number of requests processed
- **TargetResponseTime**: Response time from targets
- **HTTPCode_Target_2XX_Count**: Successful responses
- **HTTPCode_Target_4XX_Count**: Client errors
- **HTTPCode_Target_5XX_Count**: Server errors
- **HealthyHostCount**: Number of healthy targets
- **UnHealthyHostCount**: Number of unhealthy targets

### Common Issues

1. **Health Check Failures**:
   - Verify `/health` endpoint is accessible
   - Check security group rules
   - Review application logs

2. **SSL Certificate Issues**:
   - Ensure certificate is in the same region
   - Verify certificate is validated
   - Check certificate ARN format

3. **Target Registration**:
   - Verify ECS service is running
   - Check target group health status
   - Review ECS task network configuration

## Best Practices

1. **SSL Certificates**: Use AWS Certificate Manager for automatic renewal
2. **Health Checks**: Keep health check endpoint lightweight
3. **Security Groups**: Follow principle of least privilege
4. **Monitoring**: Set up CloudWatch alarms for key metrics
5. **Deletion Protection**: Enable for production environments
6. **Access Logs**: Consider enabling ALB access logs for troubleshooting

## Integration with ECS

The ALB is tightly integrated with the ECS service:

- **Service Discovery**: ECS automatically registers/deregisters tasks
- **Rolling Deployments**: ALB supports blue/green deployments
- **Auto Scaling**: Works with ECS auto scaling policies
- **Health Checks**: Failed health checks trigger task replacement

This implementation provides a robust, scalable, and secure load balancing solution for the Achievement Management application.