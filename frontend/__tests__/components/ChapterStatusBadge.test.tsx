/**
 * ChapterStatusBadge 组件测试
 * 测试章节状态徽章的渲染和显示
 */

import React from 'react';
import { render, screen } from '@testing-library/react';
import ChapterStatusBadge from '@/components/learning/ChapterStatusBadge';

describe('ChapterStatusBadge 组件', () => {
  it('应显示 "未开始" 状态当没有进度数据', () => {
    render(<ChapterStatusBadge progressData={null} />);

    expect(screen.getByText('🔘')).toBeInTheDocument();
    expect(screen.getByText('未开始')).toBeInTheDocument();
  });

  it('应显示 "已完成" 状态当 completed_at 有值', () => {
    const progressData = {
      status: 'in_progress',
      completed_at: '2024-01-01T00:00:00Z',
    };

    render(<ChapterStatusBadge progressData={progressData} />);

    expect(screen.getByText('✅')).toBeInTheDocument();
    expect(screen.getByText('已完成')).toBeInTheDocument();
  });

  it('应显示 "学习中" 状态当 first_visit_at 有值', () => {
    const progressData = {
      status: 'not_started',
      first_visit_at: '2024-01-01T00:00:00Z',
    };

    render(<ChapterStatusBadge progressData={progressData} />);

    expect(screen.getByText('📖')).toBeInTheDocument();
    expect(screen.getByText('学习中')).toBeInTheDocument();
  });

  it('应仅显示图标当 showText 为 false', () => {
    const progressData = {
      status: 'completed',
      completed_at: '2024-01-01T00:00:00Z',
    };

    render(<ChapterStatusBadge progressData={progressData} showText={false} />);

    expect(screen.getByText('✅')).toBeInTheDocument();
    expect(screen.queryByText('已完成')).not.toBeInTheDocument();
  });

  it('应仅显示文本当 showIcon 为 false', () => {
    const progressData = {
      status: 'completed',
      completed_at: '2024-01-01T00:00:00Z',
    };

    render(<ChapterStatusBadge progressData={progressData} showIcon={false} />);

    expect(screen.queryByText('✅')).not.toBeInTheDocument();
    expect(screen.getByText('已完成')).toBeInTheDocument();
  });

  it('应处理 undefined 进度数据', () => {
    render(<ChapterStatusBadge progressData={undefined} />);

    expect(screen.getByText('🔘')).toBeInTheDocument();
    expect(screen.getByText('未开始')).toBeInTheDocument();
  });

  it('应使用正确的颜色渲染 Tag 组件', () => {
    const { container } = render(<ChapterStatusBadge progressData={null} />);

    const tag = container.querySelector('.ant-tag');
    expect(tag).toHaveClass('ant-tag-default');
  });
});
