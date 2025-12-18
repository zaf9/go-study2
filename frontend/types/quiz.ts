export interface QuizOption {
  id: string;
  label: string;
}

export interface QuizItem {
  id: string;
  stem: string;
  options: QuizOption[];
  multi: boolean;
  answer: string[];
  explanation?: string;
  type?: string;
  difficulty?: string;
  codeSnippet?: string | null;
}

export interface QuizSubmitAnswer {
  id: string;
  choices: string[];
}

export interface QuizSubmitRequest {
  topic: string;
  chapter?: string;
  answers: QuizSubmitAnswer[];
}

export interface QuizResult {
  score: number;
  total: number;
  correctIds: string[];
  wrongIds: string[];
  submittedAt: string;
  durationMs: number;
}

export interface QuizAnswerDetail {
  question_id: number;
  is_correct: boolean;
  correct_answers: string[];
  explanation: string;
}

export interface QuizSubmitResult {
  score: number;
  total_questions: number;
  correct_answers: number;
  passed: boolean;
  details: QuizAnswerDetail[];
}

export interface QuizQuestion {
  id: number;
  type: string;
  difficulty: string;
  question: string;
  options: QuizOption[];
  codeSnippet?: string | null;
}

export interface QuizSessionPayload {
  sessionId: string;
  topic: string;
  chapter: string;
  questions: QuizQuestion[];
}

export interface QuizHistoryItem {
  id: number;
  sessionId?: string;
  topic: string;
  chapter?: string | null;
  score: number;
  totalQuestions?: number;
  correctAnswers?: number;
  passed?: boolean;
  completedAt?: string | null;
}

export interface QuizReviewMeta {
  sessionId: string;
  topic: string;
  chapter: string;
  score: number;
  passed: boolean;
  completedAt?: string | null;
}

export interface QuizReviewItem {
  questionId: number;
  stem: string;
  options: string[];
  userChoice: string;
  correctChoice: string;
  isCorrect: boolean;
  explanation: string;
}

export interface QuizReviewDetail {
  meta: QuizReviewMeta;
  items: QuizReviewItem[];
}

export interface QuizStats {
  total: number;
  byType: Record<string, number>;
  byDifficulty: Record<string, number>;
}
