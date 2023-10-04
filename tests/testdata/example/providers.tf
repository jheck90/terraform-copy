provider "aws" {
  region  = "us-east-1"
  profile = ""

  default_tags {
    tags = {
      Environment      = terraform.workspace
      Stack            = basename(abspath(path.root))


    }
  }
}