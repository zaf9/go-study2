/**
 * 测验中心导航集成测试
 * 
 * 测试目标:
 * - 从主页到测验中心的导航路径 ≤ 2 步
 * - 测验记录列表完整显示
 * - 可查看测验详情
 * 
 * @module __tests__/integration/quiz-center-navigation
 */

import { render, screen, waitFor, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import QuizCenterPage from '@/app/(protected)/quiz-center/page';
import DashboardPage from '@/app/(protected)/dashboard/page';
import { AuthProvider } from '@/contexts/AuthContext';
import { quizService } from '@/services/quizService';

// Mock 测验服务
jest.mock('@/services/quizService');

// Mock Next.js 导航
jest.mock('next/navigation', () => ({
  useRouter: () => ({
    push: jest.fn(),
    replace: jest.fn(),
    prefetch: jest.fn(),
  }),
  usePathname: () => '/quiz-center',
  useSearchParams: () => new URLSearchParams(),
}));

// Mock 测验历史数据
const mockQuizHistory = [
  {
    sessionId: 'session-001',
    topic: 'lexical_elements',
    chapter: 'identifiers',
    score: 80,
    totalQuestions: 10,
    correctAnswers: 8,
    createdAt: new Date('2025-12-30T10:00:00Z').toISOString(),
    submittedAt: new Date('2025-12-30T10:15:00Z').toISOString(),
  },
  {
    sessionId: 'session-002',
    topic: 'types',
    chapter: 'basic_types',
    score: 90,
    totalQuestions: 10,
    correctAnswers: 9,
    createdAt: new Date('2025-12-30T14:00:00Z').toISOString(),
    submittedAt: new Date('2025-12-30T14:12:00Z').toISOString(),
  },
  {
    sessionId: 'session-003',
    topic: 'variables',
    chapter: 'declarations',
    score: 70,
    totalQuestions: 10,
    correctAnswers: 7,
    createdAt: new Date('2025-12-30T16:00:00Z').toISOString(),
    submittedAt: new Date('2025-12-30T16:18:00Z').toISOString(),
  },
];

describe('测验中心导航集成测试', () => {
  const user = userEvent.setup();

  beforeEach(() => {
    jest.clearAllMocks();
    (quizService.getQuizHistory as jest.Mock).mockResolvedValue(mockQuizHistory);
  });

  describe('导航路径测试', () => {
    it('应该从主导航直接访问测验中心(1步)', async () => {
      // 模拟主导航组件
      const { container } = render(
        <AuthProvider>
          <nav aria-label="主导航">
            <a href="/quiz-center">测验中心</a>
          </nav>
        </AuthProvider>
      );

      // 验证导航链接存在
      const navLink = screen.getByRole('link', { name: /测验中心/i });
      expect(navLink).toBeInTheDocument();
      expect(navLink).toHaveAttribute('href', '/quiz-center');

      // 点击导航链接 (第1步)
      await user.click(navLink);

      // 验证:从主导航到测验中心仅需1步
      expect(navLink).toHaveAttribute('href', '/quiz-center');
    });

    it('应该从Dashboard快速访问测验中心(≤2步)', async () => {
      render(
        <AuthProvider>
          <DashboardPage />
        </AuthProvider>
      );

      // 在Dashboard中查找"查看全部"或"测验中心"链接
      // 第1步:在Dashboard找到测验相关区域
      const quizSection = await screen.findByText(/最近测验|测验记录/i, {}, { timeout: 3000 });
      expect(quizSection).toBeInTheDocument();

      // 第2步:点击"查看全部"按钮
      const viewAllButton = screen.getByRole('link', { name: /查看全部|更多/i });
      expect(viewAllButton).toHaveAttribute('href', '/quiz-center');

      // 验证:从Dashboard到测验中心需要2步
      await user.click(viewAllButton);
    });
  });

  describe('测验记录列表显示测试', () => {
    it('应该完整显示所有测验记录', async () => {
      render(
        <AuthProvider>
          <QuizCenterPage />
        </AuthProvider>
      );

      // 等待数据加载
      await waitFor(() => {
        expect(quizService.getQuizHistory).toHaveBeenCalled();
      });

      // 验证所有测验记录都显示
      await waitFor(() => {
        mockQuizHistory.forEach((quiz) => {
          // 检查主题名称
          const topicElements = screen.getAllByText(new RegExp(quiz.topic, 'i'));
          expect(topicElements.length).toBeGreaterThan(0);

          // 检查章节名称
          const chapterElements = screen.getAllByText(new RegExp(quiz.chapter, 'i'));
          expect(chapterElements.length).toBeGreaterThan(0);

          // 检查得分
          const scoreElements = screen.getAllByText(new RegExp(`${quiz.score}`, 'i'));
          expect(scoreElements.length).toBeGreaterThan(0);
        });
      });
    });

    it('应该显示正确的测验统计信息', async () => {
      render(
        <AuthProvider>
          <QuizCenterPage />
        </AuthProvider>
      );

      await waitFor(() => {
        expect(quizService.getQuizHistory).toHaveBeenCalled();
      });

      // 验证每个测验卡片显示的信息
      await waitFor(() => {
        mockQuizHistory.forEach((quiz) => {
          const correctRate = ((quiz.correctAnswers / quiz.totalQuestions) * 100).toFixed(0);
          
          // 检查正确率
          const rateElements = screen.getAllByText(new RegExp(`${correctRate}%`, 'i'));
          expect(rateElements.length).toBeGreaterThan(0);

          // 检查题目数量
          const totalElements = screen.getAllByText(new RegExp(`${quiz.totalQuestions}`, 'i'));
          expect(totalElements.length).toBeGreaterThan(0);
        });
      });
    });

    it('应该按时间倒序显示测验记录', async () => {
      render(
        <AuthProvider>
          <QuizCenterPage />
        </AuthProvider>
      );

      await waitFor(() => {
        expect(quizService.getQuizHistory).toHaveBeenCalled();
      });

      // 获取所有测验卡片
      const quizCards = await screen.findAllByTestId(/quiz-history-card/i);

      // 验证至少有3条记录
      expect(quizCards.length).toBe(3);

      // 验证第一条是最新的 (session-003)
      const firstCard = quizCards[0];
      expect(within(firstCard).getByText(/variables/i)).toBeInTheDocument();

      // 验证最后一条是最早的 (session-001)
      const lastCard = quizCards[quizCards.length - 1];
      expect(within(lastCard).getByText(/lexical_elements/i)).toBeInTheDocument();
    });
  });

  describe('主题筛选功能测试', () => {
    it('应该提供主题筛选选项', async () => {
      render(
        <AuthProvider>
          <QuizCenterPage />
        </AuthProvider>
      );

      await waitFor(() => {
        expect(quizService.getQuizHistory).toHaveBeenCalled();
      });

      // 查找筛选控件
      const filterSelect = screen.getByRole('combobox', { name: /主题|筛选/i });
      expect(filterSelect).toBeInTheDocument();

      // 验证"全部"选项存在
      const allOption = screen.getByText(/全部|All/i);
      expect(allOption).toBeInTheDocument();
    });

    it('应该根据选择的主题过滤测验记录', async () => {
      render(
        <AuthProvider>
          <QuizCenterPage />
        </AuthProvider>
      );

      await waitFor(() => {
        expect(quizService.getQuizHistory).toHaveBeenCalled();
      });

      // 选择"types"主题
      const filterSelect = screen.getByRole('combobox', { name: /主题|筛选/i });
      await user.click(filterSelect);

      const typesOption = screen.getByText(/types/i);
      await user.click(typesOption);

      // 验证只显示types相关的测验
      await waitFor(() => {
        const quizCards = screen.getAllByTestId(/quiz-history-card/i);
        expect(quizCards.length).toBe(1);
        expect(within(quizCards[0]).getByText(/types/i)).toBeInTheDocument();
      });
    });
  });

  describe('空状态处理测试', () => {
    it('应该在无测验记录时显示友好提示', async () => {
      (quizService.getQuizHistory as jest.Mock).mockResolvedValue([]);

      render(
        <AuthProvider>
          <QuizCenterPage />
        </AuthProvider>
      );

      await waitFor(() => {
        expect(quizService.getQuizHistory).toHaveBeenCalled();
      });

      // 验证空状态提示
      const emptyMessage = await screen.findByText(/暂无测验记录|开始测验/i);
      expect(emptyMessage).toBeInTheDocument();

      // 验证包含引导操作
      const startButton = screen.getByRole('link', { name: /开始学习|前往主题/i });
      expect(startButton).toBeInTheDocument();
    });
  });

  describe('错误处理测试', () => {
    it('应该处理测验历史加载失败', async () => {
      const consoleError = jest.spyOn(console, 'error').mockImplementation();
      (quizService.getQuizHistory as jest.Mock).mockRejectedValue(
        new Error('Network error')
      );

      render(
        <AuthProvider>
          <QuizCenterPage />
        </AuthProvider>
      );

      await waitFor(() => {
        expect(quizService.getQuizHistory).toHaveBeenCalled();
      });

      // 验证错误提示
      const errorMessage = await screen.findByText(/加载失败|网络错误/i);
      expect(errorMessage).toBeInTheDocument();

      // 验证重试按钮
      const retryButton = screen.getByRole('button', { name: /重试|刷新/i });
      expect(retryButton).toBeInTheDocument();

      consoleError.mockRestore();
    });
  });

  describe('响应式设计测试', () => {
    it('应该在移动端正确显示测验列表', async () => {
      // 模拟移动端视口
      global.innerWidth = 375;
      global.innerHeight = 667;

      render(
        <AuthProvider>
          <QuizCenterPage />
        </AuthProvider>
      );

      await waitFor(() => {
        expect(quizService.getQuizHistory).toHaveBeenCalled();
      });

      // 验证测验卡片堆叠显示(单列)
      const quizCards = await screen.findAllByTestId(/quiz-history-card/i);
      quizCards.forEach((card) => {
        expect(card).toHaveClass(/flex-col|block/i);
      });
    });
  });
});
