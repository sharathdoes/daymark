'use client'

import { useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'
import Link from 'next/link'
import { getDailyQuiz } from '@/lib/api'
import { useAuth, useQuiz } from '@/lib/store'
import { Quiz, QuizSession } from '@/lib/types'
import Header from '@/components/header'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { Button } from '@/components/ui/button'

export default function DailyQuizPage() {
  const router = useRouter()
  const { isAuthenticated } = useAuth()
  const { setSession } = useQuiz()

  const [quiz, setQuiz] = useState<Quiz | null>(null)
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState('')
  const [isStarting, setIsStarting] = useState(false)

  const today = new Date().toLocaleDateString('en-IN', {
    weekday: 'long',
    day: 'numeric',
    month: 'long',
    year: 'numeric',
  })

  useEffect(() => {
    const load = async () => {
      try {
        const data = await getDailyQuiz()
        setQuiz(data)
      } catch (err) {
        setError('Failed to load today\'s quiz.')
        console.error(err)
      } finally {
        setIsLoading(false)
      }
    }
    load()
  }, [])

  const handleStart = () => {
    if (!quiz) return

    if (!isAuthenticated) {
      router.push('/login?redirect=/daily')
      return
    }

    setIsStarting(true)
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
  }

  return (
    <div className="min-h-screen flex flex-col bg-background text-foreground">

      <main className="flex-1">
        <div className="max-w-2xl mx-auto px-4 py-12 md:py-16">
          {/* Header */}
          <section className="mb-10 text-center">
            <p className="text-xs font-medium tracking-[0.2em] text-muted-foreground uppercase mb-3">
              {today}
            </p>
            <h1 className="text-3xl md:text-4xl font-semibold tracking-tight text-balance mb-3">
              Quiz of the Day
            </h1>
            <p className="text-sm md:text-base text-muted-foreground max-w-md mx-auto">
              One quiz, every day at 6 AM — pulled from real articles across all categories.
              No setup needed. Just start.
            </p>
          </section>

          {/* Content */}
          {isLoading ? (
            <div className="flex justify-center py-16">
              <p className="text-sm text-muted-foreground animate-pulse">Loading today&apos;s quiz…</p>
            </div>
          ) : error ? (
            <Card className="bg-destructive/5 border-destructive/30">
              <CardContent className="py-4 text-sm text-destructive text-center">{error}</CardContent>
            </Card>
          ) : quiz === null ? (
            <Card>
              <CardContent className="py-12 text-center space-y-3">
                <div className="text-4xl mb-2">⏳</div>
                <p className="font-medium">No quiz yet for today.</p>
                <p className="text-sm text-muted-foreground">
                  Today&apos;s quiz is generated every day at <strong>6:00 AM IST</strong>.<br />
                  Check back soon!
                </p>
                <div className="pt-2">
                  <Button variant="outline" asChild>
                    <Link href="/">Build your own quiz instead</Link>
                  </Button>
                </div>
              </CardContent>
            </Card>
          ) : (
            <Card>
              <CardHeader>
                <div className="flex items-center justify-between flex-wrap gap-2">
                  <div>
                    <CardTitle className="text-lg">{quiz.title}</CardTitle>
                    <CardDescription className="mt-0.5">
                      {quiz.questions.length} questions · {quiz.difficulty} difficulty
                    </CardDescription>
                  </div>
                  <span className="text-xs font-medium bg-primary/10 text-primary border border-primary/20 rounded-full px-3 py-1">
                    Daily Special
                  </span>
                </div>
              </CardHeader>

              <CardContent className="space-y-6">
                {/* Stats row */}
                <div className="grid grid-cols-3 gap-3 text-center">
                  {[
                    { label: 'Questions', value: quiz.questions.length },
                    { label: 'Difficulty', value: quiz.difficulty.charAt(0).toUpperCase() + quiz.difficulty.slice(1) },
                    { label: 'Categories', value: quiz.category_ids?.length ?? '—' },
                  ].map((stat) => (
                    <div
                      key={stat.label}
                      className="rounded-lg border border-border bg-muted/40 px-3 py-4"
                    >
                      <p className="text-xl font-semibold">{stat.value}</p>
                      <p className="text-xs text-muted-foreground mt-0.5">{stat.label}</p>
                    </div>
                  ))}
                </div>

                {/* CTA */}
                <div className="flex flex-col items-center gap-3 pt-2">
                  <Button
                    size="lg"
                    className="w-full max-w-xs shadow-sm"
                    onClick={handleStart}
                    disabled={isStarting}
                  >
                    {isStarting ? 'Starting…' : "Start Today's Quiz →"}
                  </Button>

                  {!isAuthenticated && (
                    <p className="text-xs text-muted-foreground text-center">
                      <Link href="/login" className="underline underline-offset-4">Sign in</Link>{' '}
                      to track your daily streak and results.
                    </p>
                  )}
                </div>
              </CardContent>
            </Card>
          )}
        </div>
      </main>
    </div>
  )
}
