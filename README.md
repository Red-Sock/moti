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
3. Run `moti install` to fetch dependencies.
4. Run `moti generate` to generate code.

## Configuration (moti.yaml)

The `moti.yaml` file is the heart of `moti`. Here's a breakdown based on a full example:

```yaml
# Local directory where external proto dependencies will be cached
cache_path: proto_modules

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

## Commands

### install
Fetches and caches all dependencies defined in the `deps` section of `moti.yaml`. It also manages a `moti.lock` file to ensure consistent versions across different environments.

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
- [ ] **Binary Installment**: Automated installation of `protoc-gen-*` plugins of specific versions
  - [ ] Generic like `go install ...` on `npm install ...`
  - [ ] Platform specific for `choco`, `brew`, `apt`
- [ ] **moti Registry**: A central registry for easier discovery and sharing of proto modules and watching over backward compatibility and deprecated contracts
  - [ ] Github actions for pushing to registry
- [ ] **Module option with wildcard**: to allow omitting `go.redsock.ru/*/api` patterns
- [ ] **Generation and loading animations**: cuz it looks cool and fun
- [ ] **Linting**: moved from easy-p

---
`moti` is a rewriten almost from scratch fork of [easy-p](https://github.com/thefrol/easy-p).
