import { render, screen } from "@testing-library/react";
import QuizResultView from "@/components/quiz/QuizResult";
import { QuizResult } from "@/types/quiz";

const result: QuizResult = {
  score: 2,
  total: 3,
  correctIds: ["q1", "q2"],
  wrongIds: ["q3"],
  submittedAt: new Date().toISOString(),
  durationMs: 800,
};

describe("QuizResultView", () => {
  it("shows score and details", () => {
    render(<QuizResultView result={result} />);
    expect(screen.getByText(/2\/3/)).toBeInTheDocument();
    expect(screen.getByText(/q1/)).toBeInTheDocument();
    expect(screen.getByText(/q3/)).toBeInTheDocument();
  });
});

