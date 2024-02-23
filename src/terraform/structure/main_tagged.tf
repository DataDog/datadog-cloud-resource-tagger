module "complete_sg" {
  source              = "terraform-aws-modules/security-group/aws"
  name                = "my-sg-test"
  vpc_id              = "some-vpc-id"
  use_name_prefix     = true
  ingress_cidr_blocks = ["10.10.0.0/16"]
  ingress_rules       = ["https-443-tcp"]
  tags = {
    dd_git_file           = "tests/terraform/module/module/main.tf"
    dd_git_org            = "DataDog"
    dd_git_repo           = "datadog-cloud-resource-tagger"
    dd_git_resource_lines = "1:8"
  }
}