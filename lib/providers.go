// providers.go
package lib

import (
	"fmt"
	"path/filepath"
	"strings"
)

// GenerateProvidersFileContents generates the content for the providers.tf file.
func GenerateProvidersFileContents(profile, region, baseGitRemotePath string) string {
	iacTag := ""
	iacModulePathTag := ""

	if baseGitRemotePath != "" {
		// Generate the "IAC" tag
		iacTag = fmt.Sprintf(`      IAC              = "%s/%s"`, baseGitRemotePath, prefix)

		// Generate the "IAC:ModulePath" tag
		iacModulePathTag = fmt.Sprintf(`      "IAC:ModulePath" = "%s"`, strings.TrimPrefix(currentDir, prefix+string(filepath.Separator)))
	}

	providersContent := fmt.Sprintf(`provider "aws" {
  region  = "%s"
  profile = "%s"

  default_tags {
    tags = {
      Environment      = terraform.workspace
      Stack            = basename(abspath(path.root))
%s
%s
    }
  }
}`, region, profile, iacTag, iacModulePathTag)

	return providersContent
}
