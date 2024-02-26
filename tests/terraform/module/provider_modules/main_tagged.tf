module "project-factory" {
  source  = "terraform-google-modules/project-factory/google"
  version = "11.0.0"
  labels = {
    test                    = "true"
    dd_correlation_uuid     = "09168bc2-c331-4e2e-a235-4f7465b6d73d"
    dd_git_create_commit    = "6028228c8899b558e9e8e9c2dd29c6a05c3667a7"
    dd_git_created_at       = "2024-01-29 20:11:37"
    dd_git_created_by       = "bahar.shah@datadoghq.com"
    dd_git_file             = "tests__terraform__module__provider_modules__main_tf"
    dd_git_last_modified_at = "2024-01-29-20-11-37"
    dd_git_last_modified_by = "baharshah"
    dd_git_modified_commit  = "6028228c8899b558e9e8e9c2dd29c6a05c3667a7"
    dd_git_modifiers        = "baharshah"
    dd_git_org              = "DataDog"
    dd_git_repo             = "datadog-cloud-resource-tagger"
    dd_git_repo_url         = "git@github_com:DataDog__datadog-cloud-resource-tagger"
    dd_git_resource_lines   = "1:7"
  }
}

module "vpc" {
  source  = "terraform-aws-modules/vpc/aws"
  version = "3.2.0"
  tags = {
    test                    = "true"
    dd_correlation_uuid     = "8dc18d06-defd-4ba3-89c9-74da3dbd3493"
    dd_git_create_commit    = "6028228c8899b558e9e8e9c2dd29c6a05c3667a7"
    dd_git_created_at       = "2024-01-29 20:11:37"
    dd_git_created_by       = "bahar.shah@datadoghq.com"
    dd_git_file             = "tests/terraform/module/provider_modules/main.tf"
    dd_git_last_modified_at = "2024-01-29 20:11:37"
    dd_git_last_modified_by = "bahar.shah@datadoghq.com"
    dd_git_modified_commit  = "6028228c8899b558e9e8e9c2dd29c6a05c3667a7"
    dd_git_modifiers        = "bahar.shah"
    dd_git_org              = "DataDog"
    dd_git_repo             = "datadog-cloud-resource-tagger"
    dd_git_repo_url         = "git@github.com:DataDog/datadog-cloud-resource-tagger"
    dd_git_resource_lines   = "9:15"
  }
}

module "project-factory_github" {
  source = "github.com/terraform-google-modules/terraform-google-project-factory"
  labels = {
    test                    = "true"
    dd_correlation_uuid     = "4e6c65da-3fa9-45f8-b077-3c2c29ae9cbd"
    dd_git_create_commit    = "6028228c8899b558e9e8e9c2dd29c6a05c3667a7"
    dd_git_created_at       = "2024-01-29 20:11:37"
    dd_git_created_by       = "bahar.shah@datadoghq.com"
    dd_git_file             = "tests__terraform__module__provider_modules__main_tf"
    dd_git_last_modified_at = "2024-01-29-20-11-37"
    dd_git_last_modified_by = "baharshah"
    dd_git_modified_commit  = "6028228c8899b558e9e8e9c2dd29c6a05c3667a7"
    dd_git_modifiers        = "baharshah"
    dd_git_org              = "DataDog"
    dd_git_repo             = "datadog-cloud-resource-tagger"
    dd_git_repo_url         = "git@github_com:DataDog__datadog-cloud-resource-tagger"
    dd_git_resource_lines   = "17:22"
  }
}

module "project-factory_git" {
  source = "git@github.com:terraform-google-modules/terraform-google-project-factory.git"
  labels = {
    test                    = "true"
    dd_correlation_uuid     = "df3838b7-3cec-491e-b170-e211e7966012"
    dd_git_create_commit    = "6028228c8899b558e9e8e9c2dd29c6a05c3667a7"
    dd_git_created_at       = "2024-01-29 20:11:37"
    dd_git_created_by       = "bahar.shah@datadoghq.com"
    dd_git_file             = "tests__terraform__module__provider_modules__main_tf"
    dd_git_last_modified_at = "2024-01-29-20-11-37"
    dd_git_last_modified_by = "baharshah"
    dd_git_modified_commit  = "6028228c8899b558e9e8e9c2dd29c6a05c3667a7"
    dd_git_modifiers        = "baharshah"
    dd_git_org              = "DataDog"
    dd_git_repo             = "datadog-cloud-resource-tagger"
    dd_git_repo_url         = "git@github_com:DataDog__datadog-cloud-resource-tagger"
    dd_git_resource_lines   = "24:29"
  }
}

