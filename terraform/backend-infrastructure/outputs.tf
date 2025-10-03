# Outputs for Backend Infrastructure

output "s3_bucket_names" {
  description = "Names of the S3 buckets created for Terraform state storage"
  value = {
    for env, bucket in aws_s3_bucket.terraform_state : env => bucket.id
  }
}

output "s3_bucket_arns" {
  description = "ARNs of the S3 buckets created for Terraform state storage"
  value = {
    for env, bucket in aws_s3_bucket.terraform_state : env => bucket.arn
  }
}

output "dynamodb_table_names" {
  description = "Names of the DynamoDB tables created for Terraform state locking"
  value = {
    for env, table in aws_dynamodb_table.terraform_state_lock : env => table.name
  }
}

output "dynamodb_table_arns" {
  description = "ARNs of the DynamoDB tables created for Terraform state locking"
  value = {
    for env, table in aws_dynamodb_table.terraform_state_lock : env => table.arn
  }
}

output "backend_configuration_summary" {
  description = "Summary of backend configuration for each environment"
  value = {
    for env in local.environments : env => {
      bucket         = aws_s3_bucket.terraform_state[env].id
      key            = "${env}/terraform.tfstate"
      region         = var.aws_region
      dynamodb_table = aws_dynamodb_table.terraform_state_lock[env].name
      encrypt        = true
    }
  }
}