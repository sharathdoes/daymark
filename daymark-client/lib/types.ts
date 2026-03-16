// TypeScript types that mirror the Go backend models

export type Category = {
  id: number
  name: string
  slug: string
}

export type FeedSource = {
  id: number
  name: string
  url: string
  categories: Category[]
}

export type Question = {
  question: string
  options: string[]
  answer: number
  article_url: string
}

export type Quiz = {
  id: number
  title: string
  category_ids: number[]
  difficulty: string
  questions: Question[]
  created_at: string
}

export type User = {
  id: string
  name: string
  email: string
  provider: string
  provider_id: string
  avatar_url: string
  created_at: string
  updated_at: string
}

export type SignUpRequest = {
  name: string
  email: string
  password: string
}

export type SignInRequest = {
  email: string
  password: string
}

export type AuthResponse = {
  token: string
  user: User
}

export type GenerateQuizRequest = {
  category_ids: number[]
  difficulty: string
  number_of_questions: number
}

export type QuizSession = {
  quiz: Quiz
  currentIndex: number
  answers: (number | null)[]
  score: number
  completed: boolean
  timerOption: string // 'none', '5', '10'
  timerSeconds: number // Remaining time in seconds
}

export type SaveQuizResultRequest = {
  quiz_id: number
  score: number
  total_questions: number
  difficulty: string
  categories: string
}

export type UserQuizResult = {
  id: number
  user_id: number
  quiz_id: number
  score: number
  total_questions: number
  difficulty: string
  categories: string
  created_at: string
}
