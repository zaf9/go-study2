"use client";

import { useMemo, useState, useRef, useEffect } from "react";
import {
    submitQuiz,
    useQuizHistory as useQuizHistoryQuery,
    useQuizSession,
} from "@/services/quizService";
import { QuizItem, QuizSubmitResult } from "@/types/quiz";
import { mutate as globalMutate } from "swr";
import { progressKeys } from "@/services/progressService";

/**
 * 获取localStorage中的暂存答案key
 */
function getStorageKey(topic: string, chapter: string): string {
    return `quiz_answers_${topic}_${chapter}`;
}

/**
 * 从localStorage加载暂存的答案
 */
function loadStoredAnswers(topic: string, chapter: string): Record<string, string[]> {
    if (typeof window === 'undefined') return {};

    try {
        const key = getStorageKey(topic, chapter);
        const stored = localStorage.getItem(key);
        if (stored) {
            return JSON.parse(stored);
        }
    } catch (error) {
        console.warn('Failed to load stored quiz answers:', error);
    }
    return {};
}

/**
 * 保存答案到localStorage
 */
function saveAnswersToStorage(topic: string, chapter: string, answers: Record<string, string[]>): void {
    if (typeof window === 'undefined') return;

    try {
        const key = getStorageKey(topic, chapter);
        localStorage.setItem(key, JSON.stringify(answers));
    } catch (error) {
        console.warn('Failed to save quiz answers:', error);
    }
}

/**
 * 清除localStorage中的暂存答案
 */
function clearStoredAnswers(topic: string, chapter: string): void {
    if (typeof window === 'undefined') return;
    
    try {
        const key = getStorageKey(topic, chapter);
        localStorage.removeItem(key);
    } catch (error) {
        console.warn('Failed to clear stored quiz answers:', error);
    }
}

export default function useQuiz(topic: string, chapter: string) {
    const [answers, setAnswers] = useState<Record<string, string[]>>(() => {
        // 初始化时尝试加载暂存的答案
        return loadStoredAnswers(topic, chapter);
    });
    const [result, setResult] = useState<QuizSubmitResult | null>(null);
    const [submitting, setSubmitting] = useState(false);
    const [startAt, setStartAt] = useState<number>(Date.now());

    const { data: session, error, isLoading, mutate } = useQuizSession(
        topic,
        chapter,
    );

    // 检查是否有暂存的答案
    const hasStoredAnswers = useMemo(() => {
        return Object.keys(loadStoredAnswers(topic, chapter)).length > 0;
    }, [topic, chapter]);

    // 当答案改变时，自动保存到localStorage
    useEffect(() => {
        if (Object.keys(answers).length > 0 && !result) {
            saveAnswersToStorage(topic, chapter, answers);
        }
    }, [answers, topic, chapter, result]);

    // 当 topic 或 chapter 改变时，重新加载对应的暂存答案
    useEffect(() => {
        const stored = loadStoredAnswers(topic, chapter);
        setAnswers(stored);
        setResult(null); // 切换章节时重置结果
        setStartAt(Date.now()); // 切换章节时重置开始时间
    }, [topic, chapter]);

    const selectAnswer = (id: string, choices: string[]) => {
        // 直接使用字符串ID，不转换为number，避免大整数精度问题
        setAnswers((prev) => ({ ...prev, [id]: choices }));
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
                    questionId: id, // 保持为字符串，避免大整数精度问题
                    userAnswers: choices,
                })),
            };
            const res = await submitQuiz(payload);
            // 后端可能返回整章题目的判分详情；前端只需展示用户实际提交的题目详情，故在此进行过滤
            if (res && Array.isArray(res.details) && payload.answers.length > 0) {
                const answeredIds = new Set<string>(
                    payload.answers.map((a: any) => String(a.questionId)),
                );
                const filteredDetails = res.details.filter((d: any) =>
                    answeredIds.has(String(d.question_id)),
                );
                // 如果过滤后details为空，但原始details不为空，说明可能是ID类型匹配问题
                // 这种情况下保留原始details，确保解析功能可用
                if (filteredDetails.length === 0 && res.details.length > 0) {
                    console.warn('Details filtering resulted in empty array, using original details');
                    setResult({ ...res, details: res.details });
                } else {
                    setResult({ ...res, details: filteredDetails });
                }
            } else {
                setResult(res);
            }
            
            // User Story 4: 提交成功后刷新进度数据,确保所有页面显示一致
            void globalMutate(progressKeys.overview);
            void globalMutate(progressKeys.topic(topic));
            
            // 提交成功后清除localStorage中的暂存答案
            clearStoredAnswers(topic, chapter);
            
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
        clearStoredAnswers(topic, chapter); // 重置时清除暂存答案
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
        hasStoredAnswers, // 返回是否有暂存答案的标志
        selectAnswer,
        submit,
        reset,
    };
}

export function useQuizHistory(topic?: string) {
    const { data, error, isLoading, mutate } = useQuizHistoryQuery(topic);
    return { history: data ?? [], error, isLoading, refresh: mutate };
}
