package port

import (
	"net"
	"os"
	"os/exec"
	"runtime"
	"syscall"
	"testing"
	"time"
)

func TestPIDsByPort_returnsCurrentPID(t *testing.T) {
	// Start a listener on a random port
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to create listener: %v", err)
	}
	defer listener.Close()

	// Get the actual port number
	addr := listener.Addr().(*net.TCPAddr)
	port := addr.Port

	// Give the system a moment to register the connection
	time.Sleep(100 * time.Millisecond)

	// Call PIDsByPort
	pids, err := PIDsByPort(port)
	if err != nil {
		t.Fatalf("PIDsByPort(%d) failed: %v", port, err)
	}

	// Check that our PID is in the list
	myPID := int32(os.Getpid())
	found := false
	for _, pid := range pids {
		if pid == myPID {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Expected PID %d to be in PIDsByPort(%d) result, got %v", myPID, port, pids)
	}
}

func TestKillPID_onChildProcess(t *testing.T) {
	// Start the helper process
	cmd := exec.Command("go", "run", "testdata/listener.go", "-port", "0")
	cmd.Stdout = nil // Discard output
	cmd.Stderr = nil // Discard errors

	if err := cmd.Start(); err != nil {
		t.Fatalf("Failed to start helper process: %v", err)
	}

	childPID := int32(cmd.Process.Pid)

	// Give the process time to start and begin listening
	time.Sleep(500 * time.Millisecond)

	// Verify the process is running
	proc, err := os.FindProcess(int(childPID))
	if err != nil {
		t.Fatalf("Failed to find child process: %v", err)
	}

	// On Unix, signal 0 checks if process exists without killing it
	if runtime.GOOS != "windows" {
		if err := proc.Signal(syscall.Signal(0)); err != nil {
			t.Fatalf("Child process %d is not running: %v", childPID, err)
		}
	}

	// Kill the process using our function
	if err := KillPID(childPID); err != nil {
		t.Fatalf("KillPID(%d) failed: %v", childPID, err)
	}

	// Wait for the process to actually exit
	// cmd.Wait() should return quickly since we killed the process
	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case err := <-done:
		if err != nil {
			// Process exited with an error (expected since we killed it)
			t.Logf("Child process exited with error (expected): %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatalf("Child process did not exit within 5 seconds after KillPID")
	}

	// Verify the process is gone by trying signal 0 again
	if runtime.GOOS != "windows" {
		if err := proc.Signal(syscall.Signal(0)); err == nil {
			t.Errorf("Process %d is still running after KillPID", childPID)
		} else {
			t.Logf("Process successfully terminated (signal 0 returned error as expected)")
		}
	}
}

// TestPIDsByPort_nonExistentPort verifies that querying an unused port returns an empty list
func TestPIDsByPort_nonExistentPort(t *testing.T) {
	// Use a very high port number that's unlikely to be in use
	port := 65432

	pids, err := PIDsByPort(port)
	if err != nil {
		t.Fatalf("PIDsByPort(%d) failed: %v", port, err)
	}

	if len(pids) != 0 {
		t.Errorf("Expected no PIDs for unused port %d, got %d PIDs: %v", port, len(pids), pids)
	}
}

// TestKillPID_invalidPID verifies that killing a non-existent process returns an error
func TestKillPID_invalidPID(t *testing.T) {
	// Use a PID that's unlikely to exist
	invalidPID := int32(999999)

	if err := KillPID(invalidPID); err == nil {
		t.Errorf("Expected error when killing non-existent PID %d, got nil", invalidPID)
	} else {
		t.Logf("Got expected error for invalid PID: %v", err)
	}
}
