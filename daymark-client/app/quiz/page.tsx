'use client'

import { useEffect } from 'react'
import { useRouter } from 'next/navigation'
import { useQuiz } from '@/lib/store'
import QuizContent from '@/components/quiz-content'

export default function QuizPage() {
  const router = useRouter()
  const { session } = useQuiz()

  useEffect(() => {
    if (!session) {
      router.push('/')
    }
  }, [session, router])

  if (!session) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-background text-foreground px-4">
        <p className="text-sm text-muted-foreground">Loading quiz…</p>
      </div>
    )
  }

  return (
    <div className="min-h-screen flex flex-col bg-background text-foreground">
      <main className="flex-1">
        <div className="max-w-3xl mx-auto px-4 py-10 md:py-12">
          <QuizContent />
        </div>
      </main>
    </div>
  )
}
