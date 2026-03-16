'use client'

import { useRouter } from 'next/navigation'
import { useQuiz } from '@/lib/store'
import { useEffect } from 'react'
import ResultsContent from '@/components/results-content'

export default function ResultsPage() {
  const router = useRouter()
  const { session, resetSession } = useQuiz()

  useEffect(() => {
    if (!session || !session.completed) {
      router.push('/')
    }
  }, [session, router])

  if (!session || !session.completed) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-background text-foreground px-4">
        <p className="text-sm text-muted-foreground">Loading results…</p>
      </div>
    )
  }

  return (
    <div className="min-h-screen flex flex-col bg-background text-foreground">
      <main className="flex-1">
        <div className="max-w-4xl mx-auto px-4 py-10 md:py-12">
          <ResultsContent />
        </div>
      </main>
    </div>
  )
}
