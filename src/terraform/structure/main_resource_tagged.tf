resource aws_iam_role "iam_for_eks" {
  name               = "${local.resource_prefix.value}-iam-for-eks"
  assume_role_policy = data.aws_iam_policy_document.iam_policy_eks.json

  tags = {
    dd_git_resource_signature = "resource.aws_iam_role.iam_for_eks"
    dd_git_file               = "src/terraform/structure/main_resource_tagged.tf"
    dd_git_modified_commit    = "1cc401e76cb7fea95c144e804eaf06a191eb5311"
    dd_git_org                = "DataDog"
    dd_git_repo               = "datadog-cloud-resource-tagger"
  }
}