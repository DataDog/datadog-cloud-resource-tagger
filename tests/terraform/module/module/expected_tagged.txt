module "complete_sg" {
  source              = "terraform-aws-modules/security-group/aws"
  name                = "my-sg-test"
  vpc_id              = "some-vpc-id"
  use_name_prefix     = true
  ingress_cidr_blocks = ["10.10.0.0/16"]
  ingress_rules       = ["https-443-tcp"]
  tags = {
    dd_git_create_commit    = "6028228c8899b558e9e8e9c2dd29c6a05c3667a7"
    dd_git_created_at       = "2024-01-29 20:11:37"
    dd_git_created_by       = "bahar.shah@datadoghq.com"
    dd_git_file             = "tests/terraform/module/module/main.tf"
    dd_git_last_modified_at = "2024-01-29 20:11:37"
    dd_git_last_modified_by = "bahar.shah@datadoghq.com"
    dd_git_modified_commit  = "6028228c8899b558e9e8e9c2dd29c6a05c3667a7"
    dd_git_modifiers        = "bahar.shah"
    dd_git_org              = "DataDog"
    dd_git_repo             = "datadog-cloud-resource-tagger"
    dd_git_repo_url         = "git@github.com:DataDog/datadog-cloud-resource-tagger"
    dd_git_resource_lines   = "1:8"
  }
}