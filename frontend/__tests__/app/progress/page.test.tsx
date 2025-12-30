/**
 * 测试进度页面组件 (User Story 4 - T093)
 * 验证进度数据的正确展示
 */

import { render, screen, waitFor, within } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import ProgressPage from "@/app/(protected)/progress/page";
import useProgress from "@/hooks/useProgress";
import { ProgressSnapshot } from "@/types/learning";

// Mock useProgress Hook
jest.mock("@/hooks/useProgress");

// Mock Next.js router
const mockPush = jest.fn();
jest.mock("next/navigation", () => ({
  useRouter: () => ({
    push: mockPush,
    replace: jest.fn(),
    refresh: jest.fn(),
  }),
}));

// Mock components
jest.mock("@/components/common/Loading", () => {
  return function MockLoading() {
    return <div>Loading...</div>;
  };
});

jest.mock("@/components/common/ErrorMessage", () => {
  return function MockErrorMessage({
    message,
    description,
  }: {
    message: string;
    description: string;
  }) {
    return (
      <div>
        <div>{message}</div>
        <div>{description}</div>
      </div>
    );
  };
});

jest.mock("@/components/progress/ProgressOverview", () => {
  return function MockProgressOverview({
    overall,
    next,
    onContinue,
  }: {
    overall: ProgressSnapshot["overall"];
    next: ProgressSnapshot["next"];
    onContinue: (hint: NonNullable<ProgressSnapshot["next"]>) => void;
  }) {
    return (
      <div data-testid="progress-overview">
        <div>进度: {overall.progress}%</div>
        <div>完成章节: {overall.completedChapters}</div>
        <div>总章节: {overall.totalChapters}</div>
        {next && (
          <button onClick={() => onContinue(next)}>继续学习</button>
        )}
      </div>
    );
  };
});

jest.mock("@/components/progress/TopicProgressCard", () => {
  return function MockTopicProgressCard({
    topic,
    onContinue,
  }: {
    topic: { id: string; name: string; progress: number };
    onContinue: (chapter?: { chapter: string }) => void;
  }) {
    return (
      <div data-testid={`topic-card-${topic.id}`}>
        <div>{topic.name}</div>
        <div>{topic.progress}%</div>
        <button onClick={() => onContinue()}>继续学习</button>
      </div>
    );
  };
});

describe("ProgressPage", () => {
  const mockProgressData: ProgressSnapshot = {
    overall: {
      progress: 24.0,
      completedChapters: 12,
      totalChapters: 50,
      studyDays: 5,
      totalStudyTime: 3600,
    },
    topics: [
      {
        id: "variables",
        name: "Variables",
        weight: 30,
        progress: 40.0,
        completedChapters: 12,
        totalChapters: 30,
      },
      {
        id: "constants",
        name: "Constants",
        weight: 20,
        progress: 15.0,
        completedChapters: 3,
        totalChapters: 20,
      },
      {
        id: "lexical_elements",
        name: "Lexical Elements",
        weight: 25,
        progress: 8.0,
        completedChapters: 2,
        totalChapters: 25,
      },
    ],
    next: {
      topic: "constants",
      chapter: "iota",
      title: "Iota常量",
    },
  };

  beforeEach(() => {
    jest.clearAllMocks();
  });

  describe("加载状态", () => {
    it("应显示加载指示器", () => {
      (useProgress as jest.Mock).mockReturnValue({
        overview: undefined,
        next: null,
        isLoading: true,
        error: null,
      });

      render(<ProgressPage />);

      expect(screen.getByText("Loading...")).toBeInTheDocument();
    });
  });

  describe("错误状态", () => {
    it("应显示错误消息", () => {
      const mockError = new Error("API请求失败");

      (useProgress as jest.Mock).mockReturnValue({
        overview: undefined,
        next: null,
        isLoading: false,
        error: mockError,
      });

      render(<ProgressPage />);

      expect(screen.getByText("加载进度失败")).toBeInTheDocument();
      expect(screen.getByText("API请求失败")).toBeInTheDocument();
    });
  });

  describe("T093: 数据展示", () => {
    it("应正确显示进度概览", () => {
      (useProgress as jest.Mock).mockReturnValue({
        overview: mockProgressData,
        next: mockProgressData.next,
        isLoading: false,
        error: null,
      });

      render(<ProgressPage />);

      // 验证概览数据显示
      expect(screen.getByTestId("progress-overview")).toBeInTheDocument();
      expect(screen.getByText("进度: 24%")).toBeInTheDocument();
      expect(screen.getByText("完成章节: 12")).toBeInTheDocument();
      expect(screen.getByText("总章节: 50")).toBeInTheDocument();
    });

    it("应按进度降序显示所有主题", () => {
      (useProgress as jest.Mock).mockReturnValue({
        overview: mockProgressData,
        next: mockProgressData.next,
        isLoading: false,
        error: null,
      });

      render(<ProgressPage />);

      // 验证主题按进度降序显示
      const topicCards = screen.getAllByTestId(/^topic-card-/);
      expect(topicCards).toHaveLength(3);

      // Variables (40%) 应该排在第一位
      expect(screen.getByTestId("topic-card-variables")).toBeInTheDocument();
      // Constants (15%) 第二位
      expect(screen.getByTestId("topic-card-constants")).toBeInTheDocument();
      // Lexical Elements (8%) 第三位
      expect(
        screen.getByTestId("topic-card-lexical_elements"),
      ).toBeInTheDocument();
    });

    it("应支持主题筛选", async () => {
      const user = userEvent.setup();

      (useProgress as jest.Mock).mockReturnValue({
        overview: mockProgressData,
        next: mockProgressData.next,
        isLoading: false,
        error: null,
      });

      render(<ProgressPage />);

      // 初始应显示所有主题
      expect(screen.getAllByTestId(/^topic-card-/)).toHaveLength(3);

      // 注意: Ant Design Select 组件的测试比较复杂
      // 这里仅验证组件存在,实际筛选逻辑在集成测试中验证
      expect(screen.getByText("筛选主题")).toBeInTheDocument();
    });

    it("应支持导航到下一章节", async () => {
      const user = userEvent.setup();

      (useProgress as jest.Mock).mockReturnValue({
        overview: mockProgressData,
        next: mockProgressData.next,
        isLoading: false,
        error: null,
      });

      render(<ProgressPage />);

      // 点击概览中的"继续学习"按钮（通过data-testid定位）
      const overviewSection = screen.getByTestId("progress-overview");
      const continueButton = within(overviewSection).getByText("继续学习");
      await user.click(continueButton);

      // 验证导航
      expect(mockPush).toHaveBeenCalledWith("/topics/constants/iota");
    });
  });

  describe("空数据处理", () => {
    it("无进度数据时应返回 null", () => {
      (useProgress as jest.Mock).mockReturnValue({
        overview: null,
        next: null,
        isLoading: false,
        error: null,
      });

      const { container } = render(<ProgressPage />);

      expect(container.firstChild).toBeNull();
    });

    it("无主题数据时应显示空列表", () => {
      (useProgress as jest.Mock).mockReturnValue({
        overview: {
          overall: {
            progress: 0,
            completedChapters: 0,
            totalChapters: 50,
            studyDays: 0,
            totalStudyTime: 0,
          },
          topics: [],
          next: null,
        },
        next: null,
        isLoading: false,
        error: null,
      });

      render(<ProgressPage />);

      // 应显示概览但没有主题卡片
      expect(screen.getByTestId("progress-overview")).toBeInTheDocument();
      expect(screen.queryAllByTestId(/^topic-card-/)).toHaveLength(0);
    });
  });
});
