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
  question_id: string; // 改为string以避免大整数精度问题
  is_correct: boolean;
  correct_answers: string[];
  explanation: string;
  type: string;       // 题型：single/multiple
  difficulty: string; // 难度：easy/medium/hard
}

export interface QuizSubmitResult {
  score: number;
  total_questions: number;
  correct_answers: number;
  passed: boolean;
  details: QuizAnswerDetail[];
}

export interface QuizQuestion {
  id: string; // 改为string以避免大整数精度问题
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
  type: string;       // 题型：single/multiple
  difficulty: string; // 难度：easy/medium/hard
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
