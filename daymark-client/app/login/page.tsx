"use client"

import { useState, useEffect } from 'react'
import { useRouter, useSearchParams } from 'next/navigation'
import Link from 'next/link'
import { useAuth, useQuiz } from '@/lib/store'
import { signIn, loginWithGoogle, loginWithGithub, generateQuiz, getToken } from '@/lib/api'
import type { QuizSession } from '@/lib/types'
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'

export default function LoginPage() {
  const router = useRouter()
  const searchParams = useSearchParams()
  const { setAuth } = useAuth()
  const { setSession } = useQuiz()
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const [isLoading, setIsLoading] = useState(false)

  // After OAuth redirect, the user lands back here authenticated.
  // Read pendingQuiz from localStorage and auto-generate the quiz.
  const { isAuthenticated } = useAuth()
  useEffect(() => {
    if (!isAuthenticated) return
    const raw = localStorage.getItem('pendingQuiz')
    if (!raw) return
    localStorage.removeItem('pendingQuiz')
    const pending = JSON.parse(raw)
    setIsLoading(true)
    generateQuiz({
      category_ids: pending.categories,
      difficulty: pending.difficulty,
      number_of_questions: pending.count,
    }).then(quiz => {
      const session: QuizSession = {
        quiz,
        currentIndex: 0,
        answers: Array(quiz.questions.length).fill(null),
        score: 0,
        completed: false,
        timerOption: pending.timer ?? 'none',
        timerSeconds: pending.timer && pending.timer !== 'none' ? parseInt(pending.timer, 10) * 60 : 0,
      }
      setSession(session)
      router.push('/quiz')
    }).catch(() => {
      router.push('/')
    }).finally(() => setIsLoading(false))
  }, [isAuthenticated, router, setSession])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    setIsLoading(true)

    try {
      const response = await signIn({ email, password })
      setAuth(response.user, response.token)
      const categoriesParam = searchParams.get('categories')
      const difficultyParam = searchParams.get('difficulty')
      const countParam = searchParams.get('count')

      if (categoriesParam && difficultyParam) {
        const categoryIds = categoriesParam
          .split(',')
          .map(Number)
          .filter(id => !Number.isNaN(id))

        try {
          const quiz = await generateQuiz({
            category_ids: categoryIds,
            difficulty: difficultyParam,
            number_of_questions: countParam ? Number(countParam) : 10,
          })

          const session: QuizSession = {
            quiz,
            currentIndex: 0,
            answers: Array(quiz.questions.length).fill(null),
            score: 0,
            completed: false,
            timerOption: 'none',
            timerSeconds: 0,
          }

          setSession(session)
          router.push('/quiz')
        } catch (err) {
          setError(err instanceof Error ? err.message : 'Failed to start quiz')
          router.push('/')
        }
      } else {
        router.push('/')
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to sign in')
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-background px-4 py-10">
      <div className="w-full max-w-sm">
        <Card>
          <CardHeader className="space-y-1 text-center">
            <CardTitle className="text-2xl font-semibold tracking-tight">
              Sign in to Daymark
            </CardTitle>
            <CardDescription className="text-xs">
              Enter your credentials to continue where you left off.
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-6">
            {error && (
              <div className="rounded-md border border-destructive/40 bg-destructive/5 px-3 py-2 text-xs text-destructive">
                {error}
              </div>
            )}

            <form onSubmit={handleSubmit} className="space-y-4">
              <div className="space-y-1.5">
                <label className="text-xs font-medium text-muted-foreground" htmlFor="email">
                  Email
                </label>
                <Input
                  id="email"
                  type="email"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  placeholder="you@example.com"
                  required
                  autoComplete="email"
                />
              </div>

              <div className="space-y-1.5">
                <label className="text-xs font-medium text-muted-foreground" htmlFor="password">
                  Password
                </label>
                <Input
                  id="password"
                  type="password"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  placeholder="••••••••"
                  required
                  autoComplete="current-password"
                />
              </div>

              <Button
                type="submit"
                className="w-full mt-2"
                disabled={isLoading}
              >
                {isLoading ? 'Signing in…' : 'Sign in'}
              </Button>
            </form>

            <div className="relative my-2 text-center text-[10px] text-muted-foreground">
              <span className="bg-background px-2 relative z-10">
                Or continue with
              </span>
              <div className="absolute inset-x-0 top-1/2 -z-0 h-px bg-border" aria-hidden="true" />
            </div>

            <div className="space-y-2">
              <Button
                type="button"
                variant="outline"
                className="w-full justify-center text-xs"
                onClick={loginWithGoogle}
              >
                Continue with Google
              </Button>
              <Button
                type="button"
                variant="outline"
                className="w-full justify-center text-xs"
                onClick={loginWithGithub}
              >
                Continue with GitHub
              </Button>
            </div>

            <p className="pt-1 text-center text-xs text-muted-foreground">
              Don&apos;t have an account?{' '}
              <Link href="/signup" className="underline underline-offset-4">
                Sign up
              </Link>
            </p>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
