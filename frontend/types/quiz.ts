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

export interface QuizHistoryItem {
  id: number;
  topic: string;
  chapter?: string | null;
  score: number;
  total: number;
  durationMs: number;
  createdAt: string;
}

