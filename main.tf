resource "aws_ami" "ami" {
  name                = "testing-ami"
  virtualization_type = "hvm"
  root_device_name    = "/dev/xvda"
  tags = {
    dd_correlation_uuid = "6962c874-bd4a-4f8c-a601-288e69a6cd9e"
  }
}