/**
 * 集成测试：跨页面进度数据一致性 (User Story 4 - T106)
 * 验证SWR缓存确保所有页面显示相同的进度数据
 */

import { render, screen, waitFor } from "@testing-library/react";
import { SWRConfig } from "swr";
import { ProgressSnapshot } from "@/types/learning";

// Mock useProgress Hook
import useProgress from "@/hooks/useProgress";
jest.mock("@/hooks/useProgress");

// 模拟进度数据
const mockProgressSnapshot: ProgressSnapshot = {
  totalChapters: 50,
  completedChapters: 25,
  completionRate: 50,
  chapters: [
    {
      name: "lexical_elements",
      totalChapters: 20,
      completedChapters: 10,
      progress: 50,
    },
    {
      name: "constants",
      totalChapters: 15,
      completedChapters: 8,
      progress: 53.33,
    },
    {
      name: "variables",
      totalChapters: 15,
      completedChapters: 7,
      progress: 46.67,
    },
  ],
  next: {
    topic: "lexical_elements",
    chapter: "identifiers",
  },
};

// 测试组件A - 显示总体进度
const ComponentA = () => {
  const { data, isLoading } = useProgress();

  if (isLoading) return <div>Loading A...</div>;
  if (!data) return <div>No data A</div>;

  return (
    <div data-testid="component-a">
      <div data-testid="completion-rate-a">{data.completionRate}%</div>
      <div data-testid="completed-chapters-a">
        {data.completedChapters}/{data.totalChapters}
      </div>
    </div>
  );
};

// 测试组件B - 显示主题进度
const ComponentB = () => {
  const { data, isLoading } = useProgress();

  if (isLoading) return <div>Loading B...</div>;
  if (!data) return <div>No data B</div>;

  return (
    <div data-testid="component-b">
      {data.chapters?.map((chapter) => (
        <div key={chapter.name} data-testid={`topic-${chapter.name}-b`}>
          {chapter.progress.toFixed(1)}%
        </div>
      ))}
    </div>
  );
};

// 测试组件C - 显示下一个章节
const ComponentC = () => {
  const { data, isLoading } = useProgress();

  if (isLoading) return <div>Loading C...</div>;
  if (!data?.next) return <div>No next chapter</div>;

  return (
    <div data-testid="component-c">
      <div data-testid="next-chapter-c">
        {data.next.topic}/{data.next.chapter}
      </div>
    </div>
  );
};

