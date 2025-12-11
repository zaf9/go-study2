import { fireEvent, render, screen } from "@testing-library/react";
import ProgressPage from "@/app/(protected)/progress/page";
import useProgress from "@/hooks/useProgress";
import { useRouter } from "next/navigation";

jest.mock("next/navigation", () => ({
  useRouter: jest.fn(),
}));

jest.mock("@/hooks/useProgress", () => ({
  __esModule: true,
  default: jest.fn(),
}));

const mockedUseProgress = useProgress as jest.MockedFunction<typeof useProgress>;
const mockedUseRouter = useRouter as jest.MockedFunction<typeof useRouter>;

describe("Progress page", () => {
  beforeEach(() => {
    mockedUseRouter.mockReturnValue({ push: jest.fn() } as any);
    mockedUseProgress.mockReturnValue({
      overview: {
        overall: {
          progress: 60,
          completedChapters: 2,
          totalChapters: 4,
          studyDays: 3,
          totalStudyTime: 500,
        },
        topics: [
          {
            id: "variables",
            name: "Variables",
            weight: 25,
            progress: 50,
            completedChapters: 2,
            totalChapters: 4,
            lastVisitAt: "",
          },
        ],
      },
      next: {
        topic: "variables",
        chapter: "storage",
        status: "in_progress",
        progress: 30,
      },
      isLoading: false,
      error: null,
    } as any);
  });

  it("renders overview and allows continue", () => {
    render(<ProgressPage />);
    expect(screen.getByText("学习进度")).toBeInTheDocument();
    expect(screen.getByText(/Variables/)).toBeInTheDocument();
    fireEvent.click(screen.getByText(/前\s*往/));
    const router = mockedUseRouter.mock.results[0].value;
    expect(router.push).toHaveBeenCalledWith("/topics/variables/storage");
  });
});

