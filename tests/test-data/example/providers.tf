provider "aws" {
  region  = "{region}"
  profile = "{profile}"

  default_tags {
    tags = {
      Environment      = terraform.workspace
      Stack            = basename(abspath(path.root))


    }
  }
}