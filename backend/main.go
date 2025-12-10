package main

import (
	"bufio"
	"flag"
	"fmt"
	"go-study2/internal/app/constants"
	"go-study2/internal/app/http_server"
	"go-study2/internal/app/lexical_elements"
	"go-study2/internal/config"
	typescli "go-study2/src/learning/types/cli"
	varcli "go-study2/src/learning/variables/cli"
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

// MenuItem 表示一个菜单选项。
// Action 函数接收三个 I/O 流参数：
//   - stdin: 用于读取用户输入
//   - stdout: 用于输出正常信息
//   - stderr: 用于输出错误信息
//
// 这种设计使得菜单动作可以是交互式的，例如显示子菜单并读取用户选择。
type MenuItem struct {
	Description string
	Action      func(io.Reader, io.Writer, io.Writer)
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
				Action:      lexical_elements.DisplayMenu,
			},
			"1": {
				Description: "Constants",
				Action:      constants.DisplayMenu,
			},
			"2": {
				Description: "Variables",
				Action:      varcli.DisplayMenu,
			},
			"3": {
				Description: "Types",
				Action:      typescli.DisplayMenu,
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
			item.Action(a.stdin, a.stdout, a.stderr)
		} else {
			fmt.Fprintln(a.stdout, "Invalid choice, please try again.")
		}
	}
}

func main() {
	daemon := flag.Bool("d", false, "Run in daemon/HTTP mode")
	flag.BoolVar(daemon, "daemon", false, "Run in daemon/HTTP mode")
	flag.Parse()

	if *daemon {
		runHttpServer()
	} else {
		app := NewApp(os.Stdin, os.Stdout, os.Stderr)
		app.Run()
	}
}

func runHttpServer() {
	// 加载配置 (Load 内部会自动读取默认配置文件并进行验证)
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// 初始化服务器
	s, err := http_server.NewServer(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to init server: %v\n", err)
		os.Exit(1)
	}

	// 启动服务器 (Run 会阻塞直到收到停止信号)
	s.Run()
}
