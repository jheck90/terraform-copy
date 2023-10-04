package cmd

import (
    "github.com/spf13/cobra"
	"github.com/jheck90/terraform-copy/lib" // Import your lib package
)

var rootCmd = &cobra.Command{
    Use:   "terraform-copy",
    Short: "Copy Terraform files to specific workspaces/environments",
    Run: func(cmd *cobra.Command, args []string) {
        // Your logic for copying files goes here
        lib.GenerateAndWriteTFFiles(parentDir, sourceDir, dryRun, remoteURL, s3Bucket, region, profile, baseGitRemotePath, tfVersion)
    },
}

// Define flags
var parentDir string
var sourceDir string
var dryRun bool
var remoteURL string
var s3Bucket string
var region string
var profile string
var baseGitRemotePath string
var tfVersion string

func init() {
    // Define flags for the root command
    rootCmd.PersistentFlags().StringVar(&parentDir, "parent-dir", "", "Specify the parent directory")
    rootCmd.PersistentFlags().StringVar(&sourceDir, "source-dir", "", "Specify the source directory")
    rootCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", true, "Enable dry-run mode")
    rootCmd.PersistentFlags().StringVar(&remoteURL, "remote-url", "", "Specify the Git remote URL")
    rootCmd.PersistentFlags().StringVar(&s3Bucket, "s3-bucket", "", "Specify the S3 bucket name")
    rootCmd.PersistentFlags().StringVar(&region, "region", "us-east-1", "Specify the AWS region")
    rootCmd.PersistentFlags().StringVar(&profile, "profile", "", "Specify the AWS profile")
    rootCmd.PersistentFlags().StringVar(&tfVersion, "terraform-version", "", "Specify the Terraform version")

    // Add additional commands if needed
    // For example:
    // rootCmd.AddCommand(otherCommand)

    // Initialize your library with configuration if necessary
    // lib.Initialize(config)
}

// Execute the Cobra commands
func Execute() error {
    return rootCmd.Execute()
}
