'use client'

import React, { createContext, useContext, useState, ReactNode, useCallback } from 'react'
import { User, QuizSession, Quiz, Question } from './types'
import { getToken, setToken, clearToken, getMe } from './api'

// Auth Store Types
export type AuthStore = {
  user: User | null
  token: string | null
  isAuthenticated: boolean
  isLoading: boolean
  setAuth: (user: User, token: string) => void
  clearAuth: () => void
  rehydrate: () => Promise<void>
  setUser: (user: User | null) => void
}

// Quiz Store Types
export type QuizStore = {
  session: QuizSession | null
  loading: boolean
  loadingMessage: string
  error: string | null
  setSession: (session: QuizSession) => void
  answerQuestion: (index: number, answer: number) => void
  nextQuestion: () => void
  resetSession: () => void
  setError: (error: string | null) => void
  decrementTimer: () => void
}

// Create contexts
const AuthContext = createContext<AuthStore | undefined>(undefined)
const QuizContext = createContext<QuizStore | undefined>(undefined)
const ThemeContext = createContext<{ isDark: boolean; setIsDark: (v: boolean) => void } | undefined>(undefined)

// Auth Provider
export function AuthProvider({
  children,
  themeToggle,
}: {
  children: ReactNode
  themeToggle?: { isDark: boolean; setIsDark: (v: boolean) => void }
}) {
  const [user, setUser] = useState<User | null>(null)
  const [token, setTokenState] = useState<string | null>(null)
  const [isLoading, setIsLoading] = useState(false)
  const [session, setSessionState] = useState<QuizSession | null>(null)
  const [quizLoading, setQuizLoading] = useState(false)
  const [loadingMessage, setLoadingMessage] = useState('')
  const [error, setError] = useState<string | null>(null)

  // Initialize - rehydrate token from localStorage
  React.useEffect(() => {
    const storedToken = getToken()
    if (storedToken) {
      setTokenState(storedToken)
      setIsLoading(true)
      getMe()
        .then(userData => setUser(userData))
        .catch(() => {
          clearToken()
          setTokenState(null)
        })
        .finally(() => setIsLoading(false))
    }
  }, [])

  // Auth methods
  const setAuth = useCallback((user: User, token: string) => {
    setUser(user)
    setTokenState(token)
    setToken(token)
  }, [])

  const clearAuth = useCallback(() => {
    setUser(null)
    setTokenState(null)
    clearToken()
  }, [])

  const rehydrate = useCallback(async () => {
    const storedToken = getToken()
    if (storedToken) {
      setTokenState(storedToken)
      try {
        const userData = await getMe()
        setUser(userData)
      } catch {
        clearAuth()
      }
    }
    setIsLoading(false)
  }, [clearAuth])

  // Quiz methods
  const setSession = useCallback((newSession: QuizSession) => {
    setSessionState(newSession)
  }, [])

  const answerQuestion = useCallback((index: number, answer: number) => {
    setSessionState(prev => {
      if (!prev) return prev
      const newAnswers = [...prev.answers]
      newAnswers[index] = answer
      
      const isCorrect = answer === prev.quiz.questions[index].answer
      const newScore = isCorrect ? prev.score + 1 : prev.score

      return {
        ...prev,
        answers: newAnswers,
        score: newScore,
      }
    })
  }, [])

  const nextQuestion = useCallback(() => {
    setSessionState(prev => {
      if (!prev) return prev
      const isLastQuestion = prev.currentIndex === prev.quiz.questions.length - 1
      return {
        ...prev,
        currentIndex: isLastQuestion ? prev.currentIndex : prev.currentIndex + 1,
        completed: isLastQuestion,
      }
    })
  }, [])

  const resetSession = useCallback(() => {
    setSessionState(null)
  }, [])

  const decrementTimer = useCallback(() => {
    setSessionState(prev => {
      if (!prev || prev.timerSeconds <= 0) return prev
      return {
        ...prev,
        timerSeconds: prev.timerSeconds - 1
      }
    })
  }, [])

  const authValue: AuthStore = {
    user,
    token,
    isAuthenticated: !!user,
    isLoading,
    setAuth,
    clearAuth,
    rehydrate,
    setUser,
  }

  const quizValue: QuizStore = {
    session,
    loading: quizLoading,
    loadingMessage,
    error,
    setSession,
    answerQuestion,
    nextQuestion,
    resetSession,
    setError,
    decrementTimer,
  }

  return (
    <AuthContext.Provider value={authValue}>
      <QuizContext.Provider value={quizValue}>
        <ThemeContext.Provider value={themeToggle || { isDark: true, setIsDark: () => {} }}>
          {children}
        </ThemeContext.Provider>
      </QuizContext.Provider>
    </AuthContext.Provider>
  )
}

// Hooks
export function useAuth(): AuthStore {
  const context = useContext(AuthContext)
  if (context === undefined) {
    throw new Error('useAuth must be used within AuthProvider')
  }
  return context
}

export function useQuiz(): QuizStore {
  const context = useContext(QuizContext)
  if (context === undefined) {
    throw new Error('useQuiz must be used within AuthProvider')
  }
  return context
}

export function useTheme() {
  const context = useContext(ThemeContext)
  if (context === undefined) {
    throw new Error('useTheme must be used within AuthProvider')
  }
  return context
}
