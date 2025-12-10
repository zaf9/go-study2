'use client';

import useSWR from "swr";
import { useMemo, useState } from "react";
import { fetchQuizQuestions, fetchQuizHistory, submitQuiz } from "@/lib/quiz";
import { QuizHistoryItem, QuizItem, QuizResult } from "@/types/quiz";

export default function useQuiz(topic: string, chapter: string) {
  const [answers, setAnswers] = useState<Record<string, string[]>>({});
  const [result, setResult] = useState<QuizResult | null>(null);
  const [submitting, setSubmitting] = useState(false);

  const {
    data: questions,
    error,
    isLoading,
    mutate,
  } = useSWR<QuizItem[]>(["quiz", topic, chapter], () => fetchQuizQuestions(topic, chapter));

  const selectAnswer = (id: string, choices: string[]) => {
    setAnswers((prev) => ({ ...prev, [id]: choices }));
  };

  const submit = async () => {
    if (!questions || questions.length === 0) return null;
    setSubmitting(true);
    try {
      const payload = {
        topic,
        chapter,
        answers: Object.entries(answers).map(([id, choices]) => ({
          id,
          choices,
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
    void mutate();
  };

  const answeredCount = useMemo(() => Object.keys(answers).length, [answers]);

  return {
    questions: questions ?? [],
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
  const { data, error, isLoading, mutate } = useSWR<QuizHistoryItem[]>(
    ["quiz-history", topic || "all"],
    () => fetchQuizHistory(topic)
  );

  return {
    history: data ?? [],
    error,
    isLoading,
    refresh: mutate,
  };
}

