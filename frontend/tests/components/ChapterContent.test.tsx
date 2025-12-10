import "@testing-library/jest-dom";
import { render, screen } from "@testing-library/react";
import ChapterContent from "@/components/learning/ChapterContent";
import { ChapterContent as ChapterContentType } from "@/types/learning";

jest.mock("react-markdown", () => (props: any) => <div>{props.children}</div>);
jest.mock("remark-gfm", () => () => null);

describe("ChapterContent", () => {
  it("渲染章节标题与内容", () => {
    const content: ChapterContentType = {
      id: "intro",
      topicKey: "constants",
      title: "示例章节",
      markdown: "```go\nfmt.Println(\"hi\")\n```",
    };

    render(<ChapterContent content={content} />);

    expect(screen.getByText("示例章节")).toBeInTheDocument();
    expect(screen.getByText(/fmt\.Println\("hi"\)/)).toBeInTheDocument();
  });
});


