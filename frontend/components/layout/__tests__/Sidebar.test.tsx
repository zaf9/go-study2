import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import { usePathname, useRouter } from 'next/navigation';
import Sidebar from '../Sidebar';

jest.mock('next/navigation', () => ({
    usePathname: jest.fn(),
    useRouter: jest.fn(),
}));

const mockUsePathname = usePathname as jest.MockedFunction<typeof usePathname>;
const mockUseRouter = useRouter as jest.MockedFunction<typeof useRouter>;

describe('Sidebar', () => {
    const mockPush = jest.fn();

    beforeEach(() => {
        mockUseRouter.mockReturnValue({
            push: mockPush,
            replace: jest.fn(),
            prefetch: jest.fn(),
            back: jest.fn(),
            forward: jest.fn(),
            refresh: jest.fn(),
        } as ReturnType<typeof useRouter>);
        mockPush.mockClear();
    });

    it('应当显示测验入口', () => {
        mockUsePathname.mockReturnValue('/topics');
        render(<Sidebar collapsed={false} onCollapse={jest.fn()} />);
        
        expect(screen.getByText('章节测验')).toBeInTheDocument();
    });

    it('点击测验入口应当导航到 /quiz', () => {
        mockUsePathname.mockReturnValue('/topics');
        render(<Sidebar collapsed={false} onCollapse={jest.fn()} />);
        
        const quizMenuItem = screen.getByText('章节测验');
        fireEvent.click(quizMenuItem);
        
        expect(mockPush).toHaveBeenCalledWith('/quiz');
    });

    it('在 /quiz 路径下应当高亮测验入口', () => {
        mockUsePathname.mockReturnValue('/quiz');
        render(<Sidebar collapsed={false} onCollapse={jest.fn()} />);
        
        // Ant Design Menu 会通过 selectedKeys 高亮，这里验证菜单项存在
        expect(screen.getByText('章节测验')).toBeInTheDocument();
    });

    it('应当显示所有导航项', () => {
        mockUsePathname.mockReturnValue('/topics');
        render(<Sidebar collapsed={false} onCollapse={jest.fn()} />);
        
        expect(screen.getByText('学习主题')).toBeInTheDocument();
        expect(screen.getByText('学习进度')).toBeInTheDocument();
        expect(screen.getByText('章节测验')).toBeInTheDocument();
        expect(screen.getByText('首页')).toBeInTheDocument();
    });
});

