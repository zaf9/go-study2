import { fireEvent, render, screen } from "@testing-library/react";
import ProgressOverview from "@/components/progress/ProgressOverview";

describe("ProgressOverview", () => {
  it("renders overview stats and continue action", () => {
    const onContinue = jest.fn();
    render(
      <ProgressOverview
        overall={{
          progress: 75,
          completedChapters: 3,
          totalChapters: 4,
          studyDays: 5,
          totalStudyTime: 600,
        }}
        next={{
          topic: "variables",
          chapter: "storage",
          status: "in_progress",
          progress: 40,
        }}
        onContinue={onContinue}
      />,
    );

    expect(screen.getByText(/已完成章节/)).toBeInTheDocument();
    expect(screen.getByText(/继续学习：variables/)).toBeInTheDocument();
    fireEvent.click(screen.getByText(/前\s*往/));
    expect(onContinue).toHaveBeenCalled();
  });
});

