# Configuration reference

All configuration lives in `moti.yaml` (default path). Pass a custom path with `--cfg <path>`.

---

## Top-level fields

| Field        | Type       | Default          | Description                                        |
|--------------|------------|------------------|----------------------------------------------------|
| `cache_path` | string     | `proto_modules`  | Local directory where proto dependencies are cached |
| `deps`       | []string   | —                | Remote proto repositories to fetch                 |
| `replace`    | []Replace  | —                | Override a remote module with a local path          |
| `binaries`   | Binaries   | —                | Protoc plugin management                            |
| `generate`   | []Generate | —                | Code generation rules                               |

---

## `deps`

A list of git repository paths to download proto files from. Each entry is a module path with an optional `@version` suffix.

```yaml
deps:
  # latest commit
  - github.com/googleapis/googleapis

  # specific tag
  - go.redsock.ru/protoc-gen-npm@v0.0.17

  # specific commit hash
  - github.com/vervstack/Velez@220e0db758f9ce96d9b1f457234616284530622b
```

Fetched files are stored under `<cache_path>/mod/<module>/<version>/`. The resolved commit hash and a content hash are written to `moti.lock`.

---

## `replace`

Redirect a remote module to a local directory. Useful when iterating on a proto library alongside the consumer without pushing changes first.

```yaml
replace:
  - old: github.com/vervstack/Matreshka
    new: ../Matreshka
```

Relative paths in `new` are resolved relative to the directory containing `moti.yaml`. When a module is replaced, `moti install` skips fetching it entirely, and `moti generate` reads proto files from the local path.

---

## `binaries`

Controls installation and lookup of `protoc-gen-*` binaries.

```yaml
binaries:
  bin_dir: bin
  allow_custom: false
  install:
    - go:
        module: google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.11
        version_check_args: --version
```

### `bin_dir`

A directory (relative to the project root) where binaries are installed and searched first. When set, moti prepends this directory to `PATH` before running `protoc`, so project-local binaries take precedence over globally installed ones.

Leave empty to use the default `GOBIN` / `$PATH`.

### `allow_custom`

By default, only binaries whose name starts with `protoc-gen-` may be installed. Set `allow_custom: true` to allow installing any binary (required for OpenAPI generators, for example).

### `install[].go`

Install a binary via `go install`.

| Field               | Description                                                                                  |
|---------------------|----------------------------------------------------------------------------------------------|
| `module`            | Full Go module path with optional `@version` (e.g. `google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.11`) |
| `version_check_args`| Arguments passed to the binary to check its version (e.g. `--version`). The output is checked for the version string. If empty, moti only verifies the binary exists. |

`moti install` skips installation if the binary is already present at the expected version.

---

## `generate`

A list of generation rules. Each rule has a set of `inputs` and `plugins`. Rules are executed in order; each `input` within a rule produces one `protoc` invocation.

```yaml
generate:
  - inputs:
      - directory: api/grpc
    plugins:
      - name: go
        out: internal/api/server
        opts:
          paths: source_relative
```

### `inputs`

Each input is either a local directory or a git repository.

#### Local directory

```yaml
inputs:
  - directory: api/grpc
```

`directory` is the path to walk for `.proto` files. It is used as the `-I` import root for that `protoc` invocation. Set to `.` or omit to walk the whole project (the `cache_path` directory is excluded automatically).

#### Git repository

```yaml
inputs:
  - git_repo:
      url: github.com/vervstack/Matreshka
      sub_directory: api/grpc
      out: internal/api/clients  # unused field, kept for clarity
```

The module must already be present in `deps` and installed. `sub_directory` narrows the root to a specific path inside the cached repository. If omitted, the repository root is used.

All `deps` modules are automatically added as additional `-I` import paths for every `protoc` invocation, so cross-repository imports just work.

### `plugins`

Each plugin entry maps to one `--<name>_out` flag passed to `protoc`. The `protoc-gen-` prefix is omitted.

| Field  | Description                                                              |
|--------|--------------------------------------------------------------------------|
| `name` | Plugin name without `protoc-gen-` prefix (e.g. `go`, `go-grpc`, `grpc-gateway`) |
| `out`  | Output directory for generated files                                     |
| `opts` | Key/value options passed as `--<name>_out=key=value,...:<out>`           |

Common `opts` for the Go plugins:

| Option            | Description                                                                                   |
|-------------------|-----------------------------------------------------------------------------------------------|
| `paths: source_relative` | Generate files relative to the proto file location rather than the full `go_package` path |
| `module: <prefix>` | Strip the given Go module prefix from output paths                                           |

### `openapi`

Run an OpenAPI generator (or any other non-protoc binary) as part of a generate rule. The binary is invoked with `flags` verbatim. Use `allow_custom: true` in `binaries` to install non-`protoc-gen-` binaries.

```yaml
generate:
  - openapi:
      binary: ogen
      flags:
        - "--target"
        - "internal/clients/todos"
        - "--package"
        - "todos"
        - "--clean"
        - "todos/openapi.yaml"
```

---

## Full example

See [`examples/full/moti.yaml`](../examples/full/moti.yaml) for a complete, annotated configuration covering all features.
