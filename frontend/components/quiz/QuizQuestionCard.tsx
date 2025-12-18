"use client";

import React from "react";
import { Checkbox, Radio, Space, Typography, Tag } from "antd";
import { QuizQuestion } from "@/types/quiz";

const { Paragraph, Text, Title } = Typography;

/**
 * QuizQuestionCard 组件属性
 */
interface QuizQuestionProps {
    question: QuizQuestion;
    value: string[]; // 存储的是选项的文本内容 (Content-Based)
    onChange: (choices: string[]) => void;
    disabled?: boolean; // 用于回顾模式
}

/**
 * QuizQuestionCard
 * 实现测验题目的展示，支持 A-D 有序标签渲染，且标签与内容绑定关系稳定（基于数组索引）。
 */
const QuizQuestionCard: React.FC<QuizQuestionProps> = ({
    question,
    value,
    onChange,
    disabled = false,
}) => {
    const isMultiple =
        question.type === "multiple" || question.type === "code_correction";
    const isSingle =
        question.type === "single" ||
        question.type === "truefalse" ||
        question.type === "code_output";

    // 根据题型获取引导语
    const getGuidance = () => {
        if (isMultiple) return "（多选题：请选择一个或多个答案）";
        if (question.type === "truefalse") return "（判断题）";
        return "（单选题）";
    };

    /**
     * 将题目的 options 映射为带有 A, B, C... 标签的布局
     * labels 不受 shuffle 影响，因为 shuffle 是在后端完成且返回给前端时 options 已经是确定顺序的数组
     */
    const renderOptions = () => {
        if (!question.options || question.options.length === 0) return null;

        if (isMultiple) {
            return (
                <Checkbox.Group
                    disabled={disabled}
                    style={{ width: "100%" }}
                    value={value}
                    onChange={(val) => onChange(val as string[])}
                >
                    <Space direction="vertical" style={{ width: "100%" }}>
                        {question.options.map((opt, index) => {
                            const label = String.fromCharCode(65 + index); // 0 -> A, 1 -> B...
                            return (
                                <div key={opt.id} style={{ display: 'flex', alignItems: 'flex-start', padding: '8px 0' }}>
                                    <Checkbox value={opt.label}>
                                        <Text strong style={{ marginRight: 8 }}>{label}.</Text>
                                        <Text>{opt.label}</Text>
                                    </Checkbox>
                                </div>
                            );
                        })}
                    </Space>
                </Checkbox.Group>
            );
        }

        if (isSingle) {
            return (
                <Radio.Group
                    disabled={disabled}
                    style={{ width: "100%" }}
                    value={value?.[0]}
                    onChange={(e) => onChange([e.target.value])}
                >
                    <Space direction="vertical" style={{ width: "100%" }}>
                        {question.options.map((opt, index) => {
                            const label = String.fromCharCode(65 + index);
                            return (
                                <Radio key={opt.id} value={opt.label} style={{ display: 'flex', alignItems: 'flex-start', padding: '8px 0', whiteSpace: 'normal' }}>
                                    <Text strong style={{ marginRight: 8 }}>{label}.</Text>
                                    <Text>{opt.label}</Text>
                                </Radio>
                            );
                        })}
                    </Space>
                </Radio.Group>
            );
        }

        return null;
    };

    return (
        <div className="bg-white p-6 rounded shadow-sm border border-gray-100 mb-4 scale-in-center">
            {/* 题干标题与辅助信息 */}
            <div className="mb-4">
                <Space wrap>
                    <Tag color="blue">{labelByType(question.type)}</Tag>
                    <Tag color={difficultyColor(question.difficulty)}>{question.difficulty.toUpperCase()}</Tag>
                </Space>
                <Title level={4} style={{ marginTop: 16 }}>
                    {question.question}
                </Title>
                <Paragraph type="secondary" italic>
                    {getGuidance()}
                </Paragraph>
            </div>

            {/* 代码片段支持 */}
            {question.codeSnippet && (
                <div className="mb-4">
                    <pre className="bg-gray-800 text-gray-100 p-4 rounded-md overflow-x-auto font-mono text-sm leading-relaxed">
                        <code>{question.codeSnippet}</code>
                    </pre>
                </div>
            )}

            {/* 选项展示区 */}
            <div className="quiz-options-container">
                {renderOptions()}
            </div>
        </div>
    );
};

// 辅助函数：根据题型返回中文标签
function labelByType(t: string) {
    switch (t) {
        case "multiple": return "多选题";
        case "truefalse": return "判断题";
        case "code_output": return "代码输出";
        case "code_correction": return "程序改错";
        default: return "单选题";
    }
}

// 辅助函数：根据难度返回颜色
function difficultyColor(d: string) {
    switch (d.toLowerCase()) {
        case "easy": return "success";
        case "medium": return "warning";
        case "hard": return "error";
        default: return "default";
    }
}

export default QuizQuestionCard;
