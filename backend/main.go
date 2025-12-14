package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"go-study2/internal/app/constants"
	"go-study2/internal/app/http_server"
	"go-study2/internal/app/lexical_elements"
	"go-study2/internal/config"
	"go-study2/internal/domain/quiz"
	"go-study2/internal/infrastructure/database"

	logger "go-study2/internal/infrastructure/logger"
	appjwt "go-study2/internal/pkg/jwt"
	typescli "go-study2/src/learning/types/cli"
	varcli "go-study2/src/learning/variables/cli"
	"io"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/gogf/gf/v2/database/gdb"

	"github.com/gogf/gf/v2/os/gctx"
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

	// 初始化日志系统 (从 backend/configs 中加载 logger 配置)
	lcfg, err := logger.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load logger config: %v\n", err)
		os.Exit(1)
	}
	if err := logger.Initialize(lcfg); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}

	// 初始化数据库
	ctx := gctx.New()
	if _, err = database.Init(ctx, cfg.Database); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to init database: %v\n", err)
		os.Exit(1)
	}

	// 尝试从 quiz_data 加载题库并在空表时写入数据库（用于本地/首次启动）
	// 仅在配置中指定了 quiz.dataPath 时执行
	if cfg.Progress.Quiz.DataPath != "" {
		// 从文件系统加载 YAML 题库到临时内存仓储
		repo := quiz.NewRepository()
		if err := quiz.LoadAllBanks(cfg.Progress.Quiz.DataPath, repo); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to load quiz banks: %v\n", err)
		} else {
			// 如果数据库中尚无 quiz_questions，则将 YAML 中的题目作为种子插入
			db := database.Default()
			if db != nil {
				// 按 topic/chapter 做增量插入：仅将文件中数据库不存在的题目写入，避免重复导入
				var seed []gdb.Map
				for _, topic := range repo.Topics() {
					for _, ch := range repo.Chapters(topic) {
						qs, _ := repo.GetBank(topic, ch)

						// 查询已有题干用于去重
						existing := map[string]struct{}{}
						if recs, _ := db.Model("quiz_questions").Fields("question").Where("topic", topic).Where("chapter", ch).All(ctx); recs != nil {
							for _, r := range recs.List() {
								// recs.List() typically returns []gdb.Map; handle that directly
								var qtext string
								if v, found := r["question"]; found {
									if str, ok := v.(string); ok {
										qtext = strings.TrimSpace(str)
									}
								} else {
									// fallback: try asserting to an interface with Map() method
									type hasMap interface{ Map() gdb.Map }
									if recWithMap, ok2 := any(r).(hasMap); ok2 {
										if v2, found2 := recWithMap.Map()["question"]; found2 {
											if str, ok := v2.(string); ok {
												qtext = strings.TrimSpace(str)
											}
										}
									}
								}
								existing[qtext] = struct{}{}
							}
						}

						for _, q := range qs {
							stem := strings.TrimSpace(q.Stem)
							if stem == "" {
								continue
							}
							if _, ok := existing[stem]; ok {
								continue
							}
							// 将答案字符串拆分为单字母数组并编码为 JSON
							ansArr := []string{}
							for _, r := range q.Answer {
								ansArr = append(ansArr, strings.ToUpper(string(r)))
							}
							ansJSON, _ := json.Marshal(ansArr)
							optionsJSON, _ := json.Marshal(q.Options)
							seed = append(seed, gdb.Map{
								"topic":           topic,
								"chapter":         ch,
								"type":            q.Type,
								"difficulty":      q.Difficulty,
								"question":        stem,
								"options":         string(optionsJSON),
								"correct_answers": string(ansJSON),
								"explanation":     q.Explanation,
							})
						}
					}
				}
				if len(seed) > 0 {
					_, _ = db.Model("quiz_questions").Data(seed).Insert()
				}
			}
		}
	}

	// 配置 JWT
	_ = appjwt.Configure(appjwt.Options{
		Secret:             cfg.Jwt.Secret,
		Issuer:             cfg.Jwt.Issuer,
		AccessTokenExpiry:  time.Duration(cfg.Jwt.AccessTokenExpiry) * time.Second,
		RefreshTokenExpiry: time.Duration(cfg.Jwt.RefreshTokenExpiry) * time.Second,
	})

	// 初始化服务器
	s, err := http_server.NewServer(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to init server: %v\n", err)
		os.Exit(1)
	}

	// 启动服务器 (Run 会阻塞直到收到停止信号)
	s.Run()
}
