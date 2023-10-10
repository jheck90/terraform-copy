terraform {
  backend "s3" {
    bucket               = "{s3BucketName}"
    workspace_key_prefix = "{moduleEnvPath}"
    key                  = "{moduleName}"
    region               = "{region}"
    profile              = "{profile}"
  }
}