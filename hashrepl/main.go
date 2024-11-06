package main

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

const replName string = "hashREPL"

var (
	tmpDir string
)

func init() {
	// Get the temporary directory
	tmpDir = os.TempDir()
}

func prompt() {
	fmt.Print(replName, "> ")
}

func clear() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func help() {
	fmt.Printf(
		"Welcome to %v! Available commands: \n",
		replName,
	)
	fmt.Println(".help    - Show available commands.")
	fmt.Println(".clear   - Clear the screen.")
	fmt.Println(".exit    - Closes the REPL.")
}

func main() {
	commands := map[string]interface{}{
		".help":  help,
		".clear": clear,
	}
	// Begin the repl loop
	reader := bufio.NewScanner(os.Stdin)
	prompt()
	for reader.Scan() {
		text := reader.Text()
		if command, exists := commands[text]; exists {
			command.(func())()
		} else if strings.EqualFold(".exit", text) {
			return
		} else {
			// Pass the command to the parser
			handle(text)
		}
		prompt()
	}
	fmt.Println()
}

func handle(text string) {
	resp, err := http.Get(text)
	if err != nil {
		fmt.Println("Error fetching URL:", err)
		return
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}
	sha256Hash := sha256.Sum256(raw)
	fmt.Printf("SHA256: %x\n", sha256Hash)
}
