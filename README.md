# Terraform File Copy Tool

The Terraform File Copy Tool is a command-line utility written in Go that allows you to easily copy Terraform configuration files to specific workspaces or environments within your project.

## Table of Contents

- [Installation](#installation)
- [Usage](#usage)
- [Commands](#commands)
- [Configuration](#configuration)
- [Contributing](#contributing)
- [License](#license)

## Installation

You can install the Terraform File Copy Tool by following these steps:

1. Clone this repository to your local machine:

   ```bash
   git clone https://github.com/your-username/terraform-file-copy.git
   ```

2. Navigate to the project directory:

   ```bash
   cd terraform-file-copy
   ```

3. Build the application:

   ```bash
   go build -o terraform-copy main.go
   ```

4. Optionally, move the generated binary to a location in your system's PATH for easy access:

   ```bash
   mv terraform-copy /usr/local/bin/
   ```

## Usage

The Terraform File Copy Tool allows you to copy Terraform files to specific workspaces or environments. You can use it as follows:

```bash
terraform-copy [flags] [command]
```

For example, to copy files to a specific workspace:

```bash
terraform-copy --parent-dir /path/to/parent --example-dir example/workspace1
```

## Commands

The following commands are available:

- `copy`: Copy Terraform files to specific workspaces/environments.

## Configuration

You can configure the behavior of the Terraform File Copy Tool using command-line flags. Here are some of the available flags:

- `--parent-dir`: Specify the parent directory.
- `--example-dir`: Specify the example directory.
- `--dry-run`: Enable dry-run mode.
- `--remote-url`: Specify the Git remote URL.
- `--s3-bucket`: Specify the S3 bucket name.
- `--region`: Specify the AWS region.
- `--profile`: Specify the AWS profile.
- `--terraform-version`: Specify the Terraform version.

## Contributing

Contributions to this project are welcome! If you have ideas, bug reports, or feature requests, please open an issue on the [GitHub repository](https://github.com/your-username/terraform-file-copy). If you'd like to contribute code, please follow the standard Go practices and open a pull request.