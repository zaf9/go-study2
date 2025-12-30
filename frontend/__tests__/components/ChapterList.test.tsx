/**
 * ChapterList 组件测试
 * 测试章节列表的渲染和状态显示
 */

import React from 'react';
import { render, screen } from '@testing-library/react';
import ChapterList from '@/components/learning/ChapterList';
import { ChapterSummary, ChapterProgress } from '@/types/learning';

// Mock next/navigation
jest.mock('next/navigation', () => ({
  useRouter: () => ({
    push: jest.fn(),
  }),
}));

const mockChapters: ChapterSummary[] = [
  {
    id: 'chapter1',
    title: '第一章',
    summary: '第一章内容',
    order: 0,
  },
  {
    id: 'chapter2',
    title: '第二章',
    summary: '第二章内容',
    order: 1,
  },
  {
    id: 'chapter3',
    title: '第三章',
    summary: '第三章内容',
    order: 2,
  },
];

const mockProgress: ChapterProgress[] = [
  {
    userId: 1,
    topic: 'lexical_elements',
    chapter: 'chapter1',
    status: 'completed',
    readDuration: 120,
    scrollProgress: 100,
    lastPosition: 'end',
    quizScore: 100,
    quizPassed: true,
    firstVisitAt: '2024-01-01T00:00:00Z',
    lastVisitAt: '2024-01-01T01:00:00Z',
    completedAt: '2024-01-01T01:00:00Z',
    percent: 100,
  },
  {
    userId: 1,
    topic: 'lexical_elements',
    chapter: 'chapter2',
    status: 'in_progress',
    readDuration: 60,
    scrollProgress: 50,
    lastPosition: 'middle',
    quizScore: 0,
    quizPassed: false,
    firstVisitAt: '2024-01-01T02:00:00Z',
    lastVisitAt: '2024-01-01T02:01:00Z',
    percent: 50,
  },
];

describe('ChapterList 组件', () => {
  it('应渲染所有章节', () => {
    render(
      <ChapterList
        topicKey="lexical_elements"
        chapters={mockChapters}
      />
    );

    expect(screen.getByText('第一章')).toBeInTheDocument();
    expect(screen.getByText('第二章')).toBeInTheDocument();
    expect(screen.getByText('第三章')).toBeInTheDocument();
  });

  it('应显示章节序号', () => {
    render(
      <ChapterList
        topicKey="lexical_elements"
        chapters={mockChapters}
      />
    );

    expect(screen.getByText('序号：1')).toBeInTheDocument();
    expect(screen.getByText('序号：2')).toBeInTheDocument();
    expect(screen.getByText('序号：3')).toBeInTheDocument();
  });

  it('应显示章节摘要', () => {
    render(
      <ChapterList
        topicKey="lexical_elements"
        chapters={mockChapters}
      />
    );

    expect(screen.getByText('第一章内容')).toBeInTheDocument();
    expect(screen.getByText('第二章内容')).toBeInTheDocument();
    expect(screen.getByText('第三章内容')).toBeInTheDocument();
  });

  it('应显示正确的章节状态徽章', () => {
    render(
      <ChapterList
        topicKey="lexical_elements"
        chapters={mockChapters}
        progress={mockProgress}
      />
    );

    // 第一章已完成
    expect(screen.getAllByText('✅')[0]).toBeInTheDocument();
    expect(screen.getAllByText('已完成')[0]).toBeInTheDocument();

    // 第二章学习中
    expect(screen.getAllByText('📖')[0]).toBeInTheDocument();
    expect(screen.getAllByText('学习中')[0]).toBeInTheDocument();
  });

  it('应显示"未开始"状态对于没有进度记录的章节', () => {
    render(
      <ChapterList
        topicKey="lexical_elements"
        chapters={mockChapters}
        progress={mockProgress}
      />
    );

    // 第三章没有进度记录，应显示未开始
    expect(screen.getByText('🔘')).toBeInTheDocument();
    expect(screen.getByText('未开始')).toBeInTheDocument();
  });

  it('应渲染查看按钮', () => {
    render(
      <ChapterList
        topicKey="lexical_elements"
        chapters={mockChapters}
      />
    );

    const buttons = screen.getAllByText('查看');
    expect(buttons).toHaveLength(3);
  });

  it('应处理空章节列表', () => {
    const { container } = render(
      <ChapterList
        topicKey="lexical_elements"
        chapters={[]}
      />
    );

    // Ant Design List 组件渲染空列表
    const list = container.querySelector('.ant-list');
    expect(list).toBeInTheDocument();
  });

  it('应处理空进度数据', () => {
    render(
      <ChapterList
        topicKey="lexical_elements"
        chapters={mockChapters}
        progress={[]}
      />
    );

    // 所有章节都应显示未开始状态
    const notStartedBadges = screen.getAllByText('🔘');
    expect(notStartedBadges).toHaveLength(3);
  });
});
