# forge

`forge` is a CLI tool for generating Moldfile from Dockerfile.

TODO: GIF

## What is Moldfile?

## Usage

> [!WARNING]
> Currently it doesn't support new Docker image format.
>
> Use `forge` with Docker Engine older version than `v25.0`.

### `mold`
```shell
forge mold PATH [flags]
```

Generate Moldfile from Dockerfile.

> [!NOTE]
For version pinning of packages, it currently supports only packages installed via `apt`.

### `vdiff`
```shell
forge vdiff FILE_PATH1 FILE_PATH2 [flags]
```

Extract version differences between two files (Moldfile/Dockerfile).

### `check`
```shell
forge check FILE_PATH [flags]
```

Check if the Dockerfile/Moldfile is appropriate format to parse with forge.
If FILE_PATH is "-", the file content is read from stdin.

## Installation

### Go tools

```shell
go install github.com/tklab-group/docker-image-forge@latest
```

### Binary

Download from [release page](https://github.com/tklab-group/forge/releases/latest).