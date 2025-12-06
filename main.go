package main

import (
	"bufio"
	"fmt"
	"go-study2/internal/app/lexical_elements"
	"io"
	"os"
	"sort"
	"strings"
)

// App represents the application with its I/O streams and menu configuration.
type App struct {
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
	menu   map[string]MenuItem
}

// MenuItem represents a single menu option.
type MenuItem struct {
	Description string
	Action      func()
}

// NewApp creates a new App instance with configured menu items.
// To add a new learning module:
//  1. Import the package.
//  2. Add a new entry to the menu map below.
//     Key: The menu option (e.g., "1").
//     Description: The text to display in the menu.
//     Action: The function to call when selected.
func NewApp(stdin io.Reader, stdout, stderr io.Writer) *App {
	return &App{
		stdin:  stdin,
		stdout: stdout,
		stderr: stderr,
		menu: map[string]MenuItem{
			"0": {
				Description: "Lexical elements",
				Action:      lexical_elements.Display,
			},
			// Add new items here
		},
	}
}

// Run starts the application's main loop.
func (a *App) Run() {
	reader := bufio.NewReader(a.stdin)

	for {
		fmt.Fprintln(a.stdout, "\nGo Lexical Elements Learning Tool")
		fmt.Fprintln(a.stdout, "---------------------------------")
		fmt.Fprintln(a.stdout, "Please select a topic to study:")

		// Sort keys for consistent display order
		var keys []string
		for k := range a.menu {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			fmt.Fprintf(a.stdout, "%s. %s\n", k, a.menu[k].Description)
		}
		fmt.Fprintln(a.stdout, "q. Quit")
		fmt.Fprint(a.stdout, "\nEnter your choice: ")

		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintf(a.stderr, "Error reading input: %v\n", err)
			return
		}
		choice := strings.TrimSpace(input)

		if choice == "q" {
			fmt.Fprintln(a.stdout, "Goodbye!")
			return
		}

		if item, ok := a.menu[choice]; ok {
			item.Action()
		} else {
			fmt.Fprintln(a.stdout, "Invalid choice, please try again.")
		}
	}
}

func main() {
	app := NewApp(os.Stdin, os.Stdout, os.Stderr)
	app.Run()
}
