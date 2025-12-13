"use client";

import { Alert, Button, Card, Col, Row, Statistic } from "antd";
import ProgressBar from "./ProgressBar";
import { ProgressStatuses } from "@/lib/progressStatus";
import { NextChapterHint, OverallProgress } from "@/types/learning";

/** 总览卡片属性：承载整体进度与“继续学习”提示。 */
interface ProgressOverviewProps {
  overall: OverallProgress;
  next?: NextChapterHint | null;
  onContinue?: (hint: NextChapterHint) => void;
}

export default function ProgressOverview({
  overall,
  next,
  onContinue,
}: ProgressOverviewProps) {
  return (
    <Card>
      <Row gutter={16}>
        <Col span={12}>
          <ProgressBar
              percent={overall.progress}
              status={overall.progress >= 100 ? ProgressStatuses.Completed : ProgressStatuses.InProgress}
              segments={12}
              label="整体进度"
            />
        </Col>
        <Col span={12}>
          <Row gutter={16}>
            <Col span={12}>
              <Statistic
                title="已完成章节"
                value={`${overall.completedChapters}/${overall.totalChapters}`}
              />
            </Col>
            <Col span={12}>
              <Statistic title="学习天数" value={overall.studyDays} />
            </Col>
            <Col span={12}>
              <Statistic
                title="学习总时长(秒)"
                value={overall.totalStudyTime}
              />
            </Col>
          </Row>
        </Col>
      </Row>
      {next && (
        <Alert
          className="mt-4"
          type="info"
          message={`继续学习：${next.topic} / ${next.chapter}`}
          description={`当前状态：${next.status}，进度 ${next.progress}%`}
          action={
            onContinue ? (
              <Button size="small" type="primary" onClick={() => onContinue(next)}>
                前往
              </Button>
            ) : undefined
          }
          showIcon
        />
      )}
    </Card>
  );
}

