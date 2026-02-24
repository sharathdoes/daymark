'use client';

import { create } from 'zustand';
import { QuizSession, QuizResult, ApiError, QuestionResult } from '@/types';

interface QuizStore {
  // State
  currentSession: QuizSession | null;
  lastResult: QuizResult | null;
  loading: boolean;
  error: ApiError | null;
  loadingMessage: string;

  // Actions
  setCurrentSession: (session: QuizSession | null) => void;
  setLoading: (loading: boolean, message?: string) => void;
  setError: (error: ApiError | null) => void;
  recordAnswer: (questionIndex: number, answerIndex: number) => void;
  computeResult: () => QuizResult | null;
  clearSession: () => void;
}

export const useQuizStore = create<QuizStore>((set) => ({
  currentSession: null,
  lastResult: null,
  loading: false,
  error: null,
  loadingMessage: '',

  setCurrentSession: (session) => set({ currentSession: session }),
  setLoading: (loading, message = '') => set({ loading, loadingMessage: message }),
  setError: (error) => set({ error }),
  
  recordAnswer: (questionIndex, answerIndex) =>
    set((state) => {
      if (!state.currentSession) return state;
      const newAnswers = [...state.currentSession.answers];
      newAnswers[questionIndex] = answerIndex;
      return {
        currentSession: {
          ...state.currentSession,
          answers: newAnswers,
        },
      };
    }),

  clearSession: () =>
    set({
      currentSession: null,
      lastResult: null,
      loading: false,
      error: null,
      loadingMessage: '',
    }),

  computeResult: () => {
    let computed: QuizResult | null = null;
    set((state) => {
      const session = state.currentSession;
      if (!session) return state;

      const total = session.questions.length;
      let score = 0;

      const questions: QuestionResult[] = session.questions.map((q, index) => {
        const userAnswer = session.answers[index];
        const isCorrect = userAnswer !== null && userAnswer === q.correctAnswer;
        if (isCorrect) score++;

        return {
          id: q.id,
          text: q.text,
          userAnswer,
          correctAnswer: q.correctAnswer,
          options: q.options,
          explanation: q.explanation,
          isCorrect,
        };
      });

      const percentage = total > 0 ? (score / total) * 100 : 0;

      computed = {
        sessionId: session.id,
        categoryName: '',
        difficulty: session.difficulty,
        score,
        totalQuestions: total,
        percentage,
        questions,
      };

      return {
        ...state,
        lastResult: computed,
      };
    });

    return computed;
  },
}));
