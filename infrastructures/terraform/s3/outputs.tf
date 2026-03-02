
output "bucket_name" {
  description = "Name of the S3 bucket"
  value       = aws_s3_bucket.main.id
}

output "bucket_arn" {
  description = "ARN of the S3 bucket"
  value       = aws_s3_bucket.main.arn
}

output "bucket_region" {
  description = "Region of the S3 bucket"
  value       = aws_s3_bucket.main.region
}

output "iam_user_name" {
  description = "IAM user name for S3 access"
  value       = aws_iam_user.s3_user.name
}

output "access_key_id" {
  description = "AWS Access Key ID (save this securely)"
  value       = aws_iam_access_key.s3_user.id
}

output "secret_access_key" {
  description = "AWS Secret Access Key (save this securely, shown only once)"
  value       = aws_iam_access_key.s3_user.secret
  sensitive   = true
}
