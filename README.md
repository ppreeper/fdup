# fdup

A fast duplicate file finder. Recursively walks directories, groups files by BLAKE2b-256 content hash, and reports duplicates.

## Features

- **Size-based pre-filtering** -- files with unique sizes are skipped before hashing, dramatically reducing I/O
- **Streaming hashes** -- files are hashed via streaming I/O, not loaded entirely into memory
- **Parallel hashing** -- worker pool scales to available CPUs
- **Deterministic output** -- groups and paths are sorted for reproducible results
- **Safe delete mode** -- generates shell commands with proper quoting

## Installation

```
go install github.com/ppreeper/fdup@latest
```

## Usage

```
fdup [options] [directories...]
```

If no directories are specified, the current directory (`.`) is used.

### Options

| Flag | Description |
|------|-------------|
| `-d` | Print all files as commented-out `rm` commands (`#rm -vf ...`) |
| `-dd` | Keep the first file (commented), emit active `rm` commands for the rest |
| `-m`, `--matchcount N` | Minimum number of files in a duplicate group (default: 2) |
| `--skip-empty` | Skip zero-length files |

### Examples

Find duplicates in the current directory:

```
fdup
```

Find duplicates across multiple directories:

```
fdup /photos /backup/photos
```

Generate delete commands (all commented out for review):

```
fdup -d /data > cleanup.sh
```

Generate delete commands keeping the first copy of each duplicate:

```
fdup -dd /data > cleanup.sh
# Review cleanup.sh, then:
# bash cleanup.sh
```

Only show groups with 3 or more copies:

```
fdup -m 3 /data
```

## License

Apache License 2.0 -- see [LICENSE](LICENSE).
