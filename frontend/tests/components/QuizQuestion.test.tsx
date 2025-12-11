import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import QuizQuestionCard from "@/components/quiz/QuizQuestion";
import { QuizQuestion } from "@/types/quiz";

const question: QuizQuestion = {
  id: 1,
  type: "single",
  difficulty: "easy",
  question: "示例题",
  options: [
    { id: "A", label: "选项A" },
    { id: "B", label: "选项B" },
  ],
};

describe("QuizQuestionCard", () => {
  it("renders options and handles change", async () => {
    const onChange = jest.fn();
    render(<QuizQuestionCard question={question} value={[]} onChange={onChange} />);

    const option = screen.getByLabelText("选项A");
    await userEvent.click(option);
    expect(onChange).toHaveBeenCalledWith(["A"]);
  });
});

