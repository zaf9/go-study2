"use client";

import { useMemo, useState } from "react";
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

  const submit = async () => {
    if (!session || !session.sessionId || (session.questions ?? []).length === 0)
      return null;
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
      setResult(res);
      return res;
    } finally {
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
