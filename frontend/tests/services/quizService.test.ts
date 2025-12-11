import { fetchQuizSession, submitQuiz } from "@/services/quizService";

jest.mock("@/lib/api", () => {
  const get = jest.fn();
  const post = jest.fn();
  return {
    __esModule: true,
    default: {
      get,
      post,
    },
  };
});

const mockedApi = jest.requireMock("@/lib/api").default as {
  get: jest.Mock;
  post: jest.Mock;
};

describe("quizService", () => {
  beforeEach(() => {
    mockedApi.get.mockReset();
    mockedApi.post.mockReset();
  });

  it("normalizes quiz session response", async () => {
    mockedApi.get.mockResolvedValue({
      sessionId: "s1",
      topic: "variables",
      chapter: "storage",
      questions: [
        {
          id: 1,
          type: "single",
          difficulty: "easy",
          question: "Q1",
          options: [{ id: "A", label: "opt" }],
        },
      ],
    });

    const session = await fetchQuizSession("variables", "storage");
    expect(session.sessionId).toBe("s1");
    expect(session.questions[0].options[0].label).toBe("opt");
  });

  it("calls submit endpoint with payload", async () => {
    mockedApi.post.mockResolvedValue({ score: 80 });
    const res = await submitQuiz({
      sessionId: "s1",
      topic: "variables",
      chapter: "storage",
      answers: [{ questionId: 1, userAnswers: ["A"] }],
    });
    expect(mockedApi.post).toHaveBeenCalled();
    expect(res.score).toBe(80);
  });
});

