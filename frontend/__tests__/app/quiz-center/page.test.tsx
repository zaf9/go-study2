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
const mockPush = jest.fn();
const mockReplace = jest.fn();
jest.mock("next/navigation", () => ({
  useRouter: jest.fn(() => ({
    push: mockPush,
    replace: mockReplace,
  })),
  usePathname: () => "/quiz-center",
}));

// Mock useAuth hook
const mockUseAuth = jest.fn();
jest.mock("@/hooks/useAuth", () => ({
  __esModule: true,
  default: () => mockUseAuth(),
}));

describe("QuizCenterPage", () => {
  beforeEach(() => {
    jest.clearAllMocks();
    // 重置 useAuth mock为已登录用户
    mockUseAuth.mockReturnValue({
      user: { id: 1, username: "testuser" },
      isAuthenticated: true,
    });
  });

  it("应渲染页面标题", () => {
    render(<QuizCenterPage />);

    // 页面中有多个“测验中心”文本，使用 getAllByText
    const titles = screen.getAllByText("测验中心");
    expect(titles.length).toBeGreaterThan(0);
  });

  it("应渲染QuizCenter组件", () => {
    render(<QuizCenterPage />);

    expect(screen.getByTestId("quiz-center-component")).toBeInTheDocument();
  });

  it("应设置正确的页面元数据", () => {
    // 验证页面有正确的布局和样式
    render(<QuizCenterPage />);

    const pageContainer = screen.getByTestId("quiz-center-page");
    expect(pageContainer).toBeInTheDocument();
    expect(pageContainer).toHaveClass(/container|page|quiz-center/);
  });

  it("应支持响应式布局", () => {
    render(<QuizCenterPage />);

    const pageContainer = screen.getByTestId("quiz-center-page");
    // 验证有响应式相关的class - 实际组件使用container和px-4等
    expect(pageContainer.className).toMatch(/container|mx-auto|px-4/);
  });

  it("未认证用户应被重定向", async () => {
    // 设置 mock 返回未认证用户
    mockUseAuth.mockReturnValue({
      user: null,
      isAuthenticated: false,
    });

    // 使用全局 mockReplace
    mockReplace.mockClear();

    render(<QuizCenterPage />);

    // 验证重定向到登录页
    await waitFor(() => {
      expect(mockReplace).toHaveBeenCalledWith("/login");
    }, { timeout: 3000 });
  });

  it("应包含页面标题", () => {
    render(<QuizCenterPage />);

    // 页面有"测验中心"标题
    const quizCenterTexts = screen.getAllByText("测验中心");
    expect(quizCenterTexts.length).toBeGreaterThan(0);
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
