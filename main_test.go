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

// TestApp_LexicalElements 已废弃，因为 DisplayMenu 不再直接调用所有 Display 函数
// 新的测试用例请使用 TestApp_NavigateToSubMenu 和 TestApp_ReturnFromSubMenu

func TestApp_Extensibility(t *testing.T) {
	stdin := strings.NewReader("9\nq\n")
	stdout := &bytes.Buffer{}
	app := NewApp(stdin, stdout, &bytes.Buffer{})

	// Add a custom menu item dynamically to test extensibility
	called := false
	app.menu["9"] = MenuItem{
		Description: "Custom Test Module",
		Action: func(io.Reader, io.Writer, io.Writer) {
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

// TestApp_NavigateToSubMenu 测试导航到子菜单的功能。
// 输入 "0" 后，应该显示子菜单，包含选项 0-10 和 'q'。
func TestApp_NavigateToSubMenu(t *testing.T) {
	stdin := strings.NewReader("0\nq\nq\n")
	stdout := &bytes.Buffer{}
	app := NewApp(stdin, stdout, &bytes.Buffer{})

	app.Run()

	output := stdout.String()

	// 验证子菜单显示
	if !strings.Contains(output, "词法元素学习菜单") {
		t.Error("子菜单标题未显示")
	}

	// 验证子菜单选项 0-10
	for i := 0; i <= 10; i++ {
		expectedOption := fmt.Sprintf("%d.", i)
		if !strings.Contains(output, expectedOption) {
			t.Errorf("子菜单选项 %d 未显示", i)
		}
	}

	// 验证返回选项
	if !strings.Contains(output, "q. 返回上级菜单") {
		t.Error("返回上级菜单选项未显示")
	}
}

// TestApp_ReturnFromSubMenu 测试从子菜单返回主菜单的功能。
// 输入 "0" 进入子菜单，然后输入 "q" 返回主菜单，最后输入 "q" 退出。
func TestApp_ReturnFromSubMenu(t *testing.T) {
	stdin := strings.NewReader("0\nq\nq\n")
	stdout := &bytes.Buffer{}
	app := NewApp(stdin, stdout, &bytes.Buffer{})

	app.Run()

	output := stdout.String()

	// 验证进入了子菜单
	if !strings.Contains(output, "词法元素学习菜单") {
		t.Error("未进入子菜单")
	}

	// 验证返回了主菜单（主菜单应该再次显示）
	menuCount := strings.Count(output, "Go Lexical Elements Learning Tool")
	if menuCount < 2 {
		t.Errorf("主菜单应该显示至少 2 次（进入前和返回后），实际显示 %d 次", menuCount)
	}

	// 验证退出消息
	if !strings.Contains(output, "Goodbye!") {
		t.Error("退出消息未显示")
	}
}
