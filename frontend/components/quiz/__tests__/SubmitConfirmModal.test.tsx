import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import SubmitConfirmModal from '../SubmitConfirmModal';

describe('SubmitConfirmModal', () => {
    const defaultProps = {
        open: true,
        answeredCount: 8,
        totalCount: 10,
        onConfirm: jest.fn(),
        onCancel: jest.fn(),
        submitting: false,
    };

    it('应当正确显示答题状态统计', () => {
        render(<SubmitConfirmModal {...defaultProps} />);

        expect(screen.getByText(/确认提交测验？/)).toBeInTheDocument();
        expect(screen.getByText(/已答题目：8/)).toBeInTheDocument();
        expect(screen.getByText(/总题目数：10/)).toBeInTheDocument();
    });

    it('未答完时应当显示警告信息', () => {
        render(<SubmitConfirmModal {...defaultProps} />);
        expect(screen.getByText(/您还有 2 道题目未完成/)).toBeInTheDocument();
    });

    it('答完时应当显示全部完成的正面反馈', () => {
        render(
            <SubmitConfirmModal
                {...defaultProps}
                answeredCount={10}
                totalCount={10}
            />
        );
        expect(screen.getByText(/您已经完成了所有题目/)).toBeInTheDocument();
    });

    it('点击确认应当触发 onConfirm', () => {
        const onConfirm = jest.fn();
        render(<SubmitConfirmModal {...defaultProps} onConfirm={onConfirm} />);

        const confirmButton = screen.getByRole('button', { name: /确 认/i });
        fireEvent.click(confirmButton);

        expect(onConfirm).toHaveBeenCalled();
    });

    it('点击取消应当触发 onCancel', () => {
        const onCancel = jest.fn();
        render(<SubmitConfirmModal {...defaultProps} onCancel={onCancel} />);

        const cancelButton = screen.getByRole('button', { name: /取 消/i });
        fireEvent.click(cancelButton);

        expect(onCancel).toHaveBeenCalled();
    });
});
