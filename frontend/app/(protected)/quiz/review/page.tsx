import { Suspense } from "react";
import QuizReviewPageClient from "./QuizReviewPageClient";
import Loading from "@/components/common/Loading";

export default function QuizReviewPage() {
    return (
        <Suspense fallback={<Loading />}>
            <QuizReviewPageClient />
        </Suspense>
    );
}
