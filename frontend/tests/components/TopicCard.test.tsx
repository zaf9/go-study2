import "@testing-library/jest-dom";
import { render, screen, fireEvent } from "@testing-library/react";
import TopicCard from "@/components/learning/TopicCard";
import { TopicSummary } from "@/types/learning";

const pushMock = jest.fn();

jest.mock("next/navigation", () => ({
  useRouter: () => ({ push: pushMock }),
}));

const topic: TopicSummary = {
  key: "constants",
  title: "常量",
  summary: "常量基础",
  chapterCount: 5,
};

describe("TopicCard", () => {
  beforeEach(() => {
    pushMock.mockReset();
  });

  it("渲染主题信息并支持点击跳转", () => {
    render(<TopicCard topic={topic} />);

    expect(screen.getByText("常量")).toBeInTheDocument();
    expect(screen.getByText("章节数：5")).toBeInTheDocument();

    fireEvent.click(screen.getByText("常量"));
    expect(pushMock).toHaveBeenCalledWith("/topics/constants");
  });
});


