# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.0] - 2026-01-16

### Added
- TCP port process finder: locate processes using specific ports
- Process termination with graceful (SIGTERM) and forceful (SIGKILL) kill modes
- Interactive confirmation prompts for safety
- `--force` flag to skip confirmation prompts
- `--list` flag to inspect processes without killing them
- `--dry-run` flag to preview what would be killed
- `--verbose` flag for detailed logging
- Cross-platform support (Linux, macOS, Windows)
- Clear exit codes for scripting (0=success, 2=invalid args, 3=not found, 4=permission denied)
- Go package API for programmatic use

### Contributors
- Khaled Zeitar
