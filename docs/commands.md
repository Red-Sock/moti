# Commands

## Global flags

These flags apply to every command.

| Flag        | Default     | Description                    |
|-------------|-------------|--------------------------------|
| `--cfg`     | `moti.yaml` | Path to the configuration file |
| `--version` | —           | Print moti version and exit    |

Example — use a config in a non-standard location:

```bash
moti --cfg configs/moti.yaml install
```

---

## `moti install`

**Alias:** `moti i`

Fetches all proto dependencies and installs all binaries declared in `moti.yaml`.

```bash
moti install
```

### What it does

1. For each entry in `binaries.install`, checks whether the binary is already installed at the required version by
   running it with `version_check_args`. Installs via `go install` if not.
2. For each entry in `deps`, checks whether the module is already present in the local cache. If not:
    - Clones or fetches the remote git repository.
    - Reads the remote's own `moti.yaml` to discover transitive dependencies and installs them first.
    - Archives the proto files into `<cache_path>/mod/<module>/<version>/`.
    - Computes a content hash and writes an entry to `moti.lock`.

### Skipping

- A binary is skipped when it is already installed at the pinned version.
- A module is skipped when its entry in `moti.lock` matches the cached content hash on disk.
- A module listed under `replace` is skipped entirely — the local path is used as-is.

---

## `moti generate`

**Alias:** `moti g`

Generates code by assembling and running `protoc` commands from the `generate` section of `moti.yaml`.

```bash
moti generate
```

### What it does

For each entry in `generate`:

1. Creates output directories for all plugins if they do not exist.
2. For each `input`:
    - Resolves the input root (local directory or cached git repository path).
    - Walks the input root and collects all `.proto` files.
    - Appends all `deps` modules as additional `-I` import paths.
    - Builds the full `protoc` invocation and logs it.
    - Runs `protoc` via the shell.
3. If an `openapi` block is present, runs the specified binary with the given flags.

When `binaries.bin_dir` is set, the directory is prepended to `PATH` before each command so project-local binaries are
preferred.

### Requirements

- All modules referenced in `generate[].inputs[].git_repo.url` must be installed first. Run `moti install` if you see
  a "module not installed" error.

### Logged command

moti prints the exact `protoc` command it runs before executing it. This makes it easy to debug generation issues — copy
the command and run it directly.

---

## Exit codes

| Code | Meaning                           |
|------|-----------------------------------|
| `0`  | Success                           |
| `1`  | Error (message printed to stderr) |
