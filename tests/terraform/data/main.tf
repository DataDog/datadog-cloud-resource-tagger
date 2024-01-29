data "aws_s3_bucket" "primary" {
  count  = var.create_bucket == true ? 0 : 1
  bucket = "externally-created-bucket"
}

resource "aws_s3_bucket" "primary" {
  count  = var.create_bucket == true ? 1 : 0
  bucket = "cloud-resource-tagger-bug-test-bucket"
}