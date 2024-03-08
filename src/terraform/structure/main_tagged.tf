module "complete_sg" {
  source              = "terraform-aws-modules/security-group/aws"
  name                = "my-sg-test"
  vpc_id              = "some-vpc-id"
  use_name_prefix     = true
  ingress_cidr_blocks = ["10.10.0.0/16"]
  ingress_rules       = ["https-443-tcp"]
  tags = {
    dd_git_create_commit      = "6028228c8899b558e9e8e9c2dd29c6a05c3667a7"
  }
}