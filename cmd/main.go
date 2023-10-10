package main

import (
	"fmt"
	"io"
	"os"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	s3BucketName string
	moduleName   string
	region		 string
	version	 	 string
	sourceDir    string
	targetDir1   string
	dryRun       bool
	debug        bool
)

func main() {
    var rootCmd = &cobra.Command{Use: "terraform-copy"}
    rootCmd.PersistentFlags().StringVarP(&s3BucketName, "s3-bucket-name", "b", "", "S3 bucket name (required)")
	rootCmd.PersistentFlags().StringVarP(&version, "tf version", "v", "> 1.0", "Version of terraform")
    rootCmd.PersistentFlags().StringVarP(&moduleName, "module-name", "m", "", "Module name")
    rootCmd.PersistentFlags().StringVarP(&region, "region", "r", "us-east-1", "AWS Region")
    rootCmd.PersistentFlags().StringVarP(&sourceDir, "source-dir", "s", "", "Source directory")
    rootCmd.PersistentFlags().StringVarP(&targetDir1, "target-dir1", "t", "", "Target directory 1")
    rootCmd.PersistentFlags().BoolVarP(&dryRun, "dry-run", "d", true, "Enable dry run: run with -d=false")
    rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Enable debug mode: --debug")

    var copyCmd = &cobra.Command{
        Use:   "copy",
        Short: "Copy files to target directory",
        Run:   copyFiles,
    }
    rootCmd.AddCommand(copyCmd)

    // Validate the required flags before running any command
    rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
        if s3BucketName == "" {
            fmt.Println("Error: S3 bucket name is required.")
            cmd.Help() // Display command help
            os.Exit(1)
        }
    }

    if err := rootCmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}


func copyFiles(cmd *cobra.Command, args []string) {
    if sourceDir == "" || targetDir1 == "" {
        color.Red("Error: Source directory and target directory must be provided.")
        return
    }

    // Validate the existence of required files in the source directory
    requiredFiles := []string{"main.tf", "backend.tf", "providers.tf", "versions.tf"}
    for _, file := range requiredFiles {
        filePath := filepath.Join(sourceDir, file)
        if _, err := os.Stat(filePath); os.IsNotExist(err) {
            color.Red("Error: Required file %s does not exist in the source directory.", file)
            return
        }
    }

    // Iterate through subdirectories under targetDir1/workspaces
    workspacePath := filepath.Join(targetDir1, "workspaces")
    subDirs, err := getSubdirectories(workspacePath, debug)
    if err != nil {
        color.Red("Error getting subdirectories: %v", err)
        return
    }

    for _, env := range subDirs {
        // Create the moduleName directory within each environment directory
        envPath := filepath.Join(workspacePath, env)
        moduleEnvPath := filepath.Join(envPath, moduleName)

        // Check if moduleName directory already exists
        _, err := os.Stat(moduleEnvPath)
        if err == nil && !dryRun {
            color.Red("Error: Directory %s already exists in the path.", moduleName)
            return
        }

        // Debug mode
        if debug {
            color.Yellow("Debug mode is enabled. Verbose output will be provided.")
            fmt.Printf("Source directory: %s\n", sourceDir)
            fmt.Printf("Target directory for %s/%s: %s\n", env, moduleName, moduleEnvPath)
        }

        // Dry run mode
        if dryRun {
            color.Green("Dry run is enabled. Rerun with -d=false when ready. Changes to be made:")
            fmt.Printf("Copy files from %s to %s\n", sourceDir, moduleEnvPath)
        } else {
            // Create the moduleName directory if it doesn't exist
            if err := createDirectory(moduleEnvPath); err != nil {
                color.Red("Error creating %s directory: %v", moduleName, err)
                return
            }

            // Implement the actual copy-paste logic here
            if err := copyDirectoryContents(sourceDir, moduleEnvPath, debug); err != nil {
                color.Red("Error copying files to %s: %v", moduleEnvPath, err)
                return
            }
            fmt.Printf("Files copied to %s from %s\n", moduleEnvPath, sourceDir)

            // Call the editCopiedFiles function with subDirs as a parameter
            err := editCopiedFiles(moduleName, moduleEnvPath, targetDir1, env, dryRun, region, debug, s3BucketName, version)
            if err != nil {
                return
            }
        }
    }
}



func getSubdirectories(path string, debug bool) ([]string, error) {
    var subDirs []string
    err := filepath.Walk(path, func(subPath string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if info.IsDir() && subPath != path {
            subDir := strings.TrimPrefix(subPath, path+"/")
            subDirs = append(subDirs, subDir)
        }
        return nil
    })
    if err != nil {
        return nil, err
    }

    if debug {
        fmt.Println("Subdirectories:")
        for _, subDir := range subDirs {
            fmt.Println(subDir)
        }
    }

    return subDirs, nil
}


