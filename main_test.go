package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
)

func TestApp_Quit(t *testing.T) {
	stdin := strings.NewReader("q\n")
	stdout := &bytes.Buffer{}
	app := NewApp(stdin, stdout, &bytes.Buffer{})

	app.Run()

	expectedOutput := "Goodbye!"
	if !strings.Contains(stdout.String(), expectedOutput) {
		t.Errorf("expected output to contain %q, but got %q", expectedOutput, stdout.String())
	}
}

func TestApp_InvalidInput(t *testing.T) {
	// Simulate user entering an invalid command, then quitting.
	stdin := strings.NewReader("x\nq\n")
	stdout := &bytes.Buffer{}
	app := NewApp(stdin, stdout, &bytes.Buffer{})

	app.Run()

	expectedErrorMsg := "Invalid choice"
	if !strings.Contains(stdout.String(), expectedErrorMsg) {
		t.Errorf("expected output to contain %q for invalid input, but it did not. Got: %q", expectedErrorMsg, stdout.String())
	}
}

func TestApp_LexicalElements(t *testing.T) {
	// Capture os.Stdout because lexical_elements package uses fmt.Println directly
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Read captured output in a goroutine to prevent deadlock
	var buf bytes.Buffer
	done := make(chan bool)
	go func() {
		io.Copy(&buf, r)
		done <- true
	}()

	// Simulate user entering '0' (Lexical Elements), then quitting.
	stdin := strings.NewReader("0\nq\n")
	stdout := &bytes.Buffer{} // This captures menu output from App
	app := NewApp(stdin, stdout, &bytes.Buffer{})

	app.Run()

	// Close writer and restore stdout
	w.Close()
	os.Stdout = oldStdout

	// Wait for reading to complete
	<-done
	capturedOutput := buf.String()

	// Check for some expected output from the lexical elements module
	// We check for "Go 语言的注释" which is printed by DisplayComments (the first called function)
	expectedContent := "Go 语言的注释"
	if !strings.Contains(capturedOutput, expectedContent) {
		t.Errorf("expected output to contain %q when selecting option '0', but got output length %d", expectedContent, len(capturedOutput))
	}
}

func TestApp_Extensibility(t *testing.T) {
	stdin := strings.NewReader("9\nq\n")
	stdout := &bytes.Buffer{}
	app := NewApp(stdin, stdout, &bytes.Buffer{})

	// Add a custom menu item dynamically to test extensibility
	called := false
	app.menu["9"] = MenuItem{
		Description: "Custom Test Module",
		Action: func() {
			called = true
			fmt.Println("Custom module executed")
		},
	}

	// We need to capture os.Stdout for the custom action too if it uses fmt.Println
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Read captured output in a goroutine
	var buf bytes.Buffer
	done := make(chan bool)
	go func() {
		io.Copy(&buf, r)
		done <- true
	}()

	app.Run()

	w.Close()
	os.Stdout = oldStdout
	<-done
	capturedOutput := buf.String()

	if !called {
		t.Error("Custom menu item action was not executed")
	}

	if !strings.Contains(capturedOutput, "Custom module executed") {
		t.Error("Custom menu item output not captured")
	}

	// Check if the menu description was printed to stdout (App.stdout)
	if !strings.Contains(stdout.String(), "9. Custom Test Module") {
		t.Error("Custom menu item description not found in menu output")
	}
}
