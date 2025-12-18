import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import { useParams, useRouter } from 'next/navigation';
import QuizReviewPage from '@/app/(protected)/quiz/history/[sessionId]/page';
import { fetchQuizReview } from '@/lib/quiz';

jest.mock('next/navigation', () => ({
    useParams: jest.fn(),
    useRouter: jest.fn(),
}));

jest.mock('@/lib/quiz', () => ({
    fetchQuizReview: jest.fn(),
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

const mockUseParams = useParams as jest.MockedFunction<typeof useParams>;
const mockUseRouter = useRouter as jest.MockedFunction<typeof useRouter>;
const mockFetchQuizReview = fetchQuizReview as jest.MockedFunction<typeof fetchQuizReview>;

describe('QuizReviewPage', () => {
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
        mockUseParams.mockReturnValue({ sessionId: 'session-1' } as ReturnType<typeof useParams>);
        mockFetchQuizReview.mockImplementation(() => new Promise(() => {}));

        render(<QuizReviewPage />);
        expect(screen.getByTestId('loading')).toBeInTheDocument();
    });

    it('应当显示错误信息', async () => {
        mockUseParams.mockReturnValue({ sessionId: 'session-1' } as ReturnType<typeof useParams>);
        mockFetchQuizReview.mockRejectedValue(new Error('加载失败'));

        render(<QuizReviewPage />);

        await waitFor(() => {
            expect(screen.getByTestId('error')).toBeInTheDocument();
            expect(screen.getByText('加载回顾详情失败')).toBeInTheDocument();
        });
    });

    it('应当显示回顾详情', async () => {
        mockUseParams.mockReturnValue({ sessionId: 'session-1' } as ReturnType<typeof useParams>);
        const mockReview = {
            meta: {
                sessionId: 'session-1',
                topic: 'constants',
                chapter: 'boolean',
                score: 80,
                passed: true,
                completedAt: '2025-12-18T10:00:00Z',
            },
            items: [
                {
                    questionId: 1,
                    stem: '布尔常量的零值是？',
                    options: ['false', 'true'],
                    userChoice: 'false',
                    correctChoice: 'false',
                    isCorrect: true,
                    explanation: '布尔类型的零值是 false',
                },
            ],
        };

        mockFetchQuizReview.mockResolvedValue(mockReview);

        render(<QuizReviewPage />);

        await waitFor(() => {
            expect(screen.getByText('测验回顾')).toBeInTheDocument();
            expect(screen.getByText('布尔常量的零值是？')).toBeInTheDocument();
        });
    });
});

