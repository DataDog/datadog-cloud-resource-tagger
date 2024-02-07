# cloud-resources-tagger
Tagger for Terraform cloud resource files to include ownership details based on git commit info.

The project is based on the [Yor project by bridgecrew](https://github.com/bridgecrewio/yor-action).

## Features
* Apply tags and labels on infrastructure as code directory or list of selected files
* Change management: git-based tags automatically add org, repo, commit and modifier details on every resource block.

## How it works
Each Terraform file (based on file extension - *.tf) is parsed and processed into a set of blocks.
Each one of these blocks are later tagged with a set of attributes.
The following tags are being added for each resource configuration:
* **dd_correlation_uuid**: a tag to enable attribution between an IaC resource block and a running cloud resource
* **dd_git_org**: organization 
* **dd_git_repo**: repository
* **dd_git_file**: file Path
* **dd_git_modifiers**: users who modified the resource 
* **dd_git_created_by**: created by (user's email of the first commit)
* **dd_git_create_commit**: created at (date of the first commit)
* **dd_git_create_commit**: commit id of the first commit
* **dd_git_last_modified_by**: last modified by (user's email of the last commit)
* **dd_git_last_modified_at**: last modified at (date of the last commit)
* **dd_git_modified_commit**: commit id of last commit 
* **dd_git_resource_lines**: lines in the code matching the resource definition

## CI/CD Integration
You can use it in your CI/CD using our integration:
* [GitHub Action](https://github.com/DataDog/datadog-cloud-resource-tagger-action)
