import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import QuizResultPage from '../QuizResultPage';
import { QuizSubmitResult } from '@/types/quiz';

const mockResult: QuizSubmitResult = {
    score: 85,
    total_questions: 10,
    correct_answers: 8,
    passed: true,
    details: [],
};

describe('QuizResultPage', () => {
    const onRetry = jest.fn();
    const onReview = jest.fn();

    it('应当正确显示百分制得分和及格状态', () => {
        render(
            <QuizResultPage
                result={mockResult}
                onRetry={onRetry}
                onReview={onReview}
            />
        );

        // 验证得分包含 85
        expect(screen.getByText('85')).toBeInTheDocument();
        expect(screen.getByText('本次得分')).toBeInTheDocument();

        // 验证状态文本
        expect(screen.getByText(/恭喜通过/)).toBeInTheDocument();
    });

    it('不及格时应当显示鼓励文案', () => {
        const failedResult = { ...mockResult, score: 50, passed: false };
        render(
            <QuizResultPage
                result={failedResult}
                onRetry={onRetry}
                onReview={onReview}
            />
        );

        expect(screen.getByText('50')).toBeInTheDocument();
        expect(screen.getByText('请继续加油！')).toBeInTheDocument();
    });

    it('应当显示正确/错误题目数量统计', () => {
        render(<QuizResultPage result={mockResult} onRetry={onRetry} />);

        expect(screen.getByText('8')).toBeInTheDocument(); // 正确数
        expect(screen.getByText('2')).toBeInTheDocument(); // 错误数
    });

    it('点击按钮应当触发相应回调', () => {
        render(
            <QuizResultPage
                result={mockResult}
                onRetry={onRetry}
                onReview={onReview}
            />
        );

        const retryBtn = screen.getByRole('button', { name: /再测一次/i });
        fireEvent.click(retryBtn);
        expect(onRetry).toHaveBeenCalled();

        const reviewBtn = screen.getByRole('button', { name: /查看解析/i });
        fireEvent.click(reviewBtn);
        expect(onReview).toHaveBeenCalled();
    });
});
