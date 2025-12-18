"use client";

import React from "react";
import { Modal, Typography, Progress, Space, Alert } from "antd";
import { QuestionCircleOutlined, InfoCircleOutlined } from "@ant-design/icons";

const { Text, Paragraph } = Typography;

interface SubmitConfirmModalProps {
    open: boolean;
    answeredCount: number;
    totalCount: number;
    onConfirm: () => void;
    onCancel: () => void;
    submitting?: boolean;
}

/**
 * SubmitConfirmModal
 * 测验提交前的二次确认弹窗，展示答题进度统计，防止误触导致未答完即提交。
 */
const SubmitConfirmModal: React.FC<SubmitConfirmModalProps> = ({
    open,
    answeredCount,
    totalCount,
    onConfirm,
    onCancel,
    submitting = false,
}) => {
    const isComplete = totalCount > 0 && answeredCount === totalCount;
    const remainingCount = totalCount - answeredCount;
    const progressPercent = totalCount > 0 ? Math.round((answeredCount / totalCount) * 100) : 0;

    return (
        <Modal
            title={
                <Space>
                    <QuestionCircleOutlined style={{ color: '#1890ff' }} />
                    <span>确认提交测验？</span>
                </Space>
            }
            open={open}
            onOk={onConfirm}
            onCancel={onCancel}
            confirmLoading={submitting}
            okText="确 认"
            cancelText="取 消"
            centered
            width={480}
        >
            <div style={{ textAlign: 'center', padding: '20px 0' }}>
                <Progress
                    type="circle"
                    percent={progressPercent}
                    status={isComplete ? "success" : "normal"}
                    format={() => `${answeredCount}/${totalCount}`}
                    width={100}
                />

                <div style={{ marginTop: 24 }}>
                    <Paragraph strong style={{ fontSize: 16 }}>
                        {isComplete
                            ? "太棒了！您已经完成了所有题目。"
                            : `您还有 ${remainingCount} 道题目未完成。`}
                    </Paragraph>

                    <Space direction="vertical" size="small" style={{ width: '100%', textAlign: 'left' }}>
                        <div className="flex justify-between items-center text-gray-500">
                            <Text>已答题目：{answeredCount}</Text>
                        </div>
                        <div className="flex justify-between items-center text-gray-500">
                            <Text>总题目数：{totalCount}</Text>
                        </div>
                    </Space>
                </div>

                {!isComplete && (
                    <Alert
                        message="提示"
                        description="建议您检查并完成所有题目后再提交，以获得更准确的评分。"
                        type="warning"
                        showIcon
                        style={{ marginTop: 20, textAlign: 'left' }}
                        icon={<InfoCircleOutlined />}
                    />
                )}
            </div>
        </Modal>
    );
};

export default SubmitConfirmModal;
