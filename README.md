# Gopate

**Gopate** is a pure command-line file format disguise tool built with Go, enabling fast and simple file format disguise and restoration.

> This project is a rewrite of [rippod/apate](https://github.com/rippod/apate) (C#), **fully compatible** with its disguised file format and can be used interchangeably.

## ✨ Features

- 🚀 **Ultra-fast** — Handles large files instantly
- 🔒 **Secure** — Original file headers are encrypted, making detection difficult
- 🌍 **Cross-platform** — Supports Linux / macOS / Windows / FreeBSD
- 🌐 **i18n** — Auto-detects terminal language, supports English & Chinese
- 🔄 **Fully Compatible** — Interoperable with C# version apate
- 📦 **Batch Processing** — Supports glob patterns and recursive directories
- 🔍 **File Inspection** — Detect if files are disguised and view original headers
- 📋 **Preview Mode** — `--dry-run` to preview operations without modifying files
- 🛡️ **Safe Mode** — Creates new files by default, `--in-place` for in-place modification

## Installation

### Download Pre-built Binaries

Download the binary for your platform from the [Releases](https://github.com/maolei1024/gopate/releases) page.

### Build from Source

```bash
git clone https://github.com/maolei1024/gopate.git
cd gopate
go build -o gopate .
```

### Go Install

```bash
go install github.com/maolei1024/gopate@latest
```

## Usage

### Language Settings

Gopate automatically detects your terminal language from environment variables (`LANG`, `LC_ALL`, etc.). English is the default language.

```bash
# Force Chinese output
gopate --lang zh --help

# Force English output
gopate --lang en --help

# Or set via environment variable
LANG=zh_CN.UTF-8 gopate --help
```

### Disguise Files

```bash
# One-key disguise (default mode, disguise as MP4)
gopate disguise secret.zip

# Disguise as specific format
gopate disguise secret.zip --mode exe    # Disguise as EXE
gopate disguise secret.zip --mode jpg    # Disguise as JPG
gopate disguise secret.zip --mode mp4    # Disguise as MP4
gopate disguise secret.zip --mode mov    # Disguise as MOV

# Use custom mask file
gopate disguise secret.zip --mode mask --mask-file cover.png

# In-place modification (no new file created)
gopate disguise secret.zip --mode onekey --in-place

# Batch processing
gopate disguise *.zip --mode onekey --in-place

# Recursive directory processing
gopate disguise ./mydir -r --mode mp4 --in-place

# Output to specified directory
gopate disguise secret.zip -o ./output/
```

### Reveal Files

```bash
# Reveal a file (creates new file)
gopate reveal secret.zip.mp4

# In-place reveal
gopate reveal secret.zip.mp4 --in-place

# Skip confirmation prompt
gopate reveal secret.zip.mp4 --in-place -f

# Batch reveal
gopate reveal *.mp4 --in-place -f
```

### Inspect Files

```bash
# Detect if a file is disguised
gopate inspect suspicious.mp4

# Show detailed info (with original header)
gopate inspect suspicious.mp4 -v
```

### Other Options

```bash
gopate version          # Show version
gopate --help           # Show help
gopate disguise --help  # Show disguise command help

# Global options
-v, --verbose       Show verbose output
-q, --quiet         Quiet mode, show errors only
    --lang string   Set language (en, zh)
```

## Disguise Modes

| Mode | Description | Mask Source |
|------|-------------|-------------|
| `onekey` | One-key disguise as MP4 (default, suitable for most cases) | Embedded MP4 file |
| `mask` | Custom mask file disguise | User-specified file |
| `exe` | Simple EXE disguise | PE file header |
| `jpg` | Simple JPG disguise | JPEG file header |
| `mp4` | Simple MP4 disguise | MP4 file header |
| `mov` | Simple MOV disguise | MOV file header |

## ⚠️ Important Notes

1. **Always back up your data before use**
2. Revealing files that were not disguised may cause data corruption
3. This software must not be used for illegal purposes; users bear all responsibility

## Compatibility with apate

Gopate is fully compatible with [apate](https://github.com/rippod/apate) (C# version v1.4.2) using an identical binary file format:

- ✅ Files disguised by Gopate can be revealed by apate
- ✅ Files disguised by apate can be revealed by Gopate
- ✅ All mask bytes are byte-for-byte identical with apate

## License

This project is open-sourced under the [MIT License](LICENSE).
