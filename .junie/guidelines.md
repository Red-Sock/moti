### Project Overview
- `moti` is a CLI tool for managing Protocol Buffers workflows.
- It's a fork of `easy-p`.
- Key commands: `install` (manages dependencies), `generate` (generates code from proto files).

### Development Guidelines
- Use Go 1.24 idioms (e.g., `any`, `slices` package, `for i := range n`).
- Follow the existing project structure:
  - `internal/models`: Domain models and errors.
  - `internal/adapters`: Implementations of various interfaces (repository, storage, console).
  - `internal/commands`: CLI command implementations.
- Avoid using `interface{}`; use `any` instead.
- Use `rerrors` for error wrapping (e.g., `rerrors.Wrap(err, "context")`).
- Maintain existing naming conventions (e.g., `IStorage` for interfaces, `Storage` for implementations).

### Testing
- Run tests using `go test ./...`.
- Most packages currently lack extensive test coverage.
- When adding new features or fixing bugs, try to add reproduction tests if possible.
- Mocks are generated using `minimock`. If you modify an interface, run `go generate ./...` to update mocks.

### Working with Junie
- You can update these guidelines by editing this file.
- When performing refactors, check for unused code and mark it with `// TODO UNUSED` if you're not sure about deleting it yet.
- Focus on keeping the versioning logic consistent (prefer commit hashes over pseudo-versions for generated versions).
