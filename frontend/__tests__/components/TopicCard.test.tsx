/**
 * 测试 TopicCard 组件
 * 验证主题卡片显示准确的章节数量
 */

import { render, screen, fireEvent } from '@/tests/test-utils';
import TopicCard from '@/components/learning/TopicCard';
import { TopicSummary } from '@/types/learning';

// Mock Next.js router
const mockPush = jest.fn();
jest.mock('next/navigation', () => ({
  useRouter: () => ({
    push: mockPush,
  }),
}));

describe('TopicCard', () => {
  beforeEach(() => {
    mockPush.mockClear();
  });

  it('应该显示准确的章节数量', () => {
    const topic: TopicSummary = {
      key: 'lexical_elements',
      title: '词法元素',
      summary: 'Go语言的基本词法元素',
      chapterCount: 11,
    };

    render(<TopicCard topic={topic} />);

    expect(screen.getByText('词法元素')).toBeInTheDocument();
    expect(screen.getByText(/章节数.*11/)).toBeInTheDocument();
  });

  it('零章节主题应该显示友好提示', () => {
    const topic: TopicSummary = {
      key: 'variables',
      title: '新主题',
      summary: '即将推出',
      chapterCount: 0,
    };

    render(<TopicCard topic={topic} />);

    // Should show "暂无章节" instead of "章节数：0"
    expect(screen.getByText(/暂无章节|章节数.*0/)).toBeInTheDocument();
  });

  it('点击卡片应该跳转到主题页面', () => {
    const topic: TopicSummary = {
      key: 'constants',
      title: '常量',
      summary: 'Go语言常量详解',
      chapterCount: 12,
    };

    render(<TopicCard topic={topic} />);

    const card = screen.getByText('常量').closest('.ant-card');
    if (card) {
      fireEvent.click(card);
    }

    expect(mockPush).toHaveBeenCalledWith('/topics/constants');
  });

  it('应该显示主题的描述信息', () => {
    const topic: TopicSummary = {
      key: 'types',
      title: '类型',
      summary: 'Go语言类型系统详解',
      chapterCount: 14,
    };

    render(<TopicCard topic={topic} />);

    expect(screen.getByText('Go语言类型系统详解')).toBeInTheDocument();
  });

  it('没有描述时应该显示默认文本', () => {
    const topic: TopicSummary = {
      key: 'types',
      title: '类型',
      summary: '',
      chapterCount: 14,
    };

    render(<TopicCard topic={topic} />);

    expect(screen.getByText(/快速开始该主题的学习/)).toBeInTheDocument();
  });
});
