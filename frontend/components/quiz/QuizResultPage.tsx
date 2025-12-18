"use client";

import React from 'react';
import { Card, Result, Button, Statistic, Row, Col, Typography } from 'antd';
import { CheckCircleOutlined, CloseCircleOutlined, TrophyOutlined, CloseOutlined } from '@ant-design/icons';
import { QuizSubmitResult } from '@/types/quiz';

const { Title, Text, Paragraph } = Typography;

interface QuizResultPageProps {
    result: QuizSubmitResult;
    onRetry: () => void;
    onReview: () => void;
}

/**
 * QuizResultPage
 * 测验结果展示页，包含百分制得分、及格状态标识、得分统计及操作按钮。
 */
const QuizResultPage: React.FC<QuizResultPageProps> = ({ result, onRetry, onReview }) => {
    const isPassed = result.passed;

    return (
        <Card bordered={false} className="shadow-sm scale-in-center overflow-hidden">
            <div style={{ position: 'relative' }}>
                {/* 背景装饰图，可根据是否通过切换颜色 */}
                <div
                    style={{
                        position: 'absolute',
                        top: -20,
                        right: -20,
                        opacity: 0.1,
                        fontSize: '120px',
                        color: isPassed ? '#52c41a' : '#ff4d4f'
                    }}
                >
                    {isPassed ? <CheckCircleOutlined /> : <CloseCircleOutlined />}
                </div>

                <Result
                    status={isPassed ? "success" : "warning"}
                    icon={isPassed ? <TrophyOutlined style={{ color: '#faad14', fontSize: 64 }} /> : null}
                    title={
                        <div style={{ marginTop: 16 }}>
                            <Title level={2}>{isPassed ? "恭喜通过！" : "请继续加油！"}</Title>
                            <Text type="secondary">您已经完成了本章节的测验</Text>
                        </div>
                    }
                    subTitle={
                        <div className="bg-gray-50 p-6 rounded-lg mt-6">
                            <Row gutter={24} justify="center">
                                <Col span={8}>
                                    <Statistic
                                        title="本次得分"
                                        value={result.score}
                                        suffix="/ 100"
                                        valueStyle={{ color: isPassed ? '#3f8600' : '#cf1322', fontWeight: 700, fontSize: 32 }}
                                    />
                                </Col>
                                <Col span={8}>
                                    <Statistic
                                        title="正确数"
                                        value={result.correct_answers}
                                        prefix={<CheckCircleOutlined style={{ color: '#52c41a' }} />}
                                        suffix={`/ ${result.total_questions}`}
                                    />
                                </Col>
                                <Col span={8}>
                                    <Statistic
                                        title="错误数"
                                        value={result.total_questions - result.correct_answers}
                                        prefix={<CloseOutlined style={{ color: '#ff4d4f' }} />}
                                    />
                                </Col>
                            </Row>
                        </div>
                    }
                    extra={[
                        <Button type="primary" key="review" size="large" onClick={onReview} icon={<CheckCircleOutlined />}>
                            查看解析
                        </Button>,
                        <Button key="retry" size="large" onClick={onRetry}>
                            再测一次
                        </Button>,
                    ]}
                >
                    <div className="text-left mt-8">
                        <Paragraph>
                            <Text strong>学习建议：</Text>
                            {isPassed
                                ? "您已经掌握了本章节的核心知识点，可以继续挑战下一个章节或进行更深入的学习。"
                                : "本章节部分知识点可能还不够牢固，建议您回顾相关学习内容，并针对错题进行专项练习。"}
                        </Paragraph>
                    </div>
                </Result>
            </div>
        </Card>
    );
};

export default QuizResultPage;
