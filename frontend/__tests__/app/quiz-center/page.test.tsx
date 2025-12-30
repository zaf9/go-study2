/**
 * T111: 测试测验中心页面
 */

import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import QuizCenterPage from "@/app/(protected)/quiz-center/page";

// Mock QuizCenter component
jest.mock("@/components/quiz/QuizCenter", () => {
  return function MockQuizCenter() {
    return <div data-testid="quiz-center-component">Quiz Center Component</div>;
  };
});

// Mock Next.js router
jest.mock("next/navigation", () => ({
  useRouter: () => ({
    push: jest.fn(),
    replace: jest.fn(),
  }),
  usePathname: () => "/quiz-center",
}));

// Mock auth context
jest.mock("@/contexts/AuthContext", () => ({
  useAuth: () => ({
    user: { id: 1, username: "testuser" },
    isAuthenticated: true,
  }),
}));

describe("QuizCenterPage", () => {
  it("应渲染页面标题", () => {
    render(<QuizCenterPage />);

    expect(screen.getByText("测验中心")).toBeInTheDocument();
  });

  it("应渲染QuizCenter组件", () => {
    render(<QuizCenterPage />);

    expect(screen.getByTestId("quiz-center-component")).toBeInTheDocument();
  });

  it("应显示页面描述", () => {
    render(<QuizCenterPage />);

    expect(
      screen.getByText(/查看和管理你的测验记录/i)
    ).toBeInTheDocument();
  });

  it("应包含返回按钮", () => {
    render(<QuizCenterPage />);

    const backButton = screen.getByRole("button", { name: /返回|主页/i });
    expect(backButton).toBeInTheDocument();
  });

  it("点击返回按钮应导航到主页", async () => {
    const { useRouter } = require("next/navigation");
    const mockPush = jest.fn();
    useRouter.mockReturnValue({ push: mockPush });

    const user = userEvent.setup();
    render(<QuizCenterPage />);

    const backButton = screen.getByRole("button", { name: /返回|主页/i });
    await user.click(backButton);

    expect(mockPush).toHaveBeenCalledWith("/dashboard");
  });

  it("应设置正确的页面元数据", () => {
    // 验证页面有正确的布局和样式
    render(<QuizCenterPage />);

    const pageContainer = screen.getByTestId("quiz-center-page");
    expect(pageContainer).toBeInTheDocument();
    expect(pageContainer).toHaveClass(/container|page|quiz-center/);
  });

  it("应显示统计摘要区域", () => {
    render(<QuizCenterPage />);

    expect(screen.getByTestId("stats-summary")).toBeInTheDocument();
  });

  it("应支持响应式布局", () => {
    render(<QuizCenterPage />);

    const pageContainer = screen.getByTestId("quiz-center-page");
    // 验证有响应式相关的class
    expect(pageContainer.className).toMatch(/responsive|md:|lg:|sm:/);
  });

  it("未认证用户应被重定向", () => {
    const { useAuth } = require("@/contexts/AuthContext");
    useAuth.mockReturnValue({
      user: null,
      isAuthenticated: false,
    });

    const { useRouter } = require("next/navigation");
    const mockReplace = jest.fn();
    useRouter.mockReturnValue({ replace: mockReplace });

    render(<QuizCenterPage />);

    // 验证重定向到登录页
    expect(mockReplace).toHaveBeenCalledWith("/login");
  });

  it("应包含面包屑导航", () => {
    render(<QuizCenterPage />);

    expect(screen.getByRole("navigation", { name: /breadcrumb/i })).toBeInTheDocument();
    expect(screen.getByText("主页")).toBeInTheDocument();
    expect(screen.getByText("测验中心")).toBeInTheDocument();
  });

  it("应支持键盘导航", async () => {
    const user = userEvent.setup();
    render(<QuizCenterPage />);

    // Tab键应该能够聚焦到可交互元素
    await user.tab();
    
    const focusedElement = document.activeElement;
    expect(focusedElement).toBeInTheDocument();
  });
});
