import React from 'react';
import { render } from '@testing-library/react';
import { useRouter } from 'next/navigation';
import QuizHistoryPage from '@/app/(protected)/quiz/history/page';

jest.mock('next/navigation', () => ({
    useRouter: jest.fn(),
}));

const mockUseRouter = useRouter as jest.MockedFunction<typeof useRouter>;

describe('QuizHistoryPage', () => {
    beforeEach(() => {
        mockUseRouter.mockReturnValue({
            push: jest.fn(),
            replace: jest.fn(),
            prefetch: jest.fn(),
            back: jest.fn(),
            forward: jest.fn(),
            refresh: jest.fn(),
        } as ReturnType<typeof useRouter>);
    });

    afterEach(() => {
        jest.clearAllMocks();
    });

    it('应当重定向到测验中心', () => {
        const { container } = render(<QuizHistoryPage />);

        // 验证页面渲染为 null（因为会立即重定向）
        expect(container.firstChild).toBe(null);

        // 验证 replace 被调用，重定向到 /quiz-center
        expect(mockUseRouter.mock.results[0].value.replace).toHaveBeenCalledWith('/quiz-center');
    });

    it('应当在挂载后立即调用 router.replace', () => {
        render(<QuizHistoryPage />);

        // replace 应该被调用一次
        expect(mockUseRouter.mock.results[0].value.replace).toHaveBeenCalledTimes(1);
        expect(mockUseRouter.mock.results[0].value.replace).toHaveBeenCalledWith('/quiz-center');
    });
});

