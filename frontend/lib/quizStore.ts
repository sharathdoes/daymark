'use client';

import { create } from 'zustand';
import { QuizSession, ApiError } from '@/types';

interface QuizStore {
  // State
  currentSession: QuizSession | null;
  loading: boolean;
  error: ApiError | null;
  loadingMessage: string;

  // Actions
  setCurrentSession: (session: QuizSession | null) => void;
  setLoading: (loading: boolean, message?: string) => void;
  setError: (error: ApiError | null) => void;
  recordAnswer: (questionIndex: number, answerIndex: number) => void;
  clearSession: () => void;
}

export const useQuizStore = create<QuizStore>((set) => ({
  currentSession: null,
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
      loading: false,
      error: null,
      loadingMessage: '',
    }),
}));
