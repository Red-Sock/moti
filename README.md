# ProtoPack

`protopack` is a cli tool for workflows with `proto` files.

A fork with couple of additions from protopack

## Community
```bash
go install go.redsock.ru/protopack@latest
```

## Init

Creates empty `protopack` project.

Creates `protopack.yaml` (by default) and `protopack.lock` files.

### Usage

```bash
protopack init
```

### Usage

```bash
protopack lint -cfg example.protopack.yaml
```
## Breaking check

Checking your current API on backward compatibility with API from another branch.

### Usage

```bash
protopack breaking --against $BRANCH_TO_COMPARE_WITH
```

## Generate

Generate proto files. 

### Usage

There are several ways to get proto files to generate:
1. from current local repository:
```yaml
generate:
  - inputs:
      - directory: WHERE YOUR PROTO FILES ARE
```
2. from remote git repository:
```yaml
generate:
  - inputs:
      - git_repo:
          url: "URL TO REMOTE REPO"
          sub_directory: DIR WITH PROTO FILES ON REMOTE REPO
```
**NOTE:** format `url` the same as in `deps` section.

`plugins` section: config for `protoc`

Example:
```yaml
generate:
  - plugins:
      - name: go
        out: .
        opts:
          paths: source_relative
      - name: go-grpc
        out: .
        opts:
          paths: source_relative
          require_unimplemented_servers: false
```

## Package manager

Install dependence from `protopack` config (or lock file).