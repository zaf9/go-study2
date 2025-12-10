# CLI Contract for Go Lexical Elements Learning Tool

## Overview

This document defines the command-line interface (CLI) contract for the Go Lexical Elements Learning Tool. The tool operates as a menu-driven application, providing options to explore different lexical elements of the Go language.

## Commands

### Main Entry Point

The application is started by running the `main.go` file from `backend/`.

**Command**: `cd backend && go run main.go`

**Description**: Launches the interactive menu for the learning tool.

### Menu Options

Once launched, the user is presented with a menu. Input is expected via standard input (stdin).

#### Option: Lexical Elements

**Input**: `0`

**Description**: Executes the code examples and explanations related to the "Lexical elements" package. The specific output depends on the implementation within the `lexical_elements` package, but it is expected to print educational content to standard output (stdout).

#### Option: Quit

**Input**: `q`

**Description**: Exits the application gracefully.

## Input/Output

-   **Input**: User input is read from standard input (stdin) as single characters (`0`, `q`, etc.).
-   **Output**: Informational messages, menu prompts, code examples, and explanations are printed to standard output (stdout). Errors are printed to standard error (stderr).
