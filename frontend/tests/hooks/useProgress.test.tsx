import React from "react";
import { render, screen, waitFor } from "@testing-library/react";
import { SWRConfig } from "swr";
import useProgress from "@/hooks/useProgress";
import { saveProgress, fetchAllProgress } from "@/lib/progress";

jest.mock("@/lib/progress", () => ({
  fetchAllProgress: jest.fn(),
  fetchProgressByTopic: jest.fn(),
  saveProgress: jest.fn(),
}));

const mockedSaveProgress = saveProgress as jest.MockedFunction<
  typeof saveProgress
>;
const mockedFetchAll = fetchAllProgress as jest.MockedFunction<
  typeof fetchAllProgress
>;

function Wrapper() {
  const { progress, recordProgress } = useProgress();
  React.useEffect(() => {
    void recordProgress({
      topic: "variables",
      chapter: "storage",
      status: "in_progress",
      position: "{}",
    });
  }, [recordProgress]);
  return <div data-testid="count">{progress.length}</div>;
}

describe("useProgress", () => {
  beforeEach(() => {
    mockedFetchAll.mockResolvedValue([
      {
        topic: "variables",
        chapter: "storage",
        status: "done",
        lastVisit: new Date().toISOString(),
      },
    ] as any);
    mockedSaveProgress.mockResolvedValue();
  });

  it("loads progress and triggers save", async () => {
    render(
      <SWRConfig value={{ provider: () => new Map(), dedupingInterval: 0 }}>
        <Wrapper />
      </SWRConfig>,
    );

    await waitFor(() => {
      expect(screen.getByTestId("count").textContent).toBe("1");
    });
    expect(mockedSaveProgress).toHaveBeenCalled();
  });
});
