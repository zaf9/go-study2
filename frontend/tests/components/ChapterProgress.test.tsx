import { fireEvent, render, screen } from "@testing-library/react";
import ChapterProgress from "@/components/progress/ChapterProgress";
import { ProgressStatuses } from "@/lib/progressStatus";

describe("ChapterProgress", () => {
  it("renders status and triggers resume", () => {
    const onResume = jest.fn();
    render(
      <ChapterProgress
        title="章节一"
        status={ProgressStatuses.InProgress}
        scrollProgress={60}
        readDuration={120}
        estimatedSeconds={300}
        lastVisitAt="2025-12-11T10:00:00Z"
        onResume={onResume}
      />,
    );

    expect(screen.getByText(/章节一/)).toBeInTheDocument();
    expect(screen.getByText(/已阅读 120 秒/)).toBeInTheDocument();
    fireEvent.click(screen.getByText("恢复到上次阅读位置"));
    expect(onResume).toHaveBeenCalled();
  });
});

