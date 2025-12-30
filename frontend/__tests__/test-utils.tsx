/**
 * 测试工具 - React Testing Library 自定义配置
 * 提供包装了 Provider 的自定义 render 函数
 */

import React, { ReactElement } from 'react';
import { render, RenderOptions } from '@testing-library/react';
import { SWRConfig } from 'swr';

// Mock AuthContext if needed
const MockAuthProvider = ({ children }: { children: React.ReactNode }) => {
  return <>{children}</>;
};

// Custom render function with providers
interface CustomRenderOptions extends Omit<RenderOptions, 'wrapper'> {
  // Add custom options if needed
}

function customRender(
  ui: ReactElement,
  options?: CustomRenderOptions
) {
  const AllTheProviders = ({ children }: { children: React.ReactNode }) => {
    return (
      <SWRConfig value={{ provider: () => new Map(), dedupingInterval: 0 }}>
        <MockAuthProvider>{children}</MockAuthProvider>
      </SWRConfig>
    );
  };

  return render(ui, { wrapper: AllTheProviders, ...options });
}

// Re-export everything
export * from '@testing-library/react';
export { customRender as render };

// Mock data helpers
export const mockChapterProgress = {
  topic: 'lexical_elements' as const,
  chapter: 'comments',
  status: 'in_progress' as const,
  scrollProgress: 50,
  readDuration: 300,
  lastVisitAt: new Date().toISOString(),
};

export const mockTopicProgress = {
  id: 'lexical_elements' as const,
  name: '词法元素',
  weight: 1,
  progress: 50,
  completedChapters: 5,
  totalChapters: 11,
  lastVisitAt: new Date().toISOString(),
};

export const mockQuizSession = {
  sessionId: 'test-session-123',
  topic: 'lexical_elements',
  chapter: 'comments',
  questions: [
    {
      id: 1,
      type: 'single',
      difficulty: 'easy',
      question: 'What is a line comment in Go?',
      options: [
        { id: 'A', label: '//' },
        { id: 'B', label: '/* */' },
        { id: 'C', label: '#' },
      ],
      codeSnippet: null,
    },
  ],
};

export const mockQuizResult = {
  score: 80,
  total: 10,
  correctIds: ['1', '2', '3', '4', '5', '6', '7', '8'],
  wrongIds: ['9', '10'],
  submittedAt: new Date().toISOString(),
  durationMs: 120000,
};
