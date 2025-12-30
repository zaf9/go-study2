/**
 * T109: 测试测验中心组件渲染
 */

import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import QuizCenter from "@/components/quiz/QuizCenter";
import { QuizHistoryItem } from "@/types/quiz";
import * as quizService from "@/services/quizService";

// Mock Next.js router
jest.mock("next/navigation", () => ({
  useRouter: () => ({
    push: jest.fn(),
  }),
}));

// Mock quizService
jest.mock("@/services/quizService", () => ({
  quizService: {
    getQuizHistory: jest.fn(),
  },
}));

const mockQuizHistory: QuizHistoryItem[] = [
  {
    id: 1,
    sessionId: "sess-001",
    topic: "constants",
    chapter: "iota",
    score: 80,
    totalQuestions: 5,
    correctAnswers: 4,
    passed: true,
    completedAt: "2025-12-30T10:00:00Z",
  },
  {
    id: 2,
    sessionId: "sess-002",
    topic: "variables",
    chapter: "declaration",
    score: 100,
    totalQuestions: 6,
    correctAnswers: 6,
    passed: true,
    completedAt: "2025-12-30T11:00:00Z",
  },
  {
    id: 3,
    sessionId: "sess-003",
    topic: "constants",
    chapter: "types",
    score: 50,
    totalQuestions: 4,
    correctAnswers: 2,
    passed: false,
    completedAt: "2025-12-30T12:00:00Z",
  },
];

describe("QuizCenter", () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  describe("加载状态", () => {
    it("应显示加载指示器", () => {
      (quizService.quizService.getQuizHistory as jest.Mock).mockImplementation(
        () => new Promise(() => {}) // Never resolves
      );

      render(<QuizCenter />);

      expect(screen.getByTestId("loading-spinner")).toBeInTheDocument();
    });
  });

  describe("错误状态", () => {
    it("应显示错误消息", async () => {
      (quizService.quizService.getQuizHistory as jest.Mock).mockRejectedValue(
        new Error("Failed to load history")
      );

      render(<QuizCenter />);

      await waitFor(() => {
        expect(screen.getByTestId("error-message")).toBeInTheDocument();
      });
    });
  });

  describe("数据展示", () => {
    it("应正确显示测验历史列表", async () => {
      (quizService.quizService.getQuizHistory as jest.Mock).mockResolvedValue(
        mockQuizHistory
      );

      render(<QuizCenter />);

      await waitFor(() => {
        expect(screen.getByTestId("quiz-center-component")).toBeInTheDocument();
      });

      // 验证显示章节信息
      expect(screen.getByText("iota")).toBeInTheDocument();
      expect(screen.getByText("declaration")).toBeInTheDocument();
      expect(screen.getByText("types")).toBeInTheDocument();
    });
  });

  describe("空数据处理", () => {
    it("应显示空状态提示", async () => {
      (quizService.quizService.getQuizHistory as jest.Mock).mockResolvedValue([]);

      render(<QuizCenter />);

      await waitFor(() => {
        expect(screen.getByTestId("empty-state")).toBeInTheDocument();
      });

      expect(screen.getByText(/暂无测验记录/i)).toBeInTheDocument();
      expect(screen.getByText(/开始学习/i)).toBeInTheDocument();
    });
  });
});
