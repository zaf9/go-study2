"use client";

import { useEffect, useState } from "react";
import useSWR from "swr";
import { Select, Space, Typography, message } from "antd";
import QuizQuestionCard from "@/components/quiz/QuizQuestionCard";
import SubmitConfirmModal from "@/components/quiz/SubmitConfirmModal";
import QuizResultPage from "@/components/quiz/QuizResultPage";
import QuizNavigation from "@/components/quiz/QuizNavigation";
import AnswerExplanation from "@/components/quiz/AnswerExplanation";
import QuizSkeletonLoader from "@/components/quiz/QuizSkeletonLoader";
import ErrorMessage from "@/components/common/ErrorMessage";
import { fetchChapters } from "@/lib/learning";
import useQuiz from "@/hooks/useQuiz";

const { Title, Paragraph, Text } = Typography;

export default function QuizPageClient({
    params,
}: {
    params: { topic: string; chapter?: string };
}) {
    const topic = params.topic;
    const {
        data: chapters,
        error: chapterError,
        isLoading: chapterLoading,
    } = useSWR(["quiz-chapters", topic], () => fetchChapters(topic));
    const [selectedChapter, setSelectedChapter] = useState(params.chapter ?? "");
    const [currentIndex, setCurrentIndex] = useState(0);
    const [isConfirmOpen, setIsConfirmOpen] = useState(false);

    useEffect(() => {
        if (params.chapter) {
            setSelectedChapter(params.chapter);
            return;
        }
        if (chapters && chapters.length > 0 && !selectedChapter) {
            setSelectedChapter(chapters[0].id);
        }
    }, [chapters, selectedChapter, params.chapter]);

    const {
        questions,
        answers,
        selectAnswer,
        submit,
        result,
        isLoading: quizLoading,
        submitting,
        reset,
        answeredCount,
    } = useQuiz(topic, selectedChapter);

    if (chapterLoading) return <QuizSkeletonLoader />;
    if (chapterError)
        return (
            <ErrorMessage message="加载章节失败" description={chapterError.message} />
        );
    if (!selectedChapter) return <ErrorMessage message="暂无章节可用" />;

    const handleOpenConfirm = () => {
        setIsConfirmOpen(true);
    };

    const handleFinalSubmit = async () => {
        const res = await submit();
        if (res) {
            message.success("提交成功");
            setIsConfirmOpen(false);
        }
    };

    return (
        <div className="space-y-4">
            <Space direction="vertical" className="w-full">
                <div className="flex items-center justify-between bg-white p-4 rounded shadow-sm border border-gray-100">
                    <div>
                        <Title level={3} className="mb-1" style={{ fontSize: '1.5rem', fontWeight: 600 }}>
                            {topic} 测验
                        </Title>
                        <Paragraph type="secondary" style={{ marginBottom: 0 }}>
                            章节：<Text strong>{chapters?.find(c => c.id === selectedChapter)?.title || selectedChapter}</Text>
                        </Paragraph>
                    </div>
                    {!params.chapter && (
                        <Select
                            style={{ minWidth: 220 }}
                            value={selectedChapter}
                            onChange={(v) => {
                                setSelectedChapter(v);
                                reset();
                                setCurrentIndex(0);
                            }}
                            options={(chapters ?? []).map((c) => ({
                                label: c.title,
                                value: c.id,
                            }))}
                        />
                    )}
                </div>

                {quizLoading ? (
                    <QuizSkeletonLoader />
                ) : (
                    <Space direction="vertical" className="w-full" size="middle">
                        {questions.length > 0 && !result && (
                            <>
                                <QuizQuestionCard
                                    question={{
                                        id: Number(questions[currentIndex].id),
                                        type: questions[currentIndex].type ?? "single",
                                        difficulty: questions[currentIndex].difficulty ?? "easy",
                                        question: questions[currentIndex].stem,
                                        options: questions[currentIndex].options,
                                        codeSnippet: questions[currentIndex].codeSnippet,
                                    }}
                                    value={answers[Number(questions[currentIndex].id)] ?? []}
                                    onChange={(val) => selectAnswer(questions[currentIndex].id, val)}
                                />
                                <QuizNavigation
                                    current={currentIndex}
                                    total={questions.length}
                                    answered={answeredCount}
                                    onPrev={() => setCurrentIndex(Math.max(0, currentIndex - 1))}
                                    onNext={() =>
                                        setCurrentIndex(
                                            Math.min(questions.length - 1, currentIndex + 1),
                                        )
                                    }
                                    onSubmit={handleOpenConfirm}
                                    submitting={submitting}
                                />
                            </>
                        )}

                        {questions.length === 0 && !quizLoading && (
                            <ErrorMessage message="当前章节暂无测验" />
                        )}

                        {result && (
                            <div className="scale-in-center">
                                <QuizResultPage
                                    result={result}
                                    onRetry={reset}
                                    onReview={() => {
                                        const el = document.getElementById('quiz-explanation-section');
                                        el?.scrollIntoView({ behavior: 'smooth' });
                                    }}
                                />
                                {result.details && result.details.length > 0 && (
                                    <div id="quiz-explanation-section" style={{ marginTop: 24 }}>
                                        <Title level={4}>答题详解</Title>
                                        <AnswerExplanation details={result.details} />
                                    </div>
                                )}
                            </div>
                        )}
                    </Space>
                )}
            </Space>

            <SubmitConfirmModal
                open={isConfirmOpen}
                answeredCount={answeredCount}
                totalCount={questions.length}
                onConfirm={handleFinalSubmit}
                onCancel={() => setIsConfirmOpen(false)}
                submitting={submitting}
            />
        </div>
    );
}
