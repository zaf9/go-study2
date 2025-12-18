"use client";

import React from 'react';
import { Skeleton, Card, Space } from 'antd';

/**
 * QuizSkeletonLoader
 * 测验加载过程中的骨架屏组件，模拟测验卡片的布局展示。
 */
const QuizSkeletonLoader: React.FC = () => {
    return (
        <div style={{ maxWidth: 800, margin: '20px auto', padding: '0 20px' }}>
            <Card bordered={false} className="shadow-sm">
                <Space direction="vertical" size="large" style={{ width: '100%' }}>
                    {/* 题干部分骨架 */}
                    <div>
                        <Skeleton.Button active size="small" style={{ width: 80, marginBottom: 16 }} />
                        <Skeleton active title={{ width: '80%' }} paragraph={{ rows: 2 }} />
                    </div>

                    {/* 选项部分骨架 (模拟 A-D 四个选项) */}
                    <div style={{ marginTop: '20px' }}>
                        {[1, 2, 3, 4].map((i) => (
                            <div key={i} style={{ marginBottom: '16px', display: 'flex', alignItems: 'center' }}>
                                <Skeleton.Avatar active size="small" shape="square" style={{ marginRight: '16px', flexShrink: 0 }} />
                                <Skeleton.Input active size="default" style={{ width: '100%' }} />
                            </div>
                        ))}
                    </div>

                    {/* 底部按钮骨架 */}
                    <div style={{ display: 'flex', justifyContent: 'space-between', marginTop: '40px', borderTop: '1px solid #f0f0f0', paddingTop: '20px' }}>
                        <Skeleton.Button active size="default" style={{ width: 100 }} />
                        <Skeleton.Button active size="default" style={{ width: 100 }} />
                    </div>
                </Space>
            </Card>
        </div>
    );
};

export default QuizSkeletonLoader;
