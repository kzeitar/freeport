package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/kzeitar/freeport/pkg/port"
	"github.com/shirou/gopsutil/v4/process"
)

// Version is set by the build process via ldflags
var version = "dev"

const (
	exitSuccess     = 0
	exitInvalidArgs = 2
	exitNotFound    = 3
	exitPermission  = 4
)

var (
	forceFlag   bool
	listFlag    bool
	dryRunFlag  bool
	verboseFlag bool
	versionFlag bool
)

func main() {
	flag.BoolVar(&forceFlag, "force", false, "don't prompt, kill immediately")
	flag.BoolVar(&forceFlag, "f", false, "don't prompt, kill immediately (shorthand)")
	flag.BoolVar(&listFlag, "list", false, "show PIDs using the specified port")
	flag.BoolVar(&dryRunFlag, "dry-run", false, "show what would be killed but do not kill")
	flag.BoolVar(&verboseFlag, "verbose", false, "enable verbose logging")
	flag.BoolVar(&verboseFlag, "v", false, "enable verbose logging (shorthand)")
	flag.BoolVar(&versionFlag, "version", false, "print version information")
	flag.BoolVar(&versionFlag, "V", false, "print version information (shorthand)")

	flag.Usage = usage

	// Pre-process os.Args to allow flags anywhere
	// Move all flags to the front, keep positional args at the end
	reorderedArgs := reorderArgs(os.Args[1:])
	flag.CommandLine.Parse(reorderedArgs)

	// Handle version flag
	if versionFlag {
		fmt.Printf("freeport %s\n", version)
		os.Exit(exitSuccess)
	}

	args := flag.Args()
	portNum, err := parsePort(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n\n", err)
		usage()
		os.Exit(exitInvalidArgs)
	}

	// Handle list mode
	if listFlag {
		if portNum == 0 {
			listAllListeningPorts()
		} else {
			listPIDsForPort(portNum)
		}
		os.Exit(exitSuccess)
	}

	// Require port number for kill operations
	if portNum == 0 {
		fmt.Fprintln(os.Stderr, "Error: port number is required for kill operations")
		usage()
		os.Exit(exitInvalidArgs)
	}

	// Find PIDs using the port
	pids, err := port.PIDsByPort(portNum)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to lookup PIDs: %v\n", err)
		os.Exit(exitPermission)
	}

	if len(pids) == 0 {
		fmt.Printf("No processes found using port %d\n", portNum)
		os.Exit(exitNotFound)
	}

	if verboseFlag {
		fmt.Printf("Found %d process(es) using port %d\n", len(pids), portNum)
	}

	// Kill each PID
	exitCode := exitSuccess
	for _, pid := range pids {
		if err := handlePID(pid, portNum); err != nil {
			exitCode = exitPermission
		}
	}

	os.Exit(exitCode)
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: freeport [options] <port>\n\n")
	fmt.Fprintf(os.Stderr, "Options:\n")
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, "\nExamples:\n")
	fmt.Fprintf(os.Stderr, "  freeport 8080              # Kill process on port 8080 (with prompt)\n")
	fmt.Fprintf(os.Stderr, "  freeport -f 8080           # Kill process on port 8080 without prompt\n")
	fmt.Fprintf(os.Stderr, "  freeport --list 8080       # Show PIDs using port 8080\n")
	fmt.Fprintf(os.Stderr, "  freeport --version         # Show version information\n")
	fmt.Fprintf(os.Stderr, "  freeport --dry-run 8080    # Show what would be killed\n")
}

func parsePort(args []string) (int, error) {
	if len(args) > 1 {
		return 0, fmt.Errorf("too many arguments")
	}

	if len(args) == 0 {
		return 0, nil
	}

	portStr := args[0]
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return 0, fmt.Errorf("invalid port number: %q", portStr)
	}

	if port < 1 || port > 65535 {
		return 0, fmt.Errorf("port must be between 1 and 65535, got %d", port)
	}

	return port, nil
}

func listPIDsForPort(portNum int) {
	pids, err := port.PIDsByPort(portNum)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to lookup PIDs: %v\n", err)
		os.Exit(exitPermission)
	}

	if len(pids) == 0 {
		fmt.Printf("No processes found using port %d\n", portNum)
		return
	}

	fmt.Printf("Processes using port %d:\n", portNum)
	for _, pid := range pids {
		name := getProcessName(pid)
		fmt.Printf("  PID %d (%s)\n", pid, name)
	}
}

func listAllListeningPorts() {
	fmt.Fprintln(os.Stderr, "Error: listing all listening ports is not supported.")
	fmt.Fprintln(os.Stderr, "Please specify a port number: freeport --list <port>")
	os.Exit(exitInvalidArgs)
}

func handlePID(pid int32, portNum int) error {
	name := getProcessName(pid)

	if verboseFlag || dryRunFlag {
		fmt.Printf("Found: %s (PID %d) using port %d\n", name, pid, portNum)
	}

	if dryRunFlag {
		fmt.Printf("[dry-run] Would kill: %s (PID %d)\n", name, pid)
		return nil
	}

	// Prompt for confirmation unless force flag is set
	if !forceFlag {
		if !promptForConfirmation(name, pid, portNum) {
			fmt.Printf("Skipped: %s (PID %d)\n", name, pid)
			return nil
		}
	}

	if verboseFlag {
		fmt.Printf("Killing %s (PID %d)...\n", name, pid)
	}

	if err := port.KillPID(pid); err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to kill PID %d: %v\n", pid, err)
		return err
	}

	fmt.Printf("Killed: %s (PID %d)\n", name, pid)
	return nil
}

func promptForConfirmation(name string, pid int32, portNum int) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Process %s (PID %d) is using port %d. Kill it? [y/N]: ", name, pid, portNum)

	response, err := reader.ReadString('\n')
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		return false
	}

	response = strings.TrimSpace(strings.ToLower(response))
	return response == "y" || response == "yes"
}

func getProcessName(pid int32) string {
	proc, err := process.NewProcess(pid)
	if err != nil {
		return fmt.Sprintf("process_%d", pid)
	}

	name, err := proc.Name()
	if err != nil {
		return fmt.Sprintf("process_%d", pid)
	}

	return name
}

// reorderArgs reorders arguments to place flags before positional arguments.
// This allows users to specify flags anywhere (e.g., "freeport 3000 --list").
func reorderArgs(args []string) []string {
	var flags, positionals []string

	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
			flags = append(flags, arg)
		} else {
			positionals = append(positionals, arg)
		}
	}

	// Flags first, then positionals
	return append(flags, positionals...)
}
