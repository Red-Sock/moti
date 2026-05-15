# moti

`moti` is a CLI tool for managing Protocol Buffers workflows. It wraps `protoc` with declarative dependency management and code generation driven by a single `moti.yaml` config file.

## Why moti?

Maintaining raw `protoc` commands is painful. A typical command grows like this:

```bash
protoc \
  -I proto_modules/mod/github.com/googleapis/googleapis/v0.0.0-abc123 \
  -I proto_modules/mod/github.com/vervstack/Matreshka/v1.0.95 \
  --go_out=module=go.vervstack.ru:internal/api/clients \
  --go-grpc_out=module=go.vervstack.ru:internal/api/clients \
  api/grpc/matreshka_api.proto
```

With moti you declare your dependencies and generation rules once, then just run:

```bash
moti install
moti generate
```

## What moti does

- **`moti install`** — fetches proto dependencies from git repositories into a local cache and installs the required `protoc-gen-*` binaries.
- **`moti generate`** — assembles and runs `protoc` commands from your config, resolving all `-I` import paths automatically.

## Pages

- [Installation](installation.md)
- [Configuration reference](configuration.md)
- [Commands](commands.md)
- [Key concepts](concepts.md)