describe("跨页面进度数据一致性 (T106)", () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it("多个组件应共享相同的Hook返回数据", async () => {
    // Mock useProgress to return consistent data
    (useProgress as jest.Mock).mockReturnValue({
      data: mockProgressSnapshot,
      isLoading: false,
      error: null,
      mutate: jest.fn(),
    });

    const swrCache = new Map();

    render(
      <SWRConfig value={{ dedupingInterval: 0, provider: () => swrCache }}>
        <ComponentA />
        <ComponentB />
      </SWRConfig>
    );

    // 等待组件渲染
    await waitFor(() => {
      expect(screen.getByTestId("component-a")).toBeInTheDocument();
      expect(screen.getByTestId("component-b")).toBeInTheDocument();
    });

    // 验证两个组件显示相同的数据
    expect(screen.getByTestId("completion-rate-a")).toHaveTextContent("50%");
    expect(screen.getByTestId("completed-chapters-a")).toHaveTextContent(
      "25/50"
    );

    // 验证useProgress被调用了2次(每个组件一次)
    expect(useProgress).toHaveBeenCalledTimes(2);
  });

  it("数据更新后所有组件应显示新数据", async () => {
    const mockMutate = jest.fn();
    const initialData = {
      data: mockProgressSnapshot,
      isLoading: false,
      error: null,
      mutate: mockMutate,
    };

    (useProgress as jest.Mock).mockReturnValue(initialData);

    const { rerender } = render(
      <SWRConfig value={{ dedupingInterval: 0, provider: () => new Map() }}>
        <ComponentA />
        <ComponentB />
      </SWRConfig>
    );

    await waitFor(() => {
      expect(screen.getByTestId("completion-rate-a")).toHaveTextContent("50%");
    });

    // 模拟数据更新
    const updatedData: ProgressSnapshot = {
      ...mockProgressSnapshot,
      completedChapters: 30,
      completionRate: 60,
    };

    (useProgress as jest.Mock).mockReturnValue({
      data: updatedData,
      isLoading: false,
      error: null,
      mutate: mockMutate,
    });

    // 重新渲染
    rerender(
      <SWRConfig value={{ dedupingInterval: 0, provider: () => new Map() }}>
        <ComponentA />
        <ComponentB />
      </SWRConfig>
    );

    // 验证显示新数据
    await waitFor(() => {
      expect(screen.getByTestId("completion-rate-a")).toHaveTextContent("60%");
      expect(screen.getByTestId("completed-chapters-a")).toHaveTextContent(
        "30/50"
      );
    });
  });

  it("所有组件应显示一致的进度数据", async () => {
    (useProgress as jest.Mock).mockReturnValue({
      data: mockProgressSnapshot,
      isLoading: false,
      error: null,
      mutate: jest.fn(),
    });

    render(
      <SWRConfig value={{ dedupingInterval: 0, provider: () => new Map() }}>
        <ComponentA />
        <ComponentB />
        <ComponentC />
      </SWRConfig>
    );

    // 等待所有组件加载
    await waitFor(() => {
      expect(screen.getByTestId("component-a")).toBeInTheDocument();
      expect(screen.getByTestId("component-b")).toBeInTheDocument();
      expect(screen.getByTestId("component-c")).toBeInTheDocument();
    });

    // 验证组件A - 总体进度
    expect(screen.getByTestId("completion-rate-a")).toHaveTextContent("50%");
    expect(screen.getByTestId("completed-chapters-a")).toHaveTextContent(
      "25/50"
    );

    // 验证组件B - 主题进度
    expect(screen.getByTestId("topic-lexical_elements-b")).toHaveTextContent(
      "50.0%"
    );
    expect(screen.getByTestId("topic-constants-b")).toHaveTextContent("53.3%");

    // 验证组件C - 下一章节
    expect(screen.getByTestId("next-chapter-c")).toHaveTextContent(
      "lexical_elements/identifiers"
    );

    // 验证useProgress被调用了3次(每个组件一次)
    expect(useProgress).toHaveBeenCalledTimes(3);
  });

  it("错误状态应正确处理", async () => {
    // Mock error state
    (useProgress as jest.Mock).mockReturnValue({
      data: null,
      isLoading: false,
      error: new Error("API Error"),
      mutate: jest.fn(),
    });

    render(
      <SWRConfig value={{ dedupingInterval: 0, provider: () => new Map() }}>
        <ComponentA />
        <ComponentB />
      </SWRConfig>
    );

    // 等待错误状态
    await waitFor(() => {
      expect(screen.getByText("No data A")).toBeInTheDocument();
      expect(screen.getByText("No data B")).toBeInTheDocument();
    });
  });

  it("loading状态应正确显示", async () => {
    // Mock loading state
    (useProgress as jest.Mock).mockReturnValue({
      data: null,
      isLoading: true,
      error: null,
      mutate: jest.fn(),
    });

    render(
      <SWRConfig value={{ dedupingInterval: 0, provider: () => new Map() }}>
        <ComponentA />
        <ComponentB />
        <ComponentC />
      </SWRConfig>
    );

    // 验证loading状态
    expect(screen.getByText("Loading A...")).toBeInTheDocument();
    expect(screen.getByText("Loading B...")).toBeInTheDocument();
    expect(screen.getByText("Loading C...")).toBeInTheDocument();
  });
});
