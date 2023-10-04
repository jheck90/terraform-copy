terraform {
  backend "s3" {
    bucket               = "edo-terraform"
    workspace_key_prefix = "terraform-copy"
    key                  = "tests/testdata/example/"
    region               = "us-east-1"
    profile              = ""
  }
}