"use client";

import React from 'react';
import { Tag } from 'antd';

interface QuestionTypeTagProps {
    type: string;
}

/**
 * QuestionTypeTag
 * 根据题型显示不同的标签和颜色。
 */
const QuestionTypeTag: React.FC<QuestionTypeTagProps> = ({ type }) => {
    const getTypeConfig = (t: string) => {
        switch (t) {
            case 'single':
                return { color: 'blue', label: '单选题' };
            case 'multiple':
                return { color: 'purple', label: '多选题' };
            case 'truefalse':
                return { color: 'cyan', label: '判断题' };
            case 'code_output':
                return { color: 'orange', label: '代码输出' };
            case 'code_correction':
                return { color: 'magenta', label: '程序改错' };
            default:
                return { color: 'default', label: t || '未知题型' };
        }
    };

    const config = getTypeConfig(type);

    return (
        <Tag color={config.color} style={{ fontWeight: 500 }}>
            {config.label}
        </Tag>
    );
};

export default QuestionTypeTag;
