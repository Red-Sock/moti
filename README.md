# moti

`moti` is a CLI tool for managing Protocol Buffers workflows, designed to simplify dependency management and code generation. It acts as a powerful wrapper around `protoc`, allowing you to define your protobuf dependencies and generation rules in a clean, declarative configuration file.

## Key Features

- **Dependency Management**: Easily fetch proto files from external git repositories.
- **Declarative Code Generation**: Define `protoc` commands in a simple YAML configuration.
- **Smart Caching**: Local cache for external proto dependencies to speed up generation.
- **Consistent Workflow**: Replaces complex, error-prone `protoc` commands with simple `moti install` and `moti generate`.

## Installation

```bash
go install go.redsock.ru/moti@latest
```

## Quick Start

1. Initialize a `moti.yaml` file in your project root.
2. Define your dependencies and generation rules.
3. Run `moti install` to fetch dependencies and tools.
4. Run `moti generate` to generate code.

## Configuration (moti.yaml)

The `moti.yaml` file is the heart of `moti`. Here's a breakdown based on a full example:

```yaml
# Local directory where external proto dependencies will be cached
cache_path: proto_modules

# Binary installation configuration
binaries:
  # Directory where binaries will be installed
  bin_dir: bin
  # List of go packages to install as binaries
  install:
    - go:
        # go-plugin for protoc
        module: google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.5
        # argument to binary for version checking
        version_check_args: --version
    - go:
        # go-grpc plugin for protoc
        module: google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1
        version_check_args: --version

# External git repositories to download as dependencies
deps:
  # Fetches the latest commit from the repository
  - github.com/googleapis/googleapis
  # Fetches a specific tag
  - go.redsock.ru/protoc-gen-npm@v0.0.17
  # Fetches a specific commit hash
  - github.com/vervstack/Velez@220e0db758f9ce96d9b1f457234616284530622b

# Code generation rules
generate:
  - inputs:
      # Use a local directory containing proto files
      - directory: "api"
    plugins:
      # Plugins are specified without the 'protoc-gen-' prefix
      - name: go
        out: internal/api/server
        opts:
          paths: source_relative
      - name: go-grpc
        out: internal/api/server
        opts:
          paths: source_relative
      - name: grpc-gateway
        out: internal/api/server
        opts:
          paths: source_relative

  # Generation using external repositories as inputs
  - inputs:
      - git_repo:
          url: github.com/vervstack/Matreshka
          # Optional: specify a subdirectory to use as the root for proto files
          sub_directory: api/grpc
      - git_repo:
          url: github.com/vervstack/Velez
    plugins:
      - name: go
        out: internal/api/clients
        opts:
          # Use 'module' to omit a prefix from the generated package paths
          module: go.vervstack.ru
      - name: go-grpc
        out: internal/api/clients
        opts:
          module: go.vervstack.ru
```

## Key Concepts
### Proto Root

This is **crucial**:

moti (specifically protoc) generates code from proto files based on what the first import is 

For the following project structure:

- api/
  - grpc/
    - some.proto
    - some_other.proto
    
Protoc can be called with import of the `root` of the repository (-I ".") **or** with the folder `gprc` (-I "api/grpc")

Both are valid generation ways **but** with nuances

#### For root
When one proto file imports other, it will use the full path from the root of the repository 
```protobuf common.proto
syntax = "proto3";

package moti_api;

message Empty {}
```

```protobuf api.proto
syntax = "proto3";

package moti_api;
// protoc thinks that repository is located at ./ and if common.proto passed stripped - it will look for it at the root
import "api/grpc/commmon.proto";

//...
```

#### For folder gen
When one proto file imports other, it will use the full path from the folder that 
**is considered root** (first `-I ...` import)

```protobuf common.proto
syntax = "proto3";

package moti_api;

message Empty {}
```

```protobuf api.proto
syntax = "proto3";

package moti_api;

// protoc thinks that repository is located at ./api/grpc/* and look for common.proto there
import "commmon.proto";

//...
```


#### Generation results

Based on what type of generation you prefer, you also will get different results.

For folder as root import (-I api/grpc) - generated files will be saved straight to the `out` option

So if out is `internal/api/server` and root import is the root of the repository -> 
generated files will be located there (with a full package path: go_package, java_package, etc.
unless source_relative is not passed)

For folder as root import (-I api/grpc) - generated files will be stored in `internal/api/server` 
**plus** the path to proto files

e.g. `internal/api/server/api/grpc/*.go`

## Commands

### install
Fetches and caches all dependencies defined in the `deps` section and
installs binaries defined in `binaries.install` of `moti.yaml`. 
It also manages a `moti.lock` file to ensure consistent versions across different environments.

```bash
moti install
```

### generate
Executes the code generation process based on the `generate` section in `moti.yaml`. It automatically handles the `-I` flags for `protoc`, including paths to the local project and cached dependencies.

```bash
moti generate
```

## Comparison with protoc

Instead of maintaining complex Makefiles with long `protoc` commands:

```bash
# Traditional way (hard to maintain)
protoc -I proto_modules/mod/github.com/vervstack/Matreshka/v1.0.95 \
       -I proto_modules/mod/github.com/googleapis/googleapis/v0.0.0-20260324152955-59d5f2b46924714af627ac29ea6de78641a00835 \
       --go_out=module=go.vervstack.ru/matreshka/pkg:internal/api/clients \
       api/grpc/matreshka_api.proto
```

With `moti`, you just run:

```bash
moti generate
```

## Roadmap

We are constantly working to improve `moti`. Here are some features planned for the future:

- [ ] **Adhoc Operations**: Run go mod tidy, npm run generate after or before protoc generation (e.g. useful for go-import)
  - [ ] with "AllowToFail" flag
- [x] **Binary Installment**: Automated installation of `protoc-gen-*` plugins of specific versions
  - [x] Generic like `go install ...`
  - [ ] Platform specific for `choco`, `brew`, `apt`
- [ ] **moti Registry**: A central registry for easier discovery and sharing of proto modules and watching over backward compatibility and deprecated contracts
  - [ ] Github actions for pushing to registry
- [ ] **Module option with wildcard**: to allow omitting `go.redsock.ru/*/api` patterns
- [ ] **Generation and loading animations**: cuz it looks cool and fun
- [ ] **Linting**: moved from easy-p

---
`moti` is a rewriten almost from scratch fork of [easy-p](https://github.com/thefrol/easy-p).
