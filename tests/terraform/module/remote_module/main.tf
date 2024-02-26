# Defines the main configuration
locals {
  account_name = "security-sandbox"
}
module "dd-terraform-statefile" {
  source = "git::https://github.com/DataDog/cloud-inventory.git/terraform-modules/aws-bucket?ref=aws-bucket_v4.2.1"

  bucket_name = "dd-terraform-${var.account}"
  team        = "cloud-security"
  repository  = "cloud-inventory"
  tags = {
    environment         = var.account
    team                = "cloud-security"
    "terraform.managed" = "true"
  }
  local_policy = data.aws_iam_policy_document.default_policy.json
  lifecycle_rules = [
    {
      noncurrent_version_expiration = {
        noncurrent_days = 60
      },
      noncurrent_version_transition = {
        noncurrent_days = 30
        storage_class   = "STANDARD_IA"
      }
    }
  ]

}

data "aws_iam_policy_document" "default_policy" {
  statement {
    sid    = "DefaultPolicy"
    effect = "Deny"
    principals {
      type        = "*"
      identifiers = ["*"]
    }
    actions = ["*"]
    resources = [
      "arn:aws:s3:::dd-terraform-${var.account}/*"
    ]
    condition {
      test     = "StringNotLike"
      variable = "aws:arn"
      values = [
        "arn:aws:iam::${var.account-number}:root",
        "arn:aws:iam::${var.account-number}:user/*",
        "arn:aws:sts::${var.account-number}:assumed-role/*"
      ]
    }
  }

}

resource "aws_dynamodb_table" "terraform_statelock" {
  name           = "terraform-${var.account}-lock"
  read_capacity  = 20
  write_capacity = 20
  hash_key       = "LockID"

  attribute {
    name = "LockID"
    type = "S"
  }
}