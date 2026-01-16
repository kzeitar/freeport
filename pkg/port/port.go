// Package port provides utilities for working with TCP ports.
package port

import (
	"os"
	"runtime"
	"syscall"
	"time"

	gopsnet "github.com/shirou/gopsutil/v4/net"
)

// PIDsByPort returns the PIDs of processes using the given TCP port.
// It checks both listening and established connections, supporting both IPv4 and IPv6.
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

// KillPID terminates the process with the given PID.
// On Unix, it sends SIGTERM first and falls back to SIGKILL if needed.
// On Windows, it calls Process.Kill() directly.
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
