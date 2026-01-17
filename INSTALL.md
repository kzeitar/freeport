# Installation Guide

`freeport` can be installed via multiple package managers. Choose the method that best fits your workflow.

## Table of Contents

- [Go Toolchain (Recommended)](#go-toolchain-recommended)
- [Homebrew](#homebrew)
- [npm](#npm)
- [Manual Binary Installation](#manual-binary-installation)
- [Building from Source](#building-from-source)

---

## Go Toolchain (Recommended)

If you have Go installed, this is the simplest method:

```bash
go install github.com/kzeitar/freeport@latest
```

The binary will be installed to `$GOPATH/bin` (or `$HOME/go/bin` by default). Make sure this directory is in your `PATH`:

```bash
# Add to your ~/.bashrc, ~/.zshrc, or equivalent
export PATH=$PATH:$(go env GOPATH)/bin
```

**Verify installation:**

```bash
freeport --help
```

---

## Homebrew

### Option 1: Install via Tap

```bash
# Add the tap
brew tap kzeitar/freeport

# Install
brew install freeport
```

### Option 2: Install from Local Formula

```bash
brew install --formula https://raw.githubusercontent.com/kzeitar/freeport/main/.github/homebrew/freeport.rb
```

**Update:**

```bash
brew upgrade freeport
```

**Uninstall:**

```bash
brew uninstall freeport
```

---

## npm

**Note:** This method downloads pre-compiled binaries and does not require Node.js for runtime.

```bash
# Install globally
npm install -g freeport-cli

# Or using yarn
yarn global add freeport-cli
```

**Update:**

```bash
npm update -g freeport-cli
```

**Uninstall:**

```bash
npm uninstall -g freeport-cli
```

---

## Manual Binary Installation

### 1. Download the Binary

Download the appropriate binary for your platform from the [Releases](https://github.com/kzeitar/freeport/releases) page.

| Platform | Binary Name |
|----------|-------------|
| Linux (x64) | `freeport-linux-amd64` |
| Linux (ARM64) | `freeport-linux-arm64` |
| macOS (Intel) | `freeport-darwin-amd64` |
| macOS (Apple Silicon) | `freeport-darwin-arm64` |
| Windows (x64) | `freeport-windows-amd64.exe` |

### 2. Make Executable (Linux/macOS)

```bash
chmod +x freeport-*
```

### 3. Move to PATH

```bash
# Linux/macOS
sudo mv freeport-* /usr/local/bin/freeport

# Or without sudo (for current user)
mkdir -p ~/.local/bin
mv freeport-* ~/.local/bin/freeport
# Add ~/.local/bin to your PATH if not already there
```

For Windows, add the binary to a directory in your PATH or move to `C:\Program Files\`.

---

## Building from Source

### Prerequisites

- Go 1.20 or later
- Make (optional, for convenience targets)

### Build Steps

```bash
# Clone the repository
git clone https://github.com/kzeitar/freeport.git
cd freeport

# Build
go build -o freeport ./cmd/freeport

# Or use Make
make build
```

The binary will be created in the current directory (or `bin/` if using Make).

### Cross-Platform Build

Build binaries for multiple platforms:

```bash
make release
```

This creates binaries for:
- Linux (amd64, arm64)
- macOS (amd64, arm64)
- Windows (amd64)

Binaries are placed in the `bin/` directory.

---

## Verification

After installation, verify that `freeport` is working:

```bash
freeport --help
```

You should see the help output with usage instructions.

Test with an actual port:

```bash
# List processes using a port (safe, doesn't kill)
freeport --list 3000
```

---

## Troubleshooting

### "command not found: freeport"

The installation directory is not in your PATH. Add it:

```bash
# For Go installations
export PATH=$PATH:$(go env GOPATH)/bin

# For manual installations to /usr/local/bin
export PATH=$PATH:/usr/local/bin
```

Add the appropriate line to your shell configuration (`~/.bashrc`, `~/.zshrc`, etc.).

### "permission denied" when killing processes

You need elevated privileges to kill processes you don't own:

```bash
# Linux/macOS
sudo freeport 3000

# Windows (run as Administrator)
# Right-click Command Prompt/PowerShell -> Run as Administrator
```

### npm installation fails

If the automatic binary download fails, install via Go instead:

```bash
go install github.com/kzeitar/freeport@latest
```

---

## Next Steps

- Read the [Usage Guide](docs/usage.md)
- Check out the [README](README.md) for examples
- Report issues at [GitHub Issues](https://github.com/kzeitar/freeport/issues)
