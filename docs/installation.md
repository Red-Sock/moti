# Installation

## Prerequisites

- Go 1.21 or later
- `git` available on `$PATH`
- `protoc` available on `$PATH` (only required for `moti generate`)

## Install moti

```bash
go install go.redsock.ru/moti@latest
```

Verify the installation:

```bash
moti --version
```

## Project setup

Create a `moti.yaml` file in the root of your project. This is the only file moti needs.

A minimal example to get started:

```yaml
deps:
  - github.com/googleapis/googleapis

generate:
  - inputs:
      - directory: api/grpc
    plugins:
      - name: go
        out: internal/api
        opts:
          paths: source_relative
      - name: go-grpc
        out: internal/api
        opts:
          paths: source_relative
```

Then run:

```bash
moti install   # fetch deps and install protoc plugins
moti generate  # run protoc
```

## Adding to CI

Both commands are idempotent. `moti install` skips already-cached modules and already-installed binaries, so it is safe to run unconditionally in CI.

A typical CI step:

```yaml
- name: Generate protos
  run: |
    go install go.redsock.ru/moti@latest
    moti install
    moti generate
```

## Committing the lock file

Always commit `moti.lock` to version control. It pins the exact commit hash and content hash for every installed dependency, ensuring reproducible builds across all environments.