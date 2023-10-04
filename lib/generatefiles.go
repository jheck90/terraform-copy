package lib

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"io/ioutil"
	"github.com/sergi/go-diff/diffmatchpatch"
)

// GenerateAndWriteTFFiles generates Terraform configuration files and optionally outputs changes without making them.
func GenerateAndWriteTFFiles(parentDir string, sourceDir string, dryRun bool, remoteURL string, s3Bucket string, region string, profile string, baseGitRemotePath string, tfVersion string) (string, error) {
	// Generate the content for backend.tf, providers.tf, and versions.tf
	backendContent := GenerateBackendFileContents(s3Bucket, sourceDir, region, profile)
	providersContent := GenerateProvidersFileContents(profile, region, baseGitRemotePath)
	versionsContent := GenerateVersionsFileContents(tfVersion)

	// Create the output directory if it doesn't exist
	if err := os.MkdirAll(sourceDir, 0755); err != nil {
		return "", err
	}

	// Prepare variables to store changes
	var changes []string

	// Function to check for changes and append to changes slice
	checkAndAppendChange := func(filePath, newContent string) error {
		existingContent, err := readFile(filePath)
		if err != nil {
			return err
		}

		if existingContent != newContent {
			change := ComputeDiff(existingContent, newContent)
			changes = append(changes, fmt.Sprintf("Changes for %s:\n%s", filePath, change))
		}

		return nil
	}


	// Check and append changes for backend.tf
	backendFilePath := filepath.Join(sourceDir, "backend.tf")
	if err := checkAndAppendChange(backendFilePath, backendContent); err != nil {
		return "", err
	}

	// Check and append changes for providers.tf
	providersFilePath := filepath.Join(sourceDir, "providers.tf")
	if err := checkAndAppendChange(providersFilePath, providersContent); err != nil {
		return "", err
	}

	// Check and append changes for versions.tf
	versionsFilePath := filepath.Join(sourceDir, "versions.tf")
	if err := checkAndAppendChange(versionsFilePath, versionsContent); err != nil {
		return "", err
	}

	// If dryRun is true, return the changes without making any changes
	if dryRun {
		return strings.Join(changes, "\n\n"), nil
	}

	// If not a dry run, write the files
	if err := writeFile(backendFilePath, backendContent); err != nil {
		return "", err
	}
	if err := writeFile(providersFilePath, providersContent); err != nil {
		return "", err
	}
	if err := writeFile(versionsFilePath, versionsContent); err != nil {
		return "", err
	}

	return "", nil
}

// ComputeDiff computes the difference between two strings and returns it as a string.
func ComputeDiff(oldContent, newContent string) string {
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(oldContent, newContent, false)
	html := dmp.DiffPrettyHtml(diffs)
	return html
}

func writeFile(filePath, content string) error {
	// Create or open the file
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the content to the file
	_, err = file.WriteString(content)
	if err != nil {
		return err
	}

	return nil
}

// readFile reads the content of a file and returns it as a string.
func readFile(filePath string) (string, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}