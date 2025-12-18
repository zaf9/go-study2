import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import { useRouter } from 'next/navigation';
import QuizHistoryPage from '@/app/(protected)/quiz/history/page';
import { useQuizHistory } from '@/hooks/useQuiz';

jest.mock('next/navigation', () => ({
    useRouter: jest.fn(),
}));

jest.mock('@/hooks/useQuiz', () => ({
    useQuizHistory: jest.fn(),
}));

jest.mock('@/components/common/Loading', () => {
    return function Loading() {
        return <div data-testid="loading">Loading...</div>;
    };
});

jest.mock('@/components/common/ErrorMessage', () => {
    return function ErrorMessage({ message, description }: { message: string; description?: string }) {
        return (
            <div data-testid="error">
                <div>{message}</div>
                {description && <div>{description}</div>}
            </div>
        );
    };
});

const mockUseRouter = useRouter as jest.MockedFunction<typeof useRouter>;
const mockUseQuizHistory = useQuizHistory as jest.MockedFunction<typeof useQuizHistory>;

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

    it('应当显示加载状态', () => {
        mockUseQuizHistory.mockReturnValue({
            history: [],
            isLoading: true,
            error: null,
        } as ReturnType<typeof useQuizHistory>);

        render(<QuizHistoryPage />);
        expect(screen.getByTestId('loading')).toBeInTheDocument();
    });

    it('应当显示错误信息', () => {
        const error = new Error('加载失败');
        mockUseQuizHistory.mockReturnValue({
            history: [],
            isLoading: false,
            error,
        } as ReturnType<typeof useQuizHistory>);

        render(<QuizHistoryPage />);
        expect(screen.getByTestId('error')).toBeInTheDocument();
        expect(screen.getByText('加载测验历史失败')).toBeInTheDocument();
    });

    it('应当显示历史记录列表', async () => {
        const mockHistory = [
            {
                id: 1,
                sessionId: 'session-1',
                topic: 'constants',
                chapter: 'boolean',
                score: 80,
                totalQuestions: 10,
                correctAnswers: 8,
                passed: true,
                completedAt: '2025-12-18T10:00:00Z',
            },
        ];

        mockUseQuizHistory.mockReturnValue({
            history: mockHistory,
            isLoading: false,
            error: null,
        } as ReturnType<typeof useQuizHistory>);

        render(<QuizHistoryPage />);

        await waitFor(() => {
            expect(screen.getByText('测验历史记录')).toBeInTheDocument();
        });
    });
});

