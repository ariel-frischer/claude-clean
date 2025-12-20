package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/ariel-frischer/claude-clean/display"
	"github.com/ariel-frischer/claude-clean/parser"
)

// Version is set at build time
var Version = "dev"

// Command line flags
var (
	verbose     = flag.Bool("V", false, "Show verbose output (usage stats, tool IDs)")
	showVersion = flag.Bool("v", false, "Show version")
	styleFlag   = flag.String("s", "default", "Output style: default, compact, minimal, plain")
	showLineNum = flag.Bool("n", false, "Show line numbers")
	uninstall   = flag.Bool("uninstall", false, "Uninstall cclean from the system")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS] [FILE]\n\n", binaryName())
		fmt.Fprintln(os.Stderr, "Transform Claude Code's stream-json output into readable terminal output.")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "Arguments:")
		fmt.Fprintln(os.Stderr, "  FILE             JSONL file to process (optional)")
		fmt.Fprintln(os.Stderr, "  No arguments     Reads from stdin")
		fmt.Fprintln(os.Stderr, "\nOptions:")
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr, "\nStyles:")
		fmt.Fprintln(os.Stderr, "  default  - Full output with colored boxes and borders")
		fmt.Fprintln(os.Stderr, "  compact  - Single-line summaries for each message")
		fmt.Fprintln(os.Stderr, "  minimal  - Clean output without box-drawing characters")
		fmt.Fprintln(os.Stderr, "  plain    - No colors, suitable for piping")
		fmt.Fprintln(os.Stderr, "\nExamples:")
		fmt.Fprintf(os.Stderr, "  claude -p 'prompt' --output-format stream-json | %s\n", binaryName())
		fmt.Fprintf(os.Stderr, "  %s output.jsonl             # Process a JSONL file\n", binaryName())
		fmt.Fprintf(os.Stderr, "  %s -s compact output.jsonl  # Use compact style\n", binaryName())
	}

	flag.Parse()

	if *showVersion {
		fmt.Printf("%s version %s\n", binaryName(), Version)
		os.Exit(0)
	}

	if *uninstall {
		runUninstall()
		os.Exit(0)
	}

	// Determine output style
	var style display.OutputStyle
	switch *styleFlag {
	case "default":
		style = display.StyleDefault
	case "compact":
		style = display.StyleCompact
	case "minimal":
		style = display.StyleMinimal
	case "plain":
		style = display.StylePlain
	default:
		fmt.Fprintf(os.Stderr, "Unknown style: %s\n", *styleFlag)
		flag.Usage()
		os.Exit(1)
	}

	cfg := &display.Config{
		Style:       style,
		Verbose:     *verbose,
		ShowLineNum: *showLineNum,
	}

	args := flag.Args()

	switch len(args) {
	case 0:
		// No file argument - read from stdin
		processStream(os.Stdin, cfg)
	case 1:
		if args[0] == "-" {
			// Read from stdin (explicit)
			processStream(os.Stdin, cfg)
		} else if fileExists(args[0]) {
			// Process file
			processFile(args[0], cfg)
		} else {
			fmt.Fprintf(os.Stderr, "File not found: %s\n", args[0])
			os.Exit(1)
		}
	default:
		fmt.Fprintln(os.Stderr, "Too many arguments")
		flag.Usage()
		os.Exit(1)
	}
}

func binaryName() string {
	return filepath.Base(os.Args[0])
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

func processFile(filename string, cfg *display.Config) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	processStream(file, cfg)
}

func processStream(r *os.File, cfg *display.Config) {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 0, parser.MaxBufferCapacity), parser.MaxBufferCapacity)

	lineNum := 0
	var lastAssistantContent string

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		if line == "" {
			continue
		}

		var msg parser.StreamMessage
		if err := json.Unmarshal([]byte(line), &msg); err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing line %d: %v\n", lineNum, err)
			continue
		}

		// Skip duplicate result messages that contain the same content as the last assistant message
		if msg.Type == "result" && msg.Result != "" && msg.Result == lastAssistantContent {
			continue
		}

		// Track assistant message content for duplicate detection
		if msg.Type == "assistant" && msg.Message != nil && len(msg.Message.Content) > 0 {
			for _, block := range msg.Message.Content {
				if block.Type == "text" && block.Text != "" {
					lastAssistantContent = block.Text
				}
			}
		}

		display.DisplayMessage(&msg, lineNum, cfg)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading: %v\n", err)
		os.Exit(1)
	}
}

func runUninstall() {
	// Get the path to the current executable
	execPath, err := os.Executable()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error finding executable path: %v\n", err)
		os.Exit(1)
	}

	// Resolve any symlinks to get the real path
	realPath, err := filepath.EvalSymlinks(execPath)
	if err != nil {
		realPath = execPath
	}

	// Check if running from common install locations
	homeDir, _ := os.UserHomeDir()
	localBin := filepath.Join(homeDir, ".local", "bin", "cclean")
	systemBin := "/usr/local/bin/cclean"

	var pathsToRemove []string

	// Add the actual executable path
	pathsToRemove = append(pathsToRemove, realPath)

	// Check for binary in standard locations
	for _, path := range []string{localBin, systemBin} {
		if path != realPath {
			if _, err := os.Stat(path); err == nil {
				pathsToRemove = append(pathsToRemove, path)
			}
		}
	}

	if len(pathsToRemove) == 0 {
		fmt.Println("No cclean installation found.")
		return
	}

	fmt.Println("The following files will be removed:")
	for _, path := range pathsToRemove {
		fmt.Printf("  %s\n", path)
	}
	fmt.Print("\nProceed with uninstall? [y/N] ")

	var response string
	fmt.Scanln(&response)

	if response != "y" && response != "Y" {
		fmt.Println("Uninstall cancelled.")
		return
	}

	var needsSudo bool
	for _, path := range pathsToRemove {
		// Check if we need sudo (file in /usr/local/bin and not writable)
		if strings.HasPrefix(path, "/usr/local") {
			if f, err := os.OpenFile(path, os.O_WRONLY, 0); err != nil {
				needsSudo = true
			} else {
				f.Close()
			}
		}
	}

	for _, path := range pathsToRemove {
		var removeErr error

		if needsSudo && filepath.HasPrefix(path, "/usr/local") {
			// Use sudo for system paths
			cmd := exec.Command("sudo", "rm", "-f", path)
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			removeErr = cmd.Run()
		} else {
			removeErr = os.Remove(path)
		}

		if removeErr != nil {
			fmt.Fprintf(os.Stderr, "Error removing %s: %v\n", path, removeErr)
		} else {
			fmt.Printf("Removed: %s\n", path)
		}
	}

	fmt.Println("\ncclean has been uninstalled.")
}
