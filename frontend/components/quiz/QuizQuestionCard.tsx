"use client";

import React from "react";
import { Checkbox, Radio, Space, Typography, Tag } from "antd";
import { QuizQuestion } from "@/types/quiz";
import QuestionTypeTag from "./QuestionTypeTag";

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
        if (isMultiple) return "多选题：请选择一个或多个答案。选中后内容背景将高亮。";
        if (question.type === "truefalse") return "判断题：选择 T (正确) 或 F (错误)。";
        if (question.type === "code_output") return "代码输出题：阅读代码并选择正确的输出结果。";
        return "单选题：请选择一个最合适的答案。";
    };

    /**
     * 将题目的 options 映射为带有 A, B, C... 标签的布局
     * labels 不受 shuffle 影响，因为 shuffle 是在后端完成且返回给前端时 options 已经是确定顺序的数组
     * 
     * 标签生成逻辑：
     * - 1-26: A-Z（标准字母）
     * - 27-52: AA-AZ（二字母组合）
     * - 53-78: BA-BZ（二字母组合）
     * 等等，支持最多千余个选项
     */
    const getOptionLabel = (index: number): string => {
        if (index < 26) {
            // A-Z
            return String.fromCharCode(65 + index);
        } else {
            // AA, AB, ... AZ, BA, BB, ... ZZ, AAA, ...
            let label = '';
            let num = index - 26;
            while (num >= 0) {
                label = String.fromCharCode(65 + (num % 26)) + label;
                num = Math.floor(num / 26) - 1;
                if (num < 0) break;
            }
            // 确保至少两个字母
            if (label.length === 1) {
                label = 'A' + label;
            }
            return label;
        }
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
                            const label = getOptionLabel(index);
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
                            const label = getOptionLabel(index);
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
                    <QuestionTypeTag type={question.type} />
                    <Tag color={difficultyColor(question.difficulty)}>{question.difficulty.toUpperCase()}</Tag>
                </Space>
                <Title level={4} style={{ marginTop: 16 }}>
                    {question.question}
                </Title>
                <Paragraph type="secondary" italic style={{ borderLeft: '3px solid #1890ff', paddingLeft: 12, margin: '12px 0' }}>
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
