"use client";

import { useMemo, useState, useRef } from "react";
import {
    submitQuiz,
    useQuizHistory as useQuizHistoryQuery,
    useQuizSession,
} from "@/services/quizService";
import { QuizItem, QuizSubmitResult } from "@/types/quiz";

export default function useQuiz(topic: string, chapter: string) {
    const [answers, setAnswers] = useState<Record<number, string[]>>({});
    const [result, setResult] = useState<QuizSubmitResult | null>(null);
    const [submitting, setSubmitting] = useState(false);
    const [startAt, setStartAt] = useState<number>(Date.now());

    const { data: session, error, isLoading, mutate } = useQuizSession(
        topic,
        chapter,
    );

    const selectAnswer = (id: string, choices: string[]) => {
        const questionId = Number(id);
        setAnswers((prev) => ({ ...prev, [questionId]: choices }));
    };

    // 使用 ref 来立即锁定提交状态，防止闭包陈旧导致的并发提交问题
    const isSubmittingRef = useRef(false);

    const submit = async () => {
        // 双重检查：ref (同步) 和 state (渲染)
        if (isSubmittingRef.current || submitting || !session || !session.sessionId || (session.questions ?? []).length === 0)
            return null;

        isSubmittingRef.current = true;
        setSubmitting(true);

        try {
            const durationMs = Date.now() - startAt;
            const payload = {
                sessionId: session.sessionId,
                topic,
                chapter,
                durationMs,
                answers: Object.entries(answers).map(([id, choices]) => ({
                    questionId: Number(id),
                    userAnswers: choices,
                })),
            };
            const res = await submitQuiz(payload);
            // 后端可能返回整章题目的判分详情；前端只需展示用户实际提交的题目详情，故在此进行过滤
            if (res && Array.isArray(res.details) && payload.answers.length > 0) {
                const answeredIds = new Set<number>(
                    payload.answers.map((a: any) => Number(a.questionId)),
                );
                const filteredDetails = res.details.filter((d: any) =>
                    answeredIds.has(Number(d.question_id)),
                );
                setResult({ ...res, details: filteredDetails });
            } else {
                setResult(res);
            }
            return res;
        } catch (e: any) {
            console.error("Submit quiz failed:", e);
            // 如果是重复提交（409），我们也应该认为由于某种原因已经提交成功了（或者之前的请求成功了），
            // 但为了让 UI 正确响应，我们可能需要根据具体需求处理。
            // 这里抛出异常让 UI 层显示提示。
            throw e;
        } finally {
            // 注意：成功提交后通常不需要重置为 false 允许再次提交，除非是重置测试。
            // 但如果是报错了，或者逻辑允许重试，则需要重置。
            // 鉴于已有 result 状态控制显示结果页，这里重置是可以的。
            isSubmittingRef.current = false;
            setSubmitting(false);
        }
    };

    const reset = () => {
        setAnswers({});
        setResult(null);
        setStartAt(Date.now());
        void mutate();
    };

    const answeredCount = useMemo(() => Object.keys(answers).length, [answers]);

    return {
        session,
        questions:
            session?.questions?.map<QuizItem>((q) => ({
                id: String(q.id),
                stem: q.question,
                options: q.options,
                multi: q.type === "multiple",
                answer: [],
                type: q.type,
                difficulty: q.difficulty,
                codeSnippet: q.codeSnippet ?? undefined,
            })) ?? [],
        error,
        isLoading,
        answers,
        answeredCount,
        result,
        submitting,
        selectAnswer,
        submit,
        reset,
    };
}

export function useQuizHistory(topic?: string) {
    const { data, error, isLoading, mutate } = useQuizHistoryQuery(topic);
    return { history: data ?? [], error, isLoading, refresh: mutate };
}
