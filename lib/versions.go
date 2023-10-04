// versions.go
package lib

import "fmt"

// GenerateVersionsFileContents generates the content for the versions.tf file.
func GenerateVersionsFileContents(tfVersion string) string {
	if tfVersion == "" {
		tfVersion = "~> 1.0.0" // Default Terraform version if not provided
	}

	versionsContent := fmt.Sprintf(`terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
    }
  }
  required_version = "%s"
}`, tfVersion)

	return versionsContent
}
