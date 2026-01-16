// Package port provides utilities for working with TCP ports.
package port

import (
	"os"
	"runtime"
	"syscall"
	"time"

	gopsnet "github.com/shirou/gopsutil/v4/net"
)

// PIDsByPort returns the process IDs (PIDs) of all processes using the given TCP port number.
//
// It scans both listening and established connections, supporting IPv4 and IPv6 addresses.
// The function checks both local and remote endpoints of each connection, ensuring that
// any process bound to or connected on the specified port is detected.
//
// Parameters:
//   - port: The TCP port number to search for (typically 1-65535)
//
// Returns:
//   - []int32: A slice of PIDs using the specified port (deduplicated, no specific order)
//   - error: An error if the system's network connection table cannot be accessed
//
// Example:
//
//	pids, err := port.PIDsByPort(8080)
//	if err != nil {
//	    log.Fatalf("Failed to lookup PIDs: %v", err)
//	}
//	for _, pid := range pids {
//	    fmt.Printf("Process %d is using port 8080\n", pid)
//	}
func PIDsByPort(port int) ([]int32, error) {
	conns, err := gopsnet.Connections("all")
	if err != nil {
		return nil, err
	}

	pidMap := make(map[int32]struct{})
	for _, conn := range conns {
		// Check for TCP connections (SOCK_STREAM = 1)
		if conn.Type == syscall.SOCK_STREAM {
			// Check both local and remote addresses
			if conn.Laddr.Port == uint32(port) {
				if conn.Pid != 0 {
					pidMap[conn.Pid] = struct{}{}
				}
			}
			if conn.Raddr.Port == uint32(port) {
				if conn.Pid != 0 {
					pidMap[conn.Pid] = struct{}{}
				}
			}
		}
	}

	result := make([]int32, 0, len(pidMap))
	for pid := range pidMap {
		result = append(result, pid)
	}
	return result, nil
}

// KillPID terminates the process with the given PID in a platform-appropriate manner.
//
// On Unix-like systems (Linux, macOS, *BSD), it attempts graceful termination first
// by sending SIGTERM, waiting up to 3 seconds for the process to exit cleanly.
// If the process doesn't terminate within the timeout, it sends SIGKILL to force termination.
// If SIGTERM cannot be sent (e.g., due to permissions), it falls back to SIGKILL immediately.
//
// On Windows, it calls Process.Kill() directly for immediate termination.
//
// Parameters:
//   - pid: The process ID to terminate
//
// Returns:
//   - error: An error if the process cannot be found or termination fails
//
// Permission requirements:
//   - Unix: You must own the process or have root/sudo privileges to kill it
//   - Windows: You must have sufficient privileges (typically requires Administrator for system processes)
//
// Example:
//
//	pid := int32(1234)
//	if err := port.KillPID(pid); err != nil {
//	    log.Fatalf("Failed to kill process %d: %v", pid, err)
//	}
//	fmt.Printf("Successfully killed process %d\n", pid)
func KillPID(pid int32) error {
	proc, err := os.FindProcess(int(pid))
	if err != nil {
		return err
	}

	if runtime.GOOS == "windows" {
		return proc.Kill()
	}

	// Unix: try SIGTERM first, then SIGKILL
	if err := proc.Signal(syscall.SIGTERM); err != nil {
		// If SIGTERM fails, try SIGKILL immediately
		return proc.Signal(syscall.SIGKILL)
	}

	// Give the process a moment to exit gracefully
	done := make(chan error, 1)
	go func() {
		_, err := proc.Wait()
		done <- err
	}()

	select {
	case <-done:
		return nil
	case <-time.After(3 * time.Second):
		// Force kill if not terminated
		return proc.Signal(syscall.SIGKILL)
	}
}
