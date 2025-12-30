/**
 * T119: QuizCenter组件
 * 测验中心容器组件,负责获取和展示测验历史列表
 */

"use client";

import { useState, useEffect } from "react";
import { Row, Col, Spin, Empty, Button, Select, Space, Typography, Alert } from "antd";
import { quizService } from "@/services/quizService";
import { QuizHistoryItem } from "@/types/quiz";
import QuizHistoryCard from "./QuizHistoryCard";
import { useRouter } from "next/navigation";

const { Text } = Typography;
const { Option } = Select;

export default function QuizCenter() {
  const router = useRouter();
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [historyList, setHistoryList] = useState<QuizHistoryItem[]>([]);
  const [filteredList, setFilteredList] = useState<QuizHistoryItem[]>([]);
  const [topicFilter, setTopicFilter] = useState<string>("");

  // 获取所有topic选项
  const topics = Array.from(new Set(historyList.map(item => item.topic)));

  // 加载测验历史
  useEffect(() => {
    loadHistory();
  }, []);

  // 应用过滤
  useEffect(() => {
    if (topicFilter) {
      setFilteredList(historyList.filter(item => item.topic === topicFilter));
    } else {
      setFilteredList(historyList);
    }
  }, [historyList, topicFilter]);

  async function loadHistory() {
    try {
      setLoading(true);
      setError(null);
      const data = await quizService.getQuizHistory();
      setHistoryList(data);
    } catch (err) {
      setError("加载测验历史失败,请稍后重试");
      console.error("Failed to load quiz history:", err);
    } finally {
      setLoading(false);
    }
  }

  function handleCardClick(sessionId: string) {
    router.push(`/quiz/history/${sessionId}`);
  }

  function handleClearFilter() {
    setTopicFilter("");
  }

  if (loading) {
    return (
      <div className="flex justify-center items-center min-h-[400px]" data-testid="loading-spinner">
        <Spin size="large" tip="加载中..." />
      </div>
    );
  }

  if (error) {
    return (
      <div data-testid="error-message">
        <Alert
          message="加载失败"
          description={error}
          type="error"
          showIcon
          action={
            <Button size="small" onClick={loadHistory}>
              重试
            </Button>
          }
        />
      </div>
    );
  }

  return (
    <div className="quiz-center" data-testid="quiz-center-component">
      {/* 过滤器 */}
      <div className="mb-6">
        <Space>
          <Text>筛选主题:</Text>
          <Select
            style={{ width: 200 }}
            placeholder="选择主题"
            value={topicFilter || undefined}
            onChange={setTopicFilter}
            allowClear
            onClear={handleClearFilter}
          >
            {topics.map(topic => (
              <Option key={topic} value={topic}>
                {topic}
              </Option>
            ))}
          </Select>
          {topicFilter && (
            <Button onClick={handleClearFilter}>清除筛选</Button>
          )}
        </Space>
      </div>

      {/* 测验历史列表 */}
      {filteredList.length === 0 ? (
        <Empty
          data-testid="empty-state"
          description="暂无测验记录"
          image={Empty.PRESENTED_IMAGE_SIMPLE}
        >
          <Button type="primary" onClick={() => router.push("/dashboard")}>
            开始学习
          </Button>
        </Empty>
      ) : (
        <Row gutter={[16, 16]}>
          {filteredList.map(item => (
            <Col key={item.sessionId} xs={24} sm={12} md={8} lg={6}>
              <QuizHistoryCard item={item} onClick={handleCardClick} />
            </Col>
          ))}
        </Row>
      )}
    </div>
  );
}
