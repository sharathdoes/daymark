export type Category = {
  ID: number;
  Name: string;
  Slug: string;
};
export interface Question {
  id: string;
  text: string;
  options: string[];
  correctAnswer: number;
  explanation: string;
}

export interface QuizSession {
  id: string;
  categoryId: string;
  difficulty: string;
  questions: Question[];
  currentQuestionIndex: number;
  answers: (number | null)[];
  score: number;
  completed: boolean;
  createdAt: string;
}

export interface QuizResult {
  sessionId: string;
  categoryName: string;
  difficulty: string;
  score: number;
  totalQuestions: number;
  percentage: number;
  questions: QuestionResult[];
}

export interface QuestionResult {
  id: string;
  text: string;
  userAnswer: number | null;
  correctAnswer: number;
  options: string[];
  explanation: string;
  isCorrect: boolean;
}

export interface ApiError {
  message: string;
  code?: string;
}
