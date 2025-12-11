import { fireEvent, render, screen } from "@testing-library/react";
import { SWRConfig } from "swr";
import useProgress from "@/hooks/useProgress";
import {
  updateProgress,
  useProgressOverview,
  useTopicProgress,
} from "@/services/progressService";

jest.mock("@/services/progressService", () => ({
  useProgressOverview: jest.fn(),
  useTopicProgress: jest.fn(),
  updateProgress: jest.fn(),
}));

const mockedUpdate = updateProgress as jest.MockedFunction<typeof updateProgress>;
const mockedOverview = useProgressOverview as jest.MockedFunction<
  typeof useProgressOverview
>;
const mockedTopic = useTopicProgress as jest.MockedFunction<
  typeof useTopicProgress
>;

function Wrapper() {
  const { overview, topicDetail, recordProgress } = useProgress("variables");
  return (
    <div>
      <div data-testid="overall">{overview?.overall.progress}</div>
      <div data-testid="topic">{topicDetail?.id}</div>
      <button
        onClick={() =>
          recordProgress({ topic: "variables", chapter: "storage" })
        }
      >
        save
      </button>
    </div>
  );
}

describe("useProgress", () => {
  beforeEach(() => {
    mockedOverview.mockReturnValue({
      data: {
        overall: {
          progress: 30,
          completedChapters: 1,
          totalChapters: 4,
          studyDays: 2,
          totalStudyTime: 120,
        },
        topics: [],
        next: null,
      },
      error: null,
      isLoading: false,
      mutate: jest.fn().mockResolvedValue(undefined),
    } as any);
    mockedTopic.mockReturnValue({
      data: {
        id: "variables",
        name: "Variables",
        weight: 25,
        progress: 50,
        completedChapters: 2,
        totalChapters: 4,
        chapters: [],
      },
      error: null,
      isLoading: false,
      mutate: jest.fn().mockResolvedValue(undefined),
    } as any);
    mockedUpdate.mockResolvedValue({ status: "ok" } as any);
  });

  it("exposes overview and triggers update", async () => {
    render(
      <SWRConfig value={{ provider: () => new Map(), dedupingInterval: 0 }}>
        <Wrapper />
      </SWRConfig>,
    );

    expect(screen.getByTestId("overall").textContent).toBe("30");
    expect(screen.getByTestId("topic").textContent).toBe("variables");
    fireEvent.click(screen.getByText("save"));
    expect(mockedUpdate).toHaveBeenCalledWith({
      topic: "variables",
      chapter: "storage",
    });
  });
});
