# Quickstart: Go Lexical Elements Learning Tool

This document provides instructions on how to quickly get started with the Go Lexical Elements Learning Tool.

## Prerequisites

- Go 1.24.5 or later installed.
- Go Modules enabled (default for Go 1.11+).

## 1. Clone the Repository

If you haven't already, clone the project repository:

```bash
git clone <repository_url>
cd go-study2
```
(Note: Replace `<repository_url>` with the actual repository URL.)

## 2. Run the Application

Navigate to the project root directory and run the `main.go` file from `backend/`:

```bash
cd backend
go run main.go
```

### Expected Output

Upon running the application, you will see a menu similar to this:

```
Go Lexical Elements Learning Tool

Please select a topic to study:
0. Lexical elements
q. Quit

Enter your choice:
```

### 3. Interact with the Menu

-   **To study Lexical elements**: Enter `0` and press Enter. The application will then display examples and explanations for various lexical elements.
-   **To quit**: Enter `q` and press Enter. The application will exit.

## 4. (Optional) Explore the Code

The core logic and examples for "Lexical elements" are located in the `backend/internal/app/lexical_elements` package. Each sub-topic (e.g., Comments, Keywords) will have its own `.go` file containing examples and explanations.