module "caf" {
  source = "aztfmod/caf/azurerm"
  tags = {
    test                    = "true"
    dd_correlation_uuid     = "59748c52-d4f1-46aa-afc9-196aeea59e66"
    dd_git_create_commit    = "6028228c8899b558e9e8e9c2dd29c6a05c3667a7"
    dd_git_created_at       = "2024-01-29 20:11:37"
    dd_git_created_by       = "bahar.shah@datadoghq.com"
    dd_git_file             = "tests/terraform/module/provider_modules/main.tf"
    dd_git_last_modified_at = "2024-01-29 20:11:37"
    dd_git_last_modified_by = "bahar.shah@datadoghq.com"
    dd_git_modified_commit  = "6028228c8899b558e9e8e9c2dd29c6a05c3667a7"
    dd_git_modifiers        = "bahar.shah"
    dd_git_org              = "DataDog"
    dd_git_repo             = "datadog-cloud-resource-tagger"
    dd_git_repo_url         = "git@github.com:DataDog/datadog-cloud-resource-tagger"
    dd_git_resource_lines   = "38:43"
  }
}

module "caf" {
  source = "git@github.com:aztfmod/terraform-azurerm-caf.git"
  tags = {
    test                    = "true"
    dd_correlation_uuid     = "59748c52-d4f1-46aa-afc9-196aeea59e66"
    dd_git_create_commit    = "6028228c8899b558e9e8e9c2dd29c6a05c3667a7"
    dd_git_created_at       = "2024-01-29 20:11:37"
    dd_git_created_by       = "bahar.shah@datadoghq.com"
    dd_git_file             = "tests/terraform/module/provider_modules/main.tf"
    dd_git_last_modified_at = "2024-01-29 20:11:37"
    dd_git_last_modified_by = "bahar.shah@datadoghq.com"
    dd_git_modified_commit  = "6028228c8899b558e9e8e9c2dd29c6a05c3667a7"
    dd_git_modifiers        = "bahar.shah"
    dd_git_org              = "DataDog"
    dd_git_repo             = "datadog-cloud-resource-tagger"
    dd_git_repo_url         = "git@github.com:DataDog/datadog-cloud-resource-tagger"
    dd_git_resource_lines   = "38:43"
  }
}

module "bastion" {
  source = "oracle-terraform-modules/bastion/oci"
  freeform_tags = {
    test                    = "true"
    dd_correlation_uuid     = "9ef6f575-95df-4560-b7dc-7398274179e6"
    dd_git_create_commit    = "6028228c8899b558e9e8e9c2dd29c6a05c3667a7"
    dd_git_created_at       = "2024-01-29 20:11:37"
    dd_git_created_by       = "bahar.shah@datadoghq.com"
    dd_git_file             = "tests/terraform/module/provider_modules/main.tf"
    dd_git_last_modified_at = "2024-01-29 20:11:37"
    dd_git_last_modified_by = "bahar.shah@datadoghq.com"
    dd_git_modified_commit  = "6028228c8899b558e9e8e9c2dd29c6a05c3667a7"
    dd_git_modifiers        = "bahar.shah"
    dd_git_org              = "DataDog"
    dd_git_repo             = "datadog-cloud-resource-tagger"
    dd_git_repo_url         = "git@github.com:DataDog/datadog-cloud-resource-tagger"
    dd_git_resource_lines   = "45:50"
  }
}

module "run-common_logs" {
  // Tags attribute is extra_tags
  source  = "claranet/run-common/azurerm//modules/logs"
  version = "3.0.0"
  extra_tags = {
    test                    = "true"
    dd_correlation_uuid     = "dcd93d2b-5ebe-4b19-a8d4-e73333e44885"
    dd_git_create_commit    = "6028228c8899b558e9e8e9c2dd29c6a05c3667a7"
    dd_git_created_at       = "2024-01-29 20:11:37"
    dd_git_created_by       = "bahar.shah@datadoghq.com"
    dd_git_file             = "tests/terraform/module/provider_modules/main.tf"
    dd_git_last_modified_at = "2024-01-29 20:11:37"
    dd_git_last_modified_by = "bahar.shah@datadoghq.com"
    dd_git_modified_commit  = "6028228c8899b558e9e8e9c2dd29c6a05c3667a7"
    dd_git_modifiers        = "bahar.shah"
    dd_git_org              = "DataDog"
    dd_git_repo             = "datadog-cloud-resource-tagger"
    dd_git_repo_url         = "git@github.com:DataDog/datadog-cloud-resource-tagger"
    dd_git_resource_lines   = "52:59"
  }
}