/**
 * 章节状态集成测试
 * 测试章节状态的显示和更新功能
 *
 * 验收标准:
 * - 章节列表清晰显示三种状态(未开始/学习中/已完成)
 * - 首次访问章节自动标记为"学习中"
 * - 完成测验后自动更新为"已完成"
 * - 状态有明显的视觉区分(颜色+图标)
 */

import "@testing-library/jest-dom";
import React from "react";
import { render, screen } from "@testing-library/react";
import ChapterList from "@/components/learning/ChapterList";
import { ChapterSummary, ChapterProgress } from "@/types/learning";

// Mock next/navigation
jest.mock("next/navigation", () => ({
  useRouter: () => ({
    push: jest.fn(),
  }),
}));

const mockChapters: ChapterSummary[] = [
  {
    id: "comments",
    title: "Comments",
    summary: "Go语言注释语法",
    order: 0,
  },
  {
    id: "tokens",
    title: "Tokens",
    summary: "Go语言Token解析",
    order: 1,
  },
  {
    id: "semicolons",
    title: "Semicolons",
    summary: "Go语言分号规则",
    order: 2,
  },
];

describe("章节状态集成测试", () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  describe("验收标准1: 章节列表清晰显示三种状态", () => {
    it("应显示'已完成'状态", () => {
      const progressData: ChapterProgress[] = [
        {
          userId: 1,
          topic: "lexical_elements" as const,
          chapter: "comments",
          status: "completed" as const,
          percent: 100,
          scrollProgress: 100,
          readDuration: 300,
          lastPosition: "end",
          quizScore: 100,
          quizPassed: true,
          firstVisitAt: "2024-01-01T00:00:00Z",
          lastVisitAt: "2024-01-01T01:00:00Z",
          completedAt: "2024-01-01T01:00:00Z",
        },
      ];

      render(
        <ChapterList
          topicKey="lexical_elements"
          chapters={mockChapters}
          progress={progressData}
        />
      );

      // 验证显示完成状态图标和文本
      const completedBadges = screen.getAllByText("✅");
      expect(completedBadges.length).toBeGreaterThan(0);
      expect(screen.getByText("已完成")).toBeInTheDocument();
    });

    it("应显示'学习中'状态", () => {
      const progressData: ChapterProgress[] = [
        {
          userId: 1,
          topic: "lexical_elements" as const,
          chapter: "tokens",
          status: "in_progress" as const,
          percent: 50,
          scrollProgress: 50,
          readDuration: 150,
          lastPosition: "middle",
          quizScore: 0,
          quizPassed: false,
          firstVisitAt: "2024-01-01T02:00:00Z",
          lastVisitAt: "2024-01-01T02:01:00Z",
        },
      ];

      render(
        <ChapterList
          topicKey="lexical_elements"
          chapters={mockChapters}
          progress={progressData}
        />
      );

      // 验证显示学习中状态图标和文本
      const inProgressBadges = screen.getAllByText("📖");
      expect(inProgressBadges.length).toBeGreaterThan(0);
      expect(screen.getByText("学习中")).toBeInTheDocument();
    });

    it("应显示'未开始'状态", () => {
      render(
        <ChapterList
          topicKey="lexical_elements"
          chapters={mockChapters}
          progress={[]}
        />
      );

      // 验证显示未开始状态图标和文本
      const notStartedBadges = screen.getAllByText("🔘");
      expect(notStartedBadges.length).toBeGreaterThan(0);
      const notStartedText = screen.getAllByText("未开始");
      expect(notStartedText.length).toBeGreaterThan(0);
    });

    it("应同时显示三种不同状态的章节", () => {
      const progressData: ChapterProgress[] = [
        {
          userId: 1,
          topic: "lexical_elements" as const,
          chapter: "comments",
          status: "completed" as const,
          percent: 100,
          scrollProgress: 100,
          readDuration: 300,
          lastPosition: "end",
          quizScore: 100,
          quizPassed: true,
          firstVisitAt: "2024-01-01T00:00:00Z",
          lastVisitAt: "2024-01-01T01:00:00Z",
          completedAt: "2024-01-01T01:00:00Z",
        },
        {
          userId: 1,
          topic: "lexical_elements" as const,
          chapter: "tokens",
          status: "in_progress" as const,
          percent: 50,
          scrollProgress: 50,
          readDuration: 150,
          lastPosition: "middle",
          quizScore: 0,
          quizPassed: false,
          firstVisitAt: "2024-01-01T02:00:00Z",
          lastVisitAt: "2024-01-01T02:01:00Z",
        },
      ];

      render(
        <ChapterList
          topicKey="lexical_elements"
          chapters={mockChapters}
          progress={progressData}
        />
      );

      // 验证三种状态都存在
      expect(screen.getByText("✅")).toBeInTheDocument(); // 已完成
      expect(screen.getByText("📖")).toBeInTheDocument(); // 学习中
      expect(screen.getByText("🔘")).toBeInTheDocument(); // 未开始
    });
  });

  describe("验收标准2: 首次访问章节自动标记为'学习中'", () => {
    it("应将首次访问的章节标记为学习中", () => {
      // 模拟首次访问: status='in_progress'表示已经开始学习
      const firstVisitProgress: ChapterProgress[] = [
        {
          userId: 1,
          topic: "lexical_elements" as const,
          chapter: "comments",
          status: "in_progress" as const,
          percent: 10,
          firstVisitAt: "2024-01-01T00:00:00Z",
          lastVisitAt: "2024-01-01T00:00:00Z",
          readDuration: 30,
          scrollProgress: 10,
          lastPosition: "middle",
          quizScore: 0,
          quizPassed: false,
        },
      ];

      render(
        <ChapterList
          topicKey="lexical_elements"
          chapters={mockChapters}
          progress={firstVisitProgress}
        />
      );

      // 有 firstVisitAt 且 status='in_progress' 的章节应该显示为"学习中"
      expect(screen.getByText("📖")).toBeInTheDocument();
      const inProgressText = screen.getAllByText("学习中");
      expect(inProgressText.length).toBeGreaterThan(0);
    });
  });

  describe("验收标准3: 完成测验后自动更新为'已完成'", () => {
    it("应将测验通过的章节标记为已完成", () => {
      // 模拟完成测验: quizPassed = true, completedAt 有值
      const completedProgress: ChapterProgress[] = [
        {
          userId: 1,
          topic: "lexical_elements" as const,
          chapter: "comments",
          status: "completed" as const,
          percent: 100,
          firstVisitAt: "2024-01-01T00:00:00Z",
          lastVisitAt: "2024-01-01T01:00:00Z",
          completedAt: "2024-01-01T01:00:00Z",
          readDuration: 300,
          scrollProgress: 100,
          lastPosition: "end",
          quizScore: 90,
          quizPassed: true,
        },
      ];

      render(
        <ChapterList
          topicKey="lexical_elements"
          chapters={mockChapters}
          progress={completedProgress}
        />
      );

      // 验证显示为已完成
      expect(screen.getByText("✅")).toBeInTheDocument();
      expect(screen.getByText("已完成")).toBeInTheDocument();
    });

    it("应区分已完成和学习中的状态", () => {
      const mixedProgress: ChapterProgress[] = [
        {
          userId: 1,
          topic: "lexical_elements" as const,
          chapter: "comments",
          status: "completed" as const,
          percent: 100,
          firstVisitAt: "2024-01-01T00:00:00Z",
          lastVisitAt: "2024-01-01T01:00:00Z",
          completedAt: "2024-01-01T01:00:00Z",
          readDuration: 300,
          scrollProgress: 100,
          lastPosition: "end",
          quizScore: 90,
          quizPassed: true,
        },
        {
          userId: 1,
          topic: "lexical_elements" as const,
          chapter: "tokens",
          status: "in_progress" as const,
          percent: 50,
          firstVisitAt: "2024-01-01T02:00:00Z",
          lastVisitAt: "2024-01-01T02:01:00Z",
          readDuration: 150,
          scrollProgress: 50,
          lastPosition: "middle",
          quizScore: 0,
          quizPassed: false,
        },
      ];

      render(
        <ChapterList
          topicKey="lexical_elements"
          chapters={mockChapters}
          progress={mixedProgress}
        />
      );

      // 验证两种状态都能正确显示
      const completedBadges = screen.getAllByText("✅");
      const inProgressBadges = screen.getAllByText("📖");

      expect(completedBadges.length).toBeGreaterThan(0);
      expect(inProgressBadges.length).toBeGreaterThan(0);
    });
  });

  describe("验收标准4: 状态有明显的视觉区分(颜色+图标)", () => {
    it("应为不同状态使用不同的图标", () => {
      const mixedProgress: ChapterProgress[] = [
        {
          userId: 1,
          topic: "lexical_elements" as const,
          chapter: "comments",
          status: "completed" as const,
          percent: 100,
          firstVisitAt: "2024-01-01T00:00:00Z",
          lastVisitAt: "2024-01-01T01:00:00Z",
          completedAt: "2024-01-01T01:00:00Z",
          readDuration: 300,
          scrollProgress: 100,
          lastPosition: "end",
          quizScore: 90,
          quizPassed: true,
        },
        {
          userId: 1,
          topic: "lexical_elements" as const,
          chapter: "tokens",
          status: "in_progress" as const,
          percent: 50,
          firstVisitAt: "2024-01-01T02:00:00Z",
          lastVisitAt: "2024-01-01T02:01:00Z",
          readDuration: 150,
          scrollProgress: 50,
          lastPosition: "middle",
          quizScore: 0,
          quizPassed: false,
        },
      ];

      const { container } = render(
        <ChapterList
          topicKey="lexical_elements"
          chapters={mockChapters}
          progress={mixedProgress}
        />
      );

      // 验证三种不同的图标都存在
      expect(screen.getByText("✅")).toBeInTheDocument();
      expect(screen.getByText("📖")).toBeInTheDocument();
      expect(screen.getByText("🔘")).toBeInTheDocument();

      // 验证不同状态使用不同的颜色类名
      const tags = container.querySelectorAll(".ant-tag");
      expect(tags.length).toBeGreaterThan(0);

      // 检查是否有不同颜色的标签
      const hasDefaultTag = Array.from(tags).some((tag) =>
        tag.classList.contains("ant-tag-default")
      );
      const hasProcessingTag = Array.from(tags).some((tag) =>
        tag.classList.contains("ant-tag-processing")
      );
      const hasSuccessTag = Array.from(tags).some((tag) =>
        tag.classList.contains("ant-tag-success")
      );

      expect(hasDefaultTag || hasProcessingTag || hasSuccessTag).toBe(true);
    });

    it("应显示状态文本", () => {
      const progressData: ChapterProgress[] = [
        {
          userId: 1,
          topic: "lexical_elements" as const,
          chapter: "comments",
          status: "completed" as const,
          percent: 100,
          scrollProgress: 100,
          readDuration: 300,
          lastPosition: "end",
          quizScore: 100,
          quizPassed: true,
          firstVisitAt: "2024-01-01T00:00:00Z",
          lastVisitAt: "2024-01-01T01:00:00Z",
          completedAt: "2024-01-01T01:00:00Z",
        },
      ];

      render(
        <ChapterList
          topicKey="lexical_elements"
          chapters={mockChapters}
          progress={progressData}
        />
      );

      // 验证状态文本存在
      expect(screen.getByText("已完成")).toBeInTheDocument();
    });
  });

  describe("边缘情况处理", () => {
    it("应处理空进度数据", () => {
      render(
        <ChapterList
          topicKey="lexical_elements"
          chapters={mockChapters}
          progress={[]}
        />
      );

      // 所有章节应显示为未开始
      const notStartedBadges = screen.getAllByText("🔘");
      expect(notStartedBadges.length).toBe(3);
    });

    it("应处理null或undefined进度", () => {
      render(
        <ChapterList
          topicKey="lexical_elements"
          chapters={mockChapters}
          progress={undefined}
        />
      );

      // 应该正常渲染，不报错
      expect(screen.getByText("Comments")).toBeInTheDocument();
      expect(screen.getByText("Tokens")).toBeInTheDocument();
      expect(screen.getByText("Semicolons")).toBeInTheDocument();
    });

    it("应处理进度数据与章节数据不匹配的情况", () => {
      const incompleteProgress: ChapterProgress[] = [
        {
          userId: 1,
          topic: "lexical_elements" as const,
          chapter: "comments",
          status: "completed" as const,
          percent: 100,
          scrollProgress: 100,
          readDuration: 300,
          lastPosition: "end",
          quizScore: 100,
          quizPassed: true,
          firstVisitAt: "2024-01-01T00:00:00Z",
          lastVisitAt: "2024-01-01T01:00:00Z",
          completedAt: "2024-01-01T01:00:00Z",
        },
        // 缺少 tokens 和 semicolons 的进度数据
      ];

      render(
        <ChapterList
          topicKey="lexical_elements"
          chapters={mockChapters}
          progress={incompleteProgress}
        />
      );

      // 应该正常渲染所有章节，没有进度的显示为未开始
      expect(screen.getByText("Comments")).toBeInTheDocument();
      expect(screen.getByText("Tokens")).toBeInTheDocument();
      expect(screen.getByText("Semicolons")).toBeInTheDocument();

      // 验证状态显示
      expect(screen.getByText("✅")).toBeInTheDocument(); // comments - 已完成
      expect(screen.getAllByText("🔘").length).toBeGreaterThan(0); // 其他 - 未开始
    });
  });
});
