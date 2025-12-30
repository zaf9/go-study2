/**
 * 测试 useProgress Hook (User Story 4)
 * 验证进度数据的加载、缓存和刷新逻辑
 */

import { renderHook, waitFor } from "@testing-library/react";
import useProgress from "@/hooks/useProgress";
import * as progressService from "@/services/progressService";
import { ProgressSnapshot, TopicProgressDetail } from "@/types/learning";

// Mock SWR
jest.mock("swr", () => {
  const originalSWR = jest.requireActual("swr");
  return {
    __esModule: true,
    ...originalSWR,
    default: jest.fn((key, fetcher, options) => {
      // 使用原始 SWR 但允许我们控制数据
      return originalSWR.default(key, fetcher, {
        ...options,
        dedupingInterval: 0, // 禁用去重以便测试
      });
    }),
  };
});

// Mock progressService
jest.mock("@/services/progressService", () => ({
  updateProgress: jest.fn(),
  useProgressOverview: jest.fn(),
  useTopicProgress: jest.fn(),
}));

describe("useProgress Hook", () => {
  const mockOverviewData: ProgressSnapshot = {
    overall: {
      progress: 24.0,
      completedChapters: 12,
      totalChapters: 50,
      studyDays: 5,
      totalStudyTime: 3600,
    },
    topics: [
      {
        id: "variables",
        name: "Variables",
        weight: 30,
        progress: 33.3,
        completedChapters: 10,
        totalChapters: 30,
      },
    ],
    next: {
      topic: "constants",
      chapter: "iota",
      title: "Iota常量",
    },
  };

  const mockTopicData: TopicProgressDetail = {
    id: "variables",
    name: "Variables",
    weight: 30,
    progress: 33.3,
    completedChapters: 10,
    totalChapters: 30,
    chapters: [
      {
        id: "storage",
        name: "Storage",
        status: "completed",
        readDuration: 600,
        scrollProgress: 100,
        quizPassed: true,
      },
      {
        id: "pointer",
        name: "Pointer",
        status: "in_progress",
        readDuration: 120,
        scrollProgress: 50,
        quizPassed: false,
      },
    ],
  };

  beforeEach(() => {
    jest.clearAllMocks();
  });

  describe("T091: 基本功能测试", () => {
    it("应正确加载概览数据", async () => {
      // Mock useProgressOverview 返回数据
      (progressService.useProgressOverview as jest.Mock).mockReturnValue({
        data: mockOverviewData,
        error: null,
        isLoading: false,
        mutate: jest.fn(),
      });

      (progressService.useTopicProgress as jest.Mock).mockReturnValue({
        data: null,
        error: null,
        isLoading: false,
        mutate: jest.fn(),
      });

      const { result } = renderHook(() => useProgress());

      await waitFor(() => {
        expect(result.current.overview).toEqual(mockOverviewData);
        expect(result.current.isLoading).toBe(false);
        expect(result.current.error).toBeNull();
      });
    });

    it("应正确加载主题详情数据", async () => {
      (progressService.useProgressOverview as jest.Mock).mockReturnValue({
        data: mockOverviewData,
        error: null,
        isLoading: false,
        mutate: jest.fn(),
      });

      (progressService.useTopicProgress as jest.Mock).mockReturnValue({
        data: mockTopicData,
        error: null,
        isLoading: false,
        mutate: jest.fn(),
      });

      const { result } = renderHook(() => useProgress("variables"));

      await waitFor(() => {
        expect(result.current.topicDetail).toEqual(mockTopicData);
        expect(result.current.chapters).toEqual(mockTopicData.chapters);
        expect(result.current.isLoading).toBe(false);
      });
    });

    it("应正确处理加载状态", () => {
      (progressService.useProgressOverview as jest.Mock).mockReturnValue({
        data: undefined,
        error: null,
        isLoading: true,
        mutate: jest.fn(),
      });

      (progressService.useTopicProgress as jest.Mock).mockReturnValue({
        data: undefined,
        error: null,
        isLoading: false,
        mutate: jest.fn(),
      });

      const { result } = renderHook(() => useProgress());

      expect(result.current.isLoading).toBe(true);
      expect(result.current.overview).toBeUndefined();
    });

    it("应正确处理错误状态", () => {
      const mockError = new Error("API错误");

      (progressService.useProgressOverview as jest.Mock).mockReturnValue({
        data: undefined,
        error: mockError,
        isLoading: false,
        mutate: jest.fn(),
      });

      (progressService.useTopicProgress as jest.Mock).mockReturnValue({
        data: undefined,
        error: null,
        isLoading: false,
        mutate: jest.fn(),
      });

      const { result } = renderHook(() => useProgress());

      expect(result.current.error).toBe(mockError);
    });
  });

  describe("T092: SWR缓存一致性测试", () => {
    it("应支持手动刷新数据", async () => {
      const mockMutateOverview = jest.fn().mockResolvedValue(mockOverviewData);
      const mockMutateTopic = jest.fn().mockResolvedValue(mockTopicData);

      (progressService.useProgressOverview as jest.Mock).mockReturnValue({
        data: mockOverviewData,
        error: null,
        isLoading: false,
        mutate: mockMutateOverview,
      });

      (progressService.useTopicProgress as jest.Mock).mockReturnValue({
        data: mockTopicData,
        error: null,
        isLoading: false,
        mutate: mockMutateTopic,
      });

      const { result } = renderHook(() => useProgress("variables"));

      // 调用 refresh 方法
      await result.current.refresh();

      expect(mockMutateOverview).toHaveBeenCalled();
      expect(mockMutateTopic).toHaveBeenCalled();
    });

    it("应支持记录进度并自动刷新", async () => {
      const mockMutateOverview = jest.fn().mockResolvedValue(mockOverviewData);
      const mockMutateTopic = jest.fn().mockResolvedValue(mockTopicData);
      const mockUpdateProgress = jest
        .fn()
        .mockResolvedValue({ status: "in_progress" });

      (progressService.updateProgress as jest.Mock).mockImplementation(
        mockUpdateProgress,
      );

      (progressService.useProgressOverview as jest.Mock).mockReturnValue({
        data: mockOverviewData,
        error: null,
        isLoading: false,
        mutate: mockMutateOverview,
      });

      (progressService.useTopicProgress as jest.Mock).mockReturnValue({
        data: mockTopicData,
        error: null,
        isLoading: false,
        mutate: mockMutateTopic,
      });

      const { result } = renderHook(() => useProgress("variables"));

      // 记录进度
      await result.current.recordProgress({
        topic: "variables",
        chapter: "storage",
        readDuration: 120,
        scrollProgress: 50,
      });

      // 验证调用了 updateProgress
      expect(mockUpdateProgress).toHaveBeenCalledWith({
        topic: "variables",
        chapter: "storage",
        readDuration: 120,
        scrollProgress: 50,
      });

      // 验证自动刷新了缓存
      expect(mockMutateOverview).toHaveBeenCalled();
      expect(mockMutateTopic).toHaveBeenCalled();
    });

    it("无主题时不应刷新主题详情", async () => {
      const mockMutateOverview = jest.fn().mockResolvedValue(mockOverviewData);
      const mockMutateTopic = jest.fn();
      const mockUpdateProgress = jest
        .fn()
        .mockResolvedValue({ status: "in_progress" });

      (progressService.updateProgress as jest.Mock).mockImplementation(
        mockUpdateProgress,
      );

      (progressService.useProgressOverview as jest.Mock).mockReturnValue({
        data: mockOverviewData,
        error: null,
        isLoading: false,
        mutate: mockMutateOverview,
      });

      (progressService.useTopicProgress as jest.Mock).mockReturnValue({
        data: null,
        error: null,
        isLoading: false,
        mutate: mockMutateTopic,
      });

      const { result } = renderHook(() => useProgress());

      // 记录进度（无主题）
      await result.current.recordProgress({
        topic: "variables",
        chapter: "storage",
        readDuration: 120,
        scrollProgress: 50,
      });

      // 验证只刷新了概览
      expect(mockMutateOverview).toHaveBeenCalled();
      expect(mockMutateTopic).not.toHaveBeenCalled();
    });
  });

  describe("默认值处理", () => {
    it("应提供默认的 chapters 数组", () => {
      (progressService.useProgressOverview as jest.Mock).mockReturnValue({
        data: mockOverviewData,
        error: null,
        isLoading: false,
        mutate: jest.fn(),
      });

      (progressService.useTopicProgress as jest.Mock).mockReturnValue({
        data: null,
        error: null,
        isLoading: false,
        mutate: jest.fn(),
      });

      const { result } = renderHook(() => useProgress());

      expect(result.current.chapters).toEqual([]);
    });

    it("应提供默认的 next 值", () => {
      (progressService.useProgressOverview as jest.Mock).mockReturnValue({
        data: { ...mockOverviewData, next: null },
        error: null,
        isLoading: false,
        mutate: jest.fn(),
      });

      (progressService.useTopicProgress as jest.Mock).mockReturnValue({
        data: null,
        error: null,
        isLoading: false,
        mutate: jest.fn(),
      });

      const { result } = renderHook(() => useProgress());

      expect(result.current.next).toBeNull();
    });
  });
});
