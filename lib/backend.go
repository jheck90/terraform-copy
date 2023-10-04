// backend.go
package lib

import (
	"fmt"
	// "os"
	"path/filepath"
	"strings"
)

// GenerateBackendFileContents generates the content for the backend.tf file.
func GenerateBackendFileContents(s3Bucket, sourceDir, region, profile string) string {
	// Derive the workspace_key_prefix from the root directory name
	// currentDir, err := os.Getwd()
	// if err != nil {
	// 	// Handle error
	// 	return ""
	// }
	rootDir := filepath.Base(sourceDir)

	// Remove "-terraform" from the root directory name if it exists
	prefix := strings.TrimSuffix(rootDir, "-terraform")

	// Infer the environment or workspace from the current directory path
	env := ""
	parts := strings.Split(currentDir, string(filepath.Separator))
	for i, part := range parts {
		if (part == "workspaces" || part == "environments") && i < len(parts)-1 {
			env = parts[i+1]
			break
		}
	}

	// Set the profile based on the inferred environment or workspace
	if env != "" {
		profile = fmt.Sprintf("%s-%s.Deploy", prefix, env)
	}

	backendContent := fmt.Sprintf(`terraform {
  backend "s3" {
    bucket               = "%s"
    workspace_key_prefix = "%s"
    key                  = "%s"
    region               = "%s"
    profile              = "%s"
  }
}`, s3Bucket, prefix, sourceDir, region, profile)

	return backendContent
}
