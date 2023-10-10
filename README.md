# Terraform Copy Tool

A command-line utility for copying Terraform module files and customizing their content.

## Usage

```shell
terraform-copy copy [flags]
```

### Flags

- `-b`, `--s3-bucket-name`: S3 bucket name (required).
- `-v`, `--tf-version`: Version of Terraform (default: "> 1.0").
- `-m`, `--module-name`: Module name.
- `-r`, `--region`: AWS Region (default: "us-east-1").
- `-s`, `--source-dir`: Source directory (required).
- `-t`, `--target-dir1`: Target directory 1 (required).
- `-d`, `--dry-run`: Enable dry run (default: true).
- `--debug`: Enable debug mode.

### Copy Files

```shell
terraform-copy copy -b <S3_BUCKET> -m <MODULE_NAME> -s <SOURCE_DIR> -t <TARGET_DIR>
```

### Edit Files

```shell
terraform-copy edit -m <MODULE_NAME>
```

## Installation

1. Clone the repository.

2. Build the application.

3. Move the binary to your system's `PATH`.

## Contributing

Contributions welcome! Open an issue or create a pull request on [GitHub](https://github.com/jheck90/terraform-copy).

## To-Do List

- [ ] Needs to loop over an AWS Organization and list all accounts
    - [ ] Filter out accounts
- [ ] [terraform-exec](https://github.com/hashicorp/terraform-exec)
    - [ ] Terraform plan and apply
