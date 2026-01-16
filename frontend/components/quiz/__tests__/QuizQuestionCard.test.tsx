import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import QuizQuestionCard from '../QuizQuestionCard';
import { QuizQuestion } from '@/types/quiz';

const mockQuestion: QuizQuestion = {
    id: 1,
    type: 'single',
    difficulty: 'easy',
    question: 'Go 语言的作者之一是谁？',
    options: [
        { id: '1', label: 'Rob Pike' },
        { id: '2', label: 'Ken Thompson' },
        { id: '3', label: 'Robert Griesemer' },
        { id: '4', label: '以上都是' },
    ],
};

describe('QuizQuestionCard', () => {
    it('应当正确渲染题干和选项标签 A-D', () => {
        render(
            <QuizQuestionCard
                question={mockQuestion}
                value={[]}
                onChange={jest.fn()}
            />
        );

        expect(screen.getByText(/Go 语言的作者之一是谁？/)).toBeInTheDocument();

        // 验证 A-D 标签是否存在
        expect(screen.getByText(/^A\./)).toBeInTheDocument();
        expect(screen.getByText(/^B\./)).toBeInTheDocument();
        expect(screen.getByText(/^C\./)).toBeInTheDocument();
        expect(screen.getByText(/^D\./)).toBeInTheDocument();
    });

    it('单选题应当渲染并响应点击', () => {
        const onChange = jest.fn();
        render(
            <QuizQuestionCard
                question={mockQuestion}
                value={[]}
                onChange={onChange}
            />
        );

        const optionA = screen.getByLabelText(/Rob Pike/);
        fireEvent.click(optionA);

        // 注意：现在提交选项ID而不是label
        expect(onChange).toHaveBeenCalledWith(['1']);
    });

    it('多选题应当显示多选引导语并允许选择多个', () => {
        const multipleQuestion: QuizQuestion = {
            ...mockQuestion,
            type: 'multiple',
        };
        const onChange = jest.fn();
        render(
            <QuizQuestionCard
                question={multipleQuestion}
                value={['1']}
                onChange={onChange}
            />
        );

        expect(screen.getByText(/请选择一个或多个/)).toBeInTheDocument();

        const optionB = screen.getByLabelText(/Ken Thompson/);
        fireEvent.click(optionB);

        // 应该包含原来的和新选的（使用ID而不是label）
        expect(onChange).toHaveBeenCalledWith(['1', '2']);
    });

    it('应当正确处理空选项或异常情况', () => {
        const emptyOptionQuestion: QuizQuestion = {
            ...mockQuestion,
            options: [],
        };
        render(
            <QuizQuestionCard
                question={emptyOptionQuestion}
                value={[]}
                onChange={jest.fn()}
            />
        );

        expect(screen.queryByText(/A\./)).not.toBeInTheDocument();
    });
});
