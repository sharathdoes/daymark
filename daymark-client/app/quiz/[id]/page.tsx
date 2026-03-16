'use client'

import { useEffect, useState } from 'react'
import { useParams, useRouter } from 'next/navigation'
import { useQuiz } from '@/lib/store'
import { getQuizById } from '@/lib/api'
import Header from '@/components/header'
import QuizContent from '@/components/quiz-content'
import { Card, CardContent } from '@/components/ui/card'
import { Loader2 } from 'lucide-react'
import { Quiz } from '@/lib/types'

export default function SharedQuizPage() {
  const router = useRouter()
  const params = useParams()
  const idStr = params.id as string
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const { session, setSession } = useQuiz()

  useEffect(() => {
    async function loadQuiz() {
      if (!idStr) return
      
      try {
        setLoading(true)
        const quizData: Quiz = await getQuizById(idStr)
        
        // Start the shared quiz immediately with no timer by default, since they are joining a generated room
        setSession({
          quiz: quizData,
          currentIndex: 0,
          answers: Array(quizData.questions.length).fill(null),
          score: 0,
          completed: false,
          timerOption: 'none',
          timerSeconds: 0,
        })
      } catch (err) {
        console.error('Failed to load shared quiz:', err)
        const message = err instanceof Error ? err.message : 'Failed to load quiz. It may not exist.'
        setError(message)
      } finally {
        setLoading(false)
      }
    }

    loadQuiz()
  }, [idStr, setSession])

  // Once session is loaded, QuizContent assumes control of the UI flow
  useEffect(() => {
    if (session?.completed) {
      router.push('/results')
    }
  }, [session?.completed, router])

  return (
    <div className="min-h-screen bg-background flex flex-col">
      <main className="flex-1 p-4 md:p-8">
        <div className="w-full max-w-4xl mx-auto">
          {loading ? (
            <Card className="min-h-[60vh] flex flex-col items-center justify-center space-y-4">
              <Loader2 className="h-8 w-8 animate-spin text-primary" />
              <p className="text-muted-foreground text-sm font-medium animate-pulse">
                Fetching shared quiz...
              </p>
            </Card>
          ) : error ? (
            <Card className="min-h-[40vh] flex flex-col items-center justify-center p-6 text-center space-y-4">
              <div className="mx-auto flex h-16 w-16 items-center justify-center rounded-full bg-destructive/10">
                <span className="text-2xl font-bold text-destructive">!</span>
              </div>
              <p className="text-destructive font-semibold">Error loading quiz</p>
              <p className="text-muted-foreground text-sm max-w-md">{error}</p>
            </Card>
          ) : session ? (
            <div className="pt-6 md:pt-10">
              <div className="mb-6 flex flex-col gap-1.5 px-1">
                <h1 className="text-2xl font-bold tracking-tight text-foreground md:text-3xl">
                  {session.quiz.title}
                </h1>
                <p className="text-sm font-medium text-muted-foreground capitalize">
                  Difficulty: {session.quiz.difficulty}
                </p>
              </div>
              <QuizContent />
            </div>
          ) : null}
        </div>
      </main>
    </div>
  )
}
