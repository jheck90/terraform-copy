// globals.go
package lib

import (
	"os"
	"path/filepath"
	"strings"
)

var prefix string
var env string
var iacValue string
var currentDir string

// InitGlobals initializes the global variables based on the current directory.
func InitGlobals() error {
	// Get the current directory
	var err error
	currentDir, err = os.Getwd()
	if err != nil {
		return err
	}

	// Derive the workspace_key_prefix from the root directory name
	rootDir := filepath.Base(currentDir)

	// Remove "-terraform" from the root directory name if it exists
	prefix = strings.TrimSuffix(rootDir, "-terraform")

	// Infer the environment or workspace from the current directory path
	parts := strings.Split(currentDir, string(filepath.Separator))
	for i, part := range parts {
		if (part == "workspaces" || part == "environments") && i < len(parts)-1 {
			env = parts[i+1]
			break
		}
	}

	// Set the iacValue based on the first part of the path
	if len(parts) > 0 {
		iacValue = parts[0]
	}

	return nil
}

// GetPrefix returns the prefix derived from the root directory.
func GetPrefix() string {
	return prefix
}

// GetEnvironment returns the inferred environment or workspace.
func GetEnvironment() string {
	return env
}

// GetIACValue returns the first part of the path.
func GetIACValue() string {
	return iacValue
}
