import React from 'react';
import { render, screen } from '@testing-library/react';
import QuizMetaInfo from '../QuizMetaInfo';
import { QuizStats } from '@/types/quiz';

const mockStats: QuizStats = {
    total: 35,
    byType: {
        single: 18,
        multiple: 17,
    },
    byDifficulty: {
        easy: 14,
        medium: 14,
        hard: 7,
    },
};

describe('QuizMetaInfo', () => {
    it('应当显示总题量', () => {
        render(<QuizMetaInfo stats={mockStats} />);
        expect(screen.getByText('总题量')).toBeInTheDocument();
        expect(screen.getByText('35')).toBeInTheDocument();
    });

    it('应当显示预计用时', () => {
        render(<QuizMetaInfo stats={mockStats} />);
        expect(screen.getByText('预计用时')).toBeInTheDocument();
        expect(screen.getByText(/分钟/)).toBeInTheDocument();
    });

    it('应当显示难度分布', () => {
        render(<QuizMetaInfo stats={mockStats} />);
        expect(screen.getByText('难度分布:')).toBeInTheDocument();
        expect(screen.getByText(/简单.*14.*题/)).toBeInTheDocument();
        expect(screen.getByText(/中等.*14.*题/)).toBeInTheDocument();
        expect(screen.getByText(/困难.*7.*题/)).toBeInTheDocument();
    });

    it('应当显示题型分布', () => {
        render(<QuizMetaInfo stats={mockStats} />);
        expect(screen.getByText('题型分布:')).toBeInTheDocument();
        expect(screen.getByText(/单选题.*18.*题/)).toBeInTheDocument();
        expect(screen.getByText(/多选题.*17.*题/)).toBeInTheDocument();
    });

    it('加载中时应当显示加载状态', () => {
        render(<QuizMetaInfo stats={null} loading={true} />);
        expect(screen.getByText('加载中...')).toBeInTheDocument();
    });

    it('stats 为 null 时应当显示加载状态', () => {
        render(<QuizMetaInfo stats={null} />);
        expect(screen.getByText('加载中...')).toBeInTheDocument();
    });

    it('应当正确处理空的统计数据', () => {
        const emptyStats: QuizStats = {
            total: 0,
            byType: {},
            byDifficulty: {},
        };
        render(<QuizMetaInfo stats={emptyStats} />);
        expect(screen.getByText('总题量')).toBeInTheDocument();
        // 验证总题量显示为 0（使用 getAllByText 因为预计用时也可能是 0）
        const zeroElements = screen.getAllByText('0');
        expect(zeroElements.length).toBeGreaterThan(0);
    });
});

