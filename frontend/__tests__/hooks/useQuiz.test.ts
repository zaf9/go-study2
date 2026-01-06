/**
 * T062: useQuiz Hook 单元测试
 * 测试useQuiz Hook的加载和提交逻辑
 */

import { renderHook, act, waitFor } from "@testing-library/react";
import useQuiz from "@/hooks/useQuiz";
import * as quizService from "@/services/quizService";

// Mock the quiz service
jest.mock("@/services/quizService");

const mockQuizService = quizService as jest.Mocked<typeof quizService>;

describe("useQuiz Hook", () => {
    const mockSession = {
        sessionId: "test-session-123",
        topic: "variables",
        chapter: "storage",
        questions: [
            {
                id: 1,
                type: "single",
                difficulty: "easy",
                question: "What is the default value of int?",
                options: [
                    { id: "A", label: "0" },
                    { id: "B", label: "1" },
                    { id: "C", label: "nil" },
                    { id: "D", label: "undefined" },
                ],
                codeSnippet: null,
            },
            {
                id: 2,
                type: "multiple",
                difficulty: "medium",
                question: "Which are value types?",
                options: [
                    { id: "A", label: "int" },
                    { id: "B", label: "string" },
                    { id: "C", label: "slice" },
                    { id: "D", label: "array" },
                ],
                codeSnippet: null,
            },
        ],
    };

    const mockSubmitResult = {
        score: 100,
        total_questions: 2,
        correct_answers: 2,
        passed: true,
        details: [
            {
                question_id: 1,
                is_correct: true,
                correct_answers: ["A"],
                explanation: "The default value is 0",
            },
            {
                question_id: 2,
                is_correct: true,
                correct_answers: ["A", "D"],
                explanation: "int and array are value types",
            },
        ],
    };

    beforeEach(() => {
        jest.clearAllMocks();
        // 清除 localStorage 以避免测试之间的状态泄漏
        localStorage.clear();
    });

    describe("加载测验会话", () => {
        it("应该成功加载测验会话和题目", async () => {
            mockQuizService.useQuizSession = jest.fn(() => ({
                data: mockSession,
                error: undefined,
                isLoading: false,
                mutate: jest.fn(),
            })) as any;

            const { result } = renderHook(() =>
                useQuiz("variables", "storage"),
            );

            expect(result.current.session).toEqual(mockSession);
            expect(result.current.questions).toHaveLength(2);
            expect(result.current.questions[0].id).toBe("1");
            expect(result.current.questions[0].stem).toBe(
                "What is the default value of int?",
            );
            expect(result.current.questions[1].multi).toBe(true);
        });

        it("加载中时应该显示loading状态", () => {
            mockQuizService.useQuizSession = jest.fn(() => ({
                data: undefined,
                error: undefined,
                isLoading: true,
                mutate: jest.fn(),
            })) as any;

            const { result } = renderHook(() =>
                useQuiz("variables", "storage"),
            );

            expect(result.current.isLoading).toBe(true);
            expect(result.current.session).toBeUndefined();
            expect(result.current.questions).toEqual([]);
        });

        it("加载失败时应该显示错误", () => {
            const mockError = new Error("Failed to load quiz");
            mockQuizService.useQuizSession = jest.fn(() => ({
                data: undefined,
                error: mockError,
                isLoading: false,
                mutate: jest.fn(),
            })) as any;

            const { result } = renderHook(() =>
                useQuiz("variables", "storage"),
            );

            expect(result.current.error).toEqual(mockError);
            expect(result.current.questions).toEqual([]);
        });
    });

    describe("选择答案", () => {
        it("应该正确记录单选题答案", async () => {
            mockQuizService.useQuizSession = jest.fn(() => ({
                data: mockSession,
                error: undefined,
                isLoading: false,
                mutate: jest.fn(),
            })) as any;

            const { result } = renderHook(() =>
                useQuiz("variables", "storage"),
            );

            act(() => {
                result.current.selectAnswer("1", ["A"]);
            });

            expect(result.current.answers[1]).toEqual(["A"]);
            expect(result.current.answeredCount).toBe(1);
        });

        it("应该正确记录多选题答案", () => {
            mockQuizService.useQuizSession = jest.fn(() => ({
                data: mockSession,
                error: undefined,
                isLoading: false,
                mutate: jest.fn(),
            })) as any;

            const { result } = renderHook(() =>
                useQuiz("variables", "storage"),
            );

            act(() => {
                result.current.selectAnswer("2", ["A", "D"]);
            });

            expect(result.current.answers[2]).toEqual(["A", "D"]);
            expect(result.current.answeredCount).toBe(1);
        });

        it("应该允许修改已选择的答案", () => {
            mockQuizService.useQuizSession = jest.fn(() => ({
                data: mockSession,
                error: undefined,
                isLoading: false,
                mutate: jest.fn(),
            })) as any;

            const { result } = renderHook(() =>
                useQuiz("variables", "storage"),
            );

            act(() => {
                result.current.selectAnswer("1", ["A"]);
            });

            expect(result.current.answers[1]).toEqual(["A"]);

            act(() => {
                result.current.selectAnswer("1", ["B"]);
            });

            expect(result.current.answers[1]).toEqual(["B"]);
            expect(result.current.answeredCount).toBe(1);
        });

        it("应该正确统计已回答题目数量", () => {
            mockQuizService.useQuizSession = jest.fn(() => ({
                data: mockSession,
                error: undefined,
                isLoading: false,
                mutate: jest.fn(),
            })) as any;

            const { result } = renderHook(() =>
                useQuiz("variables", "storage"),
            );

            expect(result.current.answeredCount).toBe(0);

            act(() => {
                result.current.selectAnswer("1", ["A"]);
            });

            expect(result.current.answeredCount).toBe(1);

            act(() => {
                result.current.selectAnswer("2", ["A", "D"]);
            });

            expect(result.current.answeredCount).toBe(2);
        });
    });

    describe("提交测验", () => {
        it("应该成功提交测验并返回结果", async () => {
            mockQuizService.useQuizSession = jest.fn(() => ({
                data: mockSession,
                error: undefined,
                isLoading: false,
                mutate: jest.fn(),
            })) as any;

            mockQuizService.submitQuiz = jest
                .fn()
                .mockResolvedValue(mockSubmitResult);

            const { result } = renderHook(() =>
                useQuiz("variables", "storage"),
            );

            // 先选择答案
            act(() => {
                result.current.selectAnswer("1", ["A"]);
                result.current.selectAnswer("2", ["A", "D"]);
            });

            // 提交测验
            await act(async () => {
                await result.current.submit();
            });

            expect(mockQuizService.submitQuiz).toHaveBeenCalledWith({
                sessionId: "test-session-123",
                topic: "variables",
                chapter: "storage",
                durationMs: expect.any(Number),
                answers: [
                    { questionId: "1", userAnswers: ["A"] },
                    { questionId: "2", userAnswers: ["A", "D"] },
                ],
            });

            expect(result.current.result).toEqual(mockSubmitResult);
            expect(result.current.submitting).toBe(false);
        });

        it("提交过程中应该设置submitting状态", async () => {
            mockQuizService.useQuizSession = jest.fn(() => ({
                data: mockSession,
                error: undefined,
                isLoading: false,
                mutate: jest.fn(),
            })) as any;

            let resolveSubmit: (value: any) => void;
            const submitPromise = new Promise((resolve) => {
                resolveSubmit = resolve;
            });
            mockQuizService.submitQuiz = jest
                .fn()
                .mockReturnValue(submitPromise);

            const { result } = renderHook(() =>
                useQuiz("variables", "storage"),
            );

            act(() => {
                result.current.selectAnswer("1", ["A"]);
            });

            // 开始提交
            act(() => {
                void result.current.submit();
            });

            // 提交中应该为true
            await waitFor(() => {
                expect(result.current.submitting).toBe(true);
            });

            // 完成提交
            act(() => {
                resolveSubmit!(mockSubmitResult);
            });

            await waitFor(() => {
                expect(result.current.submitting).toBe(false);
            });
        });

        it("允许提交空答案", async () => {
            mockQuizService.useQuizSession = jest.fn(() => ({
                data: mockSession,
                error: undefined,
                isLoading: false,
                mutate: jest.fn(),
            })) as any;

            mockQuizService.submitQuiz = jest
                .fn()
                .mockResolvedValue(mockSubmitResult);

            const { result } = renderHook(() =>
                useQuiz("variables", "storage"),
            );

            // 不选择答案直接提交（代码允许空答案提交）
            await act(async () => {
                await result.current.submit();
            });

            // 应该调用提交，但答案列表为空
            expect(mockQuizService.submitQuiz).toHaveBeenCalledWith({
                sessionId: "test-session-123",
                topic: "variables",
                chapter: "storage",
                durationMs: expect.any(Number),
                answers: [],
            });
        });

        it("会话未加载时不应该提交", async () => {
            mockQuizService.useQuizSession = jest.fn(() => ({
                data: undefined,
                error: undefined,
                isLoading: true,
                mutate: jest.fn(),
            })) as any;

            mockQuizService.submitQuiz = jest.fn();

            const { result } = renderHook(() =>
                useQuiz("variables", "storage"),
            );

            await act(async () => {
                const res = await result.current.submit();
                expect(res).toBeNull();
            });

            expect(mockQuizService.submitQuiz).not.toHaveBeenCalled();
        });

        it("提交失败时应该抛出异常", async () => {
            mockQuizService.useQuizSession = jest.fn(() => ({
                data: mockSession,
                error: undefined,
                isLoading: false,
                mutate: jest.fn(),
            })) as any;

            const submitError = new Error("Submit failed");
            mockQuizService.submitQuiz = jest
                .fn()
                .mockRejectedValue(submitError);

            const { result } = renderHook(() =>
                useQuiz("variables", "storage"),
            );

            act(() => {
                result.current.selectAnswer("1", ["A"]);
            });

            await expect(
                act(async () => {
                    await result.current.submit();
                }),
            ).rejects.toThrow("Submit failed");

            expect(result.current.submitting).toBe(false);
        });

        it("防止重复提交 - 提交中时再次提交应该被忽略", async () => {
            mockQuizService.useQuizSession = jest.fn(() => ({
                data: mockSession,
                error: undefined,
                isLoading: false,
                mutate: jest.fn(),
            })) as any;

            let resolveSubmit: (value: any) => void;
            const submitPromise = new Promise((resolve) => {
                resolveSubmit = resolve;
            });
            mockQuizService.submitQuiz = jest
                .fn()
                .mockReturnValue(submitPromise);

            const { result } = renderHook(() =>
                useQuiz("variables", "storage"),
            );

            act(() => {
                result.current.selectAnswer("1", ["A"]);
            });

            // 第一次提交
            act(() => {
                void result.current.submit();
            });

            await waitFor(() => {
                expect(result.current.submitting).toBe(true);
            });

            // 尝试第二次提交（应该被忽略）
            await act(async () => {
                const res = await result.current.submit();
                expect(res).toBeNull();
            });

            // submitQuiz应该只被调用一次
            expect(mockQuizService.submitQuiz).toHaveBeenCalledTimes(1);

            // 完成第一次提交
            act(() => {
                resolveSubmit!(mockSubmitResult);
            });

            await waitFor(() => {
                expect(result.current.submitting).toBe(false);
            });
        });
    });

    describe("重置测验", () => {
        it("应该清除所有答案和结果", async () => {
            const mockMutate = jest.fn();
            mockQuizService.useQuizSession = jest.fn(() => ({
                data: mockSession,
                error: undefined,
                isLoading: false,
                mutate: mockMutate,
            })) as any;

            mockQuizService.submitQuiz = jest
                .fn()
                .mockResolvedValue(mockSubmitResult);

            const { result } = renderHook(() =>
                useQuiz("variables", "storage"),
            );

            // 选择答案并提交
            act(() => {
                result.current.selectAnswer("1", ["A"]);
                result.current.selectAnswer("2", ["A", "D"]);
            });

            await act(async () => {
                await result.current.submit();
            });

            expect(result.current.answeredCount).toBe(2);
            expect(result.current.result).toEqual(mockSubmitResult);

            // 重置
            act(() => {
                result.current.reset();
            });

            expect(result.current.answers).toEqual({});
            expect(result.current.answeredCount).toBe(0);
            expect(result.current.result).toBeNull();
            expect(mockMutate).toHaveBeenCalled();
        });
    });
});
