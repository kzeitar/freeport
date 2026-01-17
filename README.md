# freeport

`freeport` is a command-line tool that helps you quickly free up TCP ports by identifying and terminating the processes using them. It's particularly useful during development when you encounter "port already in use" errors.

## Features

- 🔍 **Find processes** by port number across all network connections
- ⚡ **Kill processes** gracefully (SIGTERM) or forcefully (SIGKILL)
- 🎯 **Interactive confirmation** prompts for safety
- 🔒 **Dry-run mode** to preview what would be killed
- 📋 **List mode** to inspect without taking action
- 🖥️ **Cross-platform** support (Linux, macOS, Windows)
- 🚦 **Clear exit codes** for scripting and automation

## Installation

`freeport` can be installed via multiple methods. Choose the one that works best for you:

### Go Toolchain (Recommended)

```bash
go install github.com/kzeitar/freeport@latest
```

This will install the `freeport` binary to `$GOPATH/bin` (or `$HOME/go/bin` by default). Make sure this directory is in your `$PATH`:

```bash
# Add to your ~/.bashrc, ~/.zshrc, or equivalent
export PATH=$PATH:$(go env GOPATH)/bin
```

### Homebrew

```bash
brew tap kzeitar/freeport
brew install freeport
```

### npm

```bash
npm install -g freeport-cli
```

### Binary Releases

Download the appropriate binary for your platform from the [releases page](https://github.com/kzeitar/freeport/releases).

**For more detailed installation instructions, see [INSTALL.md](INSTALL.md).**

## Quick Start

```bash
# Kill the process using port 3000 (with confirmation)
freeport 3000

# Kill without confirmation
freeport -f 3000

# See what's using the port
freeport --list 3000

# Preview what would be killed
freeport --dry-run 3000
```

## Usage

```bash
freeport [options] <port>
```

### Options

| Flag | Shorthand | Description |
|------|-----------|-------------|
| `--force` | `-f` | Kill processes without confirmation prompt |
| `--list` | | List processes using the specified port (no killing) |
| `--dry-run` | | Show what would be killed without actually killing |
| `--verbose` | `-v` | Enable verbose logging |
| `--version` | `-V` | Print version information |
| `--help` | `-h` | Show help message |

### Arguments

| Argument | Description |
|----------|-------------|
| `port` | TCP port number (1-65535) |

### Exit Codes

- `0` - Success
- `2` - Invalid arguments
- `3` - No processes found using the port
- `4` - Permission denied

## Examples

### Interactive Kill (Default)

```bash
$ freeport 8080
Process node (PID 12345) is using port 8080. Kill it? [y/N]: y
Killed: node (PID 12345)
```

### Force Kill

```bash
$ freeport -f 3000
Found: python3 (PID 54321) using port 3000
Killing python3 (PID 54321)...
Killed: python3 (PID 54321)
```

### List Processes

```bash
$ freeport --list 5432
Processes using port 5432:
  PID 1234 (postgres)
  PID 5678 (postgres)
```

### Show Version

```bash
$ freeport --version
freeport 0.1.0
```

### Dry Run

```bash
$ freeport --dry-run 3000
Found: java (PID 9876) using port 3000
[dry-run] Would kill: java (PID 9876)
```

### Combined Flags

```bash
# Verbose force kill
freeport -v -f 8080

# Flags can be specified in any order
freeport 3000 --dry-run --verbose
```

## Platform-Specific Notes

### Linux & macOS

- Viewing ports: No special permissions required
- Killing your own processes: No special permissions required
- Killing other users' processes: Requires `sudo`

```bash
sudo freeport -f 80  # Kill process on port 80 (typically requires root)
```

### Windows

- Run as Administrator when killing system processes or services
- Standard permissions suffice for your own processes

## Use Cases

- **Development:** Quickly free ports when your development server crashes and leaves the port occupied
- **Testing:** Clean up ports between test runs
- **Automation:** Integrate into CI/CD pipelines to ensure clean environments
- **Debugging:** Investigate which processes are using specific ports

## Documentation

- [Installation Guide](INSTALL.md) - Detailed installation instructions for all platforms
- [Usage Guide](docs/usage.md) - Comprehensive usage documentation with troubleshooting
- [Go Package Documentation](pkg/port/) - API documentation for using `freeport` as a Go library

## Contributing

Contributions are welcome! Here's how to get started:

### Development Setup

```bash
# Clone the repository
git clone https://github.com/kzeitar/freeport.git
cd freeport

# Run tests
make test

# Build locally
make build

# Format code
make fmt
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests with coverage
go test -cover ./...
```

### Code Style

- Follow standard Go conventions (use `gofmt`)
- Add doc comments to exported functions
- Write tests for new features
- Keep commits small and focused

### Submitting Changes

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

[MIT License](LICENSE) - feel free to use this tool in your projects.

## Acknowledgments

Built with [gopsutil](https://github.com/shirou/gopsutil) for cross-platform process and network inspection.
