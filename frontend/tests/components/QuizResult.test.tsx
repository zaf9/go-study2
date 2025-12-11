import { render, screen } from "@testing-library/react";
import QuizResultView from "@/components/quiz/QuizResult";
import { QuizSubmitResult } from "@/types/quiz";

const result: QuizSubmitResult = {
  score: 80,
  total_questions: 3,
  correct_answers: 2,
  passed: true,
  details: [],
};

describe("QuizResultView", () => {
  it("shows score and details", () => {
    render(<QuizResultView result={result} />);
    expect(screen.getByText(/80/)).toBeInTheDocument();
    expect(screen.getByText(/正确题数：2\/3/)).toBeInTheDocument();
  });
});
