import type React from "react";

interface PageContainerProps {
    children: React.ReactNode;
    className?: string;
    maxWidth?: "max-w-7xl" | "max-w-6xl" | "max-w-5xl" | "max-w-4xl" | "max-w-full" | "max-w-screen-2xl";
    "data-testid"?: string;
}

/**
 * 统一的页面容器组件
 * 提供响应式内边距和最大宽度限制
 */
export default function PageContainer({
    children,
    className = "",
    maxWidth = "max-w-screen-2xl",
    "data-testid": dataTestId,
}: PageContainerProps) {
    return (
        <div className={`${maxWidth} mx-auto px-4 sm:px-6 lg:px-8 ${className}`} data-testid={dataTestId}>
            {children}
        </div>
    );
}
