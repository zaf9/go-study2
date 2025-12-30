/**
 * T110: 测试历史记录卡片显示
 */

import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import QuizHistoryCard from "@/components/quiz/QuizHistoryCard";
import { QuizHistoryItem } from "@/types/quiz";

// Mock Next.js router
jest.mock("next/navigation", () => ({
  useRouter: () => ({
    push: jest.fn(),
  }),
}));

const mockHistoryItem: QuizHistoryItem = {
  id: 1,
  sessionId: "sess-001",
  topic: "constants",
  chapter: "iota",
  score: 80,
  totalQuestions: 5,
  correctAnswers: 4,
  passed: true,
  completedAt: "2025-12-30T10:30:00Z",
};

describe("QuizHistoryCard", () => {
  it("应显示章节名称", () => {
    render(<QuizHistoryCard item={mockHistoryItem} />);

    expect(screen.getByText("iota")).toBeInTheDocument();
  });

  it("应显示主题标签", () => {
    render(<QuizHistoryCard item={mockHistoryItem} />);

    expect(screen.getByText("constants")).toBeInTheDocument();
  });

  it("应显示得分", () => {
    render(<QuizHistoryCard item={mockHistoryItem} />);

    expect(screen.getByText("80分")).toBeInTheDocument();
  });

  it("应显示正确率", () => {
    render(<QuizHistoryCard item={mockHistoryItem} />);

    // 正确率 = 4/5 = 80%
    expect(screen.getByText("80%")).toBeInTheDocument();
  });

  it("应显示题目数量", () => {
    render(<QuizHistoryCard item={mockHistoryItem} />);

    expect(screen.getByText(/5题/)).toBeInTheDocument();
  });

  it("应显示完成时间", () => {
    render(<QuizHistoryCard item={mockHistoryItem} />);

    // 验证显示相对时间
    expect(screen.getByText(/前|ago/i)).toBeInTheDocument();
  });

  it("通过的测验应显示绿色通过标识", () => {
    render(<QuizHistoryCard item={mockHistoryItem} />);

    const passedBadge = screen.getByText(/已通过/);
    expect(passedBadge).toBeInTheDocument();
    // Ant Design Tag组件使用className "quiz-passed-badge"
    expect(passedBadge.closest('.ant-tag')).toHaveClass('quiz-passed-badge');
  });

  it("未通过的测验应显示红色未通过标识", () => {
    const failedItem: QuizHistoryItem = {
      ...mockHistoryItem,
      score: 50,
      correctAnswers: 2,
      passed: false,
    };

    render(<QuizHistoryCard item={failedItem} />);

    const failedBadge = screen.getByText(/未通过/);
    expect(failedBadge).toBeInTheDocument();
    // Ant Design Tag组件使用className "quiz-failed-badge"
    expect(failedBadge.closest('.ant-tag')).toHaveClass('quiz-failed-badge');
  });

  it("应包含正确的data-testid属性", () => {
    render(<QuizHistoryCard item={mockHistoryItem} />);

    const card = screen.getByTestId("quiz-history-card");
    expect(card).toBeInTheDocument();
  });

  it("点击卡片应触发onClick回调", async () => {
    const handleClick = jest.fn();
    const user = userEvent.setup();

    render(<QuizHistoryCard item={mockHistoryItem} onClick={handleClick} />);

    const card = screen.getByTestId("quiz-history-card");
    await user.click(card);

    expect(handleClick).toHaveBeenCalledWith("sess-001");
  });

  it("应支持hover效果", () => {
    render(<QuizHistoryCard item={mockHistoryItem} />);

    const card = screen.getByTestId("quiz-history-card");
    expect(card).toHaveClass('cursor-pointer');
  });

  it("缺少章节信息时应使用默认显示", () => {
    const itemWithoutChapter: QuizHistoryItem = {
      ...mockHistoryItem,
      chapter: null,
    };

    render(<QuizHistoryCard item={itemWithoutChapter} />);

    expect(screen.getByText(/测验/)).toBeInTheDocument();
  });

  it("缺少完成时间时应不显示时间", () => {
    const itemWithoutTime: QuizHistoryItem = {
      ...mockHistoryItem,
      completedAt: null,
    };

    render(<QuizHistoryCard item={itemWithoutTime} />);

    // 应该不会抛出错误
    expect(screen.getByTestId("quiz-history-card")).toBeInTheDocument();
    // 没有时间文本时不应该显示时间元素
    expect(screen.queryByText(/前|ago/i)).not.toBeInTheDocument();
  });

  it("满分测验应有特殊标识", () => {
    const perfectItem: QuizHistoryItem = {
      ...mockHistoryItem,
      score: 100,
      correctAnswers: 5,
      totalQuestions: 5,
    };

    render(<QuizHistoryCard item={perfectItem} />);

    expect(screen.getByText("100分")).toBeInTheDocument();
    expect(screen.getByText(/完美通过/i)).toBeInTheDocument();
  });

  it("应正确处理多选题的正确率计算", () => {
    const multiChoiceItem: QuizHistoryItem = {
      ...mockHistoryItem,
      totalQuestions: 10,
      correctAnswers: 7,
      score: 70,
    };

    render(<QuizHistoryCard item={multiChoiceItem} />);

    // 正确率 = 7/10 = 70%
    expect(screen.getByText("70%")).toBeInTheDocument();
  });
});
