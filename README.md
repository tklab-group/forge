# forge

![GitHub Release](https://img.shields.io/github/v/release/tklab-group/forge)
![GitHub License](https://img.shields.io/github/license/tklab-group/forge)

`forge` is a CLI tool for generating Moldfile from Dockerfile.

![](./assets/demo.gif)

## What is Moldfile?

Moldfile is a Dockerfile compatible file format proposed in [CDCM](https://github.com/tklab-group/CDCM).
It strictly records versions of all dependencies in Dockerfile, e.g. base image and packages, and provides high configuration reproducibility.
Its recommended file name is `Moldfile` or `Dockerfile.mold`.

### Dockerfile
For writing Dockerfile, developers usually use loose version specification.
`latest` or a meaningful named tag for a base image and no version specification for installing packages in RUN instructions.
It provides high readability and maintainability.

```Dockerfile
FROM ubuntu:latest

RUN apt-get update && apt-get install -y \
    curl
```

### Moldfile
Moldfile is generated from Dockerfile and pins versions of all dependencies.
Digest for a base image and version pinning for packages.
It provides high configuration reproducibility for building Docker image and records the configuration at the time of its generation.

```Dockerfile
FROM ubuntu@sha256:e6173d4dc55e76b87c4af8db8821b1feae4146dd47341e4d431118c7dd060a74

RUN apt-get update && apt-get install -y \
    curl=7.81.0-1ubuntu1.15
```

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

Extract specified version differences between two files (Moldfile/Dockerfile).

It only supports comparison between the same structure instruction files.
It is expected to use for Dockerfile and Moldfile generated from the Dockerfile, or Moldfiles generated from the same Dockerfile.

<details>
<summary>Example</summary>

Moldfile1:
```Dockerfile
FROM ubuntu@sha256:ed4a42283d9943135ed87d4ee34e542f7f5ad9ecf2f244870e23122f703f91c2

RUN apt-get update && apt-get install -y \
    wget=1.20.3-1ubuntu2
```

Moldfile2:

```Dockerfile
FROM ubuntu@sha256:4c32aacd0f7d1d3a29e82bee76f892ba9bb6a63f17f9327ca0d97c3d39b9b0ee

RUN apt-get update && apt-get install -y \
    wget=1.21.3-1ubuntu1
```

Output:
```json
{
  "buildStages": [
    {
      "stageName": "",
      "baseImage": {
        "name": "ubuntu",
        "moldfile1": "@sha256:ed4a42283d9943135ed87d4ee34e542f7f5ad9ecf2f244870e23122f703f91c2",
        "moldfile2": "@sha256:4c32aacd0f7d1d3a29e82bee76f892ba9bb6a63f17f9327ca0d97c3d39b9b0ee"
      },
      "packages": [
        {
          "packageManager": "apt",
          "name": "wget",
          "moldfile1": "1.20.3-1ubuntu2",
          "moldfile2": "1.21.3-1ubuntu1"
        }
      ]
    }
  ]
}
```
</details>

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