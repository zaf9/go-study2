import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import QuizItem from "@/components/quiz/QuizItem";
import { QuizItem as QuizItemType } from "@/types/quiz";

const sample: QuizItemType = {
  id: "q1",
  stem: "示例题目",
  options: [
    { id: "A", label: "选项A" },
    { id: "B", label: "选项B" },
  ],
  multi: false,
  answer: ["A"],
};

describe("QuizItem", () => {
  it("renders options and triggers change", async () => {
    const handleChange = jest.fn();
    render(<QuizItem question={sample} value={[]} onChange={handleChange} />);

    const optionA = screen.getByLabelText("选项A");
    await userEvent.click(optionA);
    expect(handleChange).toHaveBeenCalledWith(["A"]);
  });
});

