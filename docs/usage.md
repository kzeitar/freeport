# freeport Usage Guide

`freeport` is a command-line tool for freeing TCP ports by terminating the processes using them. It helps you quickly resolve "port already in use" errors during development and server management.

## Overview

When you encounter errors like "bind: address already in use" or "port 3000 is already in use," `freeport` helps you identify and terminate the process occupying that port.

## Quick Start

```bash
# Kill the process using port 8080 (with confirmation prompt)
freeport 8080

# Kill without confirmation
freeport -f 8080

# List processes using a specific port
freeport --list 8080

# Preview what would be killed (dry run)
freeport --dry-run 8080
```

## Usage Examples

### Interactive Mode (Default)

By default, `freeport` prompts for confirmation before killing each process:

```bash
$ freeport 3000
Process node (PID 12345) is using port 3000. Kill it? [y/N]: y
Killed: node (PID 12345)
```

**When to use:** Daily development where you want to verify before killing processes.

### Force Mode

Use the `-f` or `--force` flag to skip confirmation prompts:

```bash
$ freeport -f 8080
Found: python3 (PID 54321) using port 8080
Killing python3 (PID 54321)...
Killed: python3 (PID 54321)
```

**When to use:** Automation scripts, CI/CD pipelines, or when you're certain about the action.

### Dry Run Mode

Preview what would be killed without actually terminating processes:

```bash
$ freeport --dry-run 3000
Found: java (PID 9876) using port 3000
[dry-run] Would kill: java (PID 9876)
```

**When to use:** Safety check before running destructive operations, debugging, or understanding what's running.

### List Mode

List processes using a specific port:

```bash
$ freeport --list 5432
Processes using port 5432:
  PID 1234 (postgres)
  PID 5678 (postgres)
```

**When to use:** Investigating port conflicts without taking any action.

### Verbose Mode

Enable detailed logging with `-v` or `--verbose`:

```bash
$ freeport -v -f 8080
Found 1 process(es) using port 8080
Found: nginx (PID 23456) using port 8080
Killing nginx (PID 23456)...
Killed: nginx (PID 23456)
```

**When to use:** Debugging, scripts, or when you need detailed feedback.

## Combined Flags

Flags can be combined and specified in any order:

```bash
# Verbose dry run with confirmation
freeport 3000 --dry-run --verbose

# Force kill with verbose output
freeport -v -f 8080

# List with verbose output
freeport --list -v 3000
```

## Exit Codes

`freeport` uses the following exit codes:

- `0` - Success: Processes killed successfully or no processes found in list mode
- `2` - Invalid arguments: Bad port number, too many arguments, etc.
- `3` - Not found: No processes using the specified port
- `4` - Permission error: Insufficient privileges to kill processes or inspect network connections

## Platform Notes

### Linux

**Permissions:**
- Viewing network connections: No special permissions required
- Killing your own processes: No special permissions required
- Killing other users' processes: Requires `sudo`

**Example:**
```bash
# Kill your own process
freeport 3000

# Kill a system process or another user's process
sudo freeport 80
```

### macOS

**Permissions:**
- Viewing network connections: No special permissions required
- Killing your own processes: No special permissions required
- Killing system processes: Requires `sudo`

**Example:**
```bash
# Standard usage
freeport 3000

# Kill system process (e.g., Apache)
sudo freeport -f 80
```

### Windows

**Permissions:**
- Viewing network connections: No special permissions required
- Killing your own processes: No special permissions required
- Killing system services or other users' processes: Requires Administrator privileges

**Example:**
```cmd
REM Run as normal user for your own processes
freeport 3000

REM Run as Administrator for system processes
freeport -f 80
```

## Common Workflows

### Development: Restarting a Development Server

```bash
# Your Node.js app crashed but port 3000 is still occupied
freeport 3000
# Now you can restart your app
npm start
```

### Quick Port Cleanup

```bash
# Kill whatever is using port 8080 without asking
freeport -f 8080
```

### Investigating Port Conflicts

```bash
# First, see what's using the port
freeport --list 3000

# Then dry-run to verify
freeport --dry-run 3000

# Finally, kill it if you're sure
freeport 3000
```

### Scripting and Automation

```bash
#!/bin/bash
PORT=3000

# Check if port is in use
if freeport --list $PORT | grep -q "PID"; then
    echo "Port $PORT is in use, killing process..."
    freeport -f $PORT
fi

# Start your application
./my-app --port $PORT
```

### CI/CD Pipeline

```yaml
# Example GitHub Actions step
- name: Free up port 8080
  run: |
    freeport --dry-run 8080
    freeport -f 8080 || true  # Continue even if no process was found
```

## Troubleshooting

### "Permission denied" or "failed to lookup PIDs"

**Cause:** Insufficient permissions to view network connections or kill processes.

**Solution:** Run with elevated privileges:
```bash
sudo freeport 80
```

### "No processes found using port X"

**Cause:** Either the port is not in use, or the process has already terminated.

**Solution:** Verify the port number:
```bash
# List what's listening (Linux/macOS)
netstat -an | grep :3000
# or
lsof -i :3000

# Then use freeport
freeport --list 3000
```

### Process reappears after killing

**Cause:** The process is managed by a service manager (systemd, launchd, etc.) or is being respawned.

**Solution:** Stop the service instead:
```bash
# Linux with systemd
sudo systemctl stop nginx

# macOS with launchd
sudo launchctl unload /Library/LaunchDaemons/com.example.service.plist
```

## Best Practices

1. **Use `--dry-run` first** when uncertain about what processes will be affected
2. **Use `--list` mode** to investigate before taking destructive action
3. **Avoid `--force` in interactive sessions** unless you're absolutely certain
4. **Check exit codes** in scripts to handle different scenarios appropriately
5. **Run with `sudo` only when necessary** to minimize security risks

## See Also

- [README.md](../README.md) - Project overview and installation
- [Go package documentation](../pkg/port/) - API documentation for integrating `freeport` into your Go programs
