# cloud-resources-tagger
Tagger for Terraform cloud resource files to include ownership details based on git commit info.

The project is based on the [Yor project by bridgecrew](https://github.com/bridgecrewio/yor).

## Features
* Apply tags and labels on infrastructure as code directory or list of selected files
* Change management: git-based tags automatically add org, repo, commit and file details on every resource block.

## How it works
Each Terraform file (based on file extension - *.tf) is parsed and processed into a set of blocks.
Each one of these blocks are later tagged with a set of attributes.
The following tags are being added for each resource configuration:
<u>These are the minimum set of tags we look to collect:</u>
* **dd_git_org**: organization 
* **dd_git_repo**: repository
* **dd_git_file**: filepath
* **dd_git_modified_commit**: commit id of last commit 
* **dd_git_resource_signature**: resource signature from Terraform
<u>These are the remaining set of tags it is possible for us to collect:</u>
* **dd_correlation_uuid**: a tag to enable attribution between an IaC resource block and a running cloud resource
* **dd_git_resource_lines**: lines in the code matching the resource definition
* **dd_git_repo_url**: repository url
* **dd_git_modifiers**: users who modified the resource 
* **dd_git_created_by**: created by (user's email of the first commit)
* **dd_git_create_commit**: created at (date of the first commit)
* **dd_git_create_commit**: commit id of the first commit
* **dd_git_last_modified_by**: last modified by (user's email of the last commit)
* **dd_git_last_modified_at**: last modified at (date of the last commit)

## CI/CD Integration
You can use it in your CI/CD using our integration:
* [GitHub Action](https://github.com/DataDog/datadog-cloud-resource-tagger-action)
* Gitlab: 
```
stages:
 - datadog-cloud-resource-tagger

run-datadog-cloud-resource-tagger:   
 stage: datadog-cloud-resource-tagger
 script:
   - git checkout ${CI_COMMIT_REF_NAME}
   - wget -q -O - https://github.com/DataDog/datadog-cloud-resource-tagger/releases/latest/download/datadog-cloud-resource-tagger_Linux_x86_64.tar.gz | tar -xvz -C /tmp
   - /tmp/datadog-cloud-resource-tagger tag -d <.|directory path> -t dd_git_org,dd_git_repo,dd_git_file,dd_git_modified_commit,dd_git_resource_signature
```

## Running locally
You can brew install the CLI by running the following commands:
1. `brew tap datadog/datadog-cloud-resource-tagger https://github.com/DataDog/datadog-cloud-resource-tagger`
2. `brew install datadog-cloud-resource-tagger`

You may need to run `sudo launchctl config user path "$(brew --prefix)/bin:${PATH}"` and relaunch your terminal if running on MacOS Mountain Lion or later. See [this](https://docs.brew.sh/FAQ#my-mac-apps-dont-find-homebrew-utilities) for more information.

## Command flags
The command to run when invoking the cloud resource tagger is:

`datadog-cloud-resource-tagger tag`

The following flags are available when running:
* --directory (alias -d): specify the directory to scope tagging over. By default will use `.` if no value is provided (ie tag everything)
* --tags (alias -t): specify the exact list of tags to add. By default will apply the entire list of tags specified above if no value provided. To scope to the minimum ones needed use the following argument:`-t "dd_git_org,dd_git_repo,dd_git_file,dd_git_modified_commit,dd_git_resource_signature"`
* --tag-groups (alias -g): specify the tag groups to generate tags from. By default we will use `"git,code2cloud"`.
* --changed-files: only run the tagger on the specified comma separated list of absolute filepaths
* --include-resource-types: specify the comma separated resource types to tag and skip all others ie `--include-resource-types="aws_s3_bucket,gcp_compute_instance"`
* --include-providers: specify the comma separated list of providers to tag and skip all others ie `--include-providers="aws,gcp"`