func createDirectory(path string) error {
    _, err := os.Stat(path)
    if os.IsNotExist(err) {
        if err := os.MkdirAll(path, os.ModePerm); err != nil {
            return err
        }
    }
    return nil
}


func copyDirectoryContents(sourceDir, targetDir string, debug bool) error {
	return filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relativePath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return err
		}

		targetFilePath := filepath.Join(targetDir, relativePath)

		if info.IsDir() {
			// Create the corresponding directory in the target path
			if err := os.MkdirAll(targetFilePath, os.ModePerm); err != nil {
				return err
			}
		} else {
			// Copy the file
			if err := copyFile(path, targetFilePath); err != nil {
				return err
			}
		}

		if debug {
			fmt.Printf("Copied: %s -> %s\n", path, targetFilePath)
		}

		return nil
	})
}


func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	if debug {
		fmt.Printf("Copied file from %s to %s\n", src, dst)
	}

	return nil
}

func editCopiedFiles(moduleName string, moduleEnvPath string, targetDir1 string, env string, dryRun bool, region string, debug bool, s3BucketName string, version string) error {
    // Parse the root directory name from targetDir1
    targetDirName := filepath.Base(targetDir1)
    targetDirName = strings.TrimSuffix(targetDirName, "-terraform")

    // Create the profile using newDirName and targetDirName
    profile := fmt.Sprintf("%s-%s.Deploy", targetDirName, env)

    // Keep track of changes made
    changesMade := false

    // Iterate through files in moduleEnvPath
    err := filepath.Walk(moduleEnvPath, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }

        // Check if it's a file
        if !info.IsDir() {
            // Read the content of the file
            content, err := ioutil.ReadFile(path)
            if err != nil {
                return err
            }

            // Replace tokens with actual values
            editedContent := string(content)
            editedContent = strings.ReplaceAll(editedContent, "{moduleName}", moduleName)
            editedContent = strings.ReplaceAll(editedContent, "{moduleEnvPath}", targetDirName)
            editedContent = strings.ReplaceAll(editedContent, "{profile}", profile)
			editedContent = strings.ReplaceAll(editedContent, "{region}", region)
			editedContent = strings.ReplaceAll(editedContent, "{s3BucketName}", s3BucketName)
			editedContent = strings.ReplaceAll(editedContent, "{version}", version)



            // Check if any changes were made
            if string(content) != editedContent {
                changesMade = true
                if debug {
                    fmt.Printf("After creating the directory, %s, replaced token(s) in file: %s\n", moduleName, path)
                }

                // Write the edited content back to the file
                if !dryRun {
                    err = ioutil.WriteFile(path, []byte(editedContent), info.Mode())
                    if err != nil {
                        return err
                    }
                }
            }
        }
        return nil
    })

    if err != nil && !dryRun {
        color.Red("Error editing copied files: %v", err)
        return err
    }

    // Include the code to interpolate values in backend.tf
    backendTFPath := filepath.Join(moduleEnvPath, "backend.tf")
    err = interpolateBackendTF(backendTFPath, moduleName, moduleEnvPath, targetDirName, profile, region, s3BucketName, version)
    if err != nil && !dryRun {
        color.Red("Error interpolating backend.tf: %v", err)
        return err
    }

    if changesMade {
        color.Green("Files edited successfully.")
    } else {
        color.Yellow("No changes made to files.")
    }

    return nil
}

func interpolateBackendTF(filePath, moduleName, moduleEnvPath, targetDirName, profile, region, s3BucketName, version string) error {
    content, err := ioutil.ReadFile(filePath)
    if err != nil {
        return err
    }

    // Replace tokens with actual values
    editedContent := string(content)
    editedContent = strings.ReplaceAll(editedContent, "{moduleName}", moduleName)
	editedContent = strings.ReplaceAll(editedContent, "{moduleEnvPath}", targetDirName)
    editedContent = strings.ReplaceAll(editedContent, "{profile}", profile)
	editedContent = strings.ReplaceAll(editedContent, "{region}", region)
	editedContent = strings.ReplaceAll(editedContent, "{s3BucketName}", s3BucketName)
	editedContent = strings.ReplaceAll(editedContent, "{version}", version)




    // Write the edited content back to backend.tf
    err = ioutil.WriteFile(filePath, []byte(editedContent), 0644)
    if err != nil {
        return err
    }

    return nil
}
