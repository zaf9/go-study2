/**
 * 测试主题列表页面
 * 验证所有主题及章节数正确显示
 */

import { render, screen, waitFor } from '@/tests/test-utils';
import TopicsPage from '@/app/(protected)/topics/page';
import { fetchTopics } from '@/lib/learning';

// Mock the learning module
jest.mock('@/lib/learning', () => ({
  fetchTopics: jest.fn(),
}));

// Mock useProgress hook
jest.mock('@/hooks/useProgress', () => ({
  __esModule: true,
  default: () => ({
    next: null,
    isLoading: false,
  }),
}));

// Mock Next.js router
jest.mock('next/navigation', () => ({
  useRouter: () => ({
    push: jest.fn(),
  }),
}));

describe('TopicsPage', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('应该显示所有主题及准确的章节数量', async () => {
    // Mock fetch topics response
    (fetchTopics as jest.Mock).mockResolvedValue([
      {
        key: 'lexical_elements',
        title: '词法元素',
        summary: 'Go语言的基本词法元素',
        chapterCount: 11,
      },
      {
        key: 'constants',
        title: '常量',
        summary: 'Go语言常量详解',
        chapterCount: 12,
      },
      {
        key: 'variables',
        title: '变量',
        summary: 'Go语言变量系统',
        chapterCount: 4,
      },
      {
        key: 'types',
        title: '类型',
        summary: 'Go语言类型系统',
        chapterCount: 14,
      },
    ]);

    render(<TopicsPage />);

    // Wait for loading to complete
    await waitFor(() => {
      expect(screen.getByText('词法元素')).toBeInTheDocument();
    });

    // Verify chapter counts are displayed correctly
    expect(screen.getByText(/章节数.*11/)).toBeInTheDocument();
    expect(screen.getByText(/章节数.*12/)).toBeInTheDocument();
    expect(screen.getByText(/章节数.*4/)).toBeInTheDocument();
    expect(screen.getByText(/章节数.*14/)).toBeInTheDocument();
  });

  it('应该显示正确的主题标题', async () => {
    (fetchTopics as jest.Mock).mockResolvedValue([
      {
        key: 'lexical_elements',
        title: '词法元素',
        summary: 'Go语言的基本词法元素',
        chapterCount: 11,
      },
    ]);

    render(<TopicsPage />);

    await waitFor(() => {
      expect(screen.getByText('词法元素')).toBeInTheDocument();
    });
  });

  it('加载失败时应该显示错误信息', async () => {
    (fetchTopics as jest.Mock).mockRejectedValue(new Error('网络错误'));

    render(<TopicsPage />);

    await waitFor(() => {
      expect(screen.getByText(/加载主题列表失败/)).toBeInTheDocument();
    });
  });
});
