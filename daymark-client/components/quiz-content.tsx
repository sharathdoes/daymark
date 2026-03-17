'use client'

import React, { useEffect } from 'react'
import { useRouter } from 'next/navigation'
import { useQuiz } from '@/lib/store'
import { Question } from '@/lib/types'
import { Progress } from '@/components/ui/progress'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { saveQuizResult } from '@/lib/api'
import { Clock } from 'lucide-react'

const LETTER_LABELS = ['A', 'B', 'C', 'D']

// Formatting helper for mm:ss
function formatTime(seconds: number) {
  const m = Math.floor(seconds / 60)
  const s = seconds % 60
  return `${m}:${s < 10 ? '0' : ''}${s}`
}

export default function QuizContent() {
  const router = useRouter()
  const { session, answerQuestion, nextQuestion, decrementTimer } = useQuiz()

  // Timer interval
  useEffect(() => {
    if (!session || session.timerOption === 'none' || session.completed) return

    if (session.timerSeconds <= 0) {
      // Time is up, force navigate to results and mark complete
      router.push('/results')
      return
    }

    const interval = setInterval(() => {
      decrementTimer()
    }, 1000)

    return () => clearInterval(interval)
  }, [session?.timerSeconds, session?.timerOption, session?.completed, decrementTimer, router])
  
  if (!session) return null

  const totalQuestions = session.quiz.questions?.length ?? 0

  // Gracefully handle quizzes that have no questions (e.g. malformed legacy data)
  if (totalQuestions === 0) {
    return (
      <div className="space-y-6">
        <Card>
          <CardHeader>
            <CardTitle className="text-base md:text-lg leading-snug">
              This quiz has no questions
            </CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-sm text-muted-foreground">
              Something went wrong loading this quiz&apos;s questions. Please try generating a new quiz.
            </p>
          </CardContent>
        </Card>
        <div className="flex justify-center">
          <Button type="button" onClick={() => router.push('/')}>Go back home</Button>
        </div>
      </div>
    )
  }

  const currentQuestion = session.quiz.questions[session.currentIndex]
  const selectedAnswer = session.answers[session.currentIndex]
  const isAnswered = selectedAnswer !== null
  const isLastQuestion = session.currentIndex === totalQuestions - 1
  const progress = ((session.currentIndex + 1) / totalQuestions) * 100

  const handleOptionClick = (optionIndex: number) => {
    // Allow reselecting a different option before clicking Next
    answerQuestion(session.currentIndex, optionIndex)
  }

  const handleNext = () => {
    if (isAnswered) {
      nextQuestion()
      if (isLastQuestion) {
        // Navigate to results on next call
        setTimeout(() => router.push('/results'), 0)
      }
    }
  }

  return (
    <div className="space-y-8">
      {/* Top Bar: Progress & Timer */}
      <div className="space-y-2">
        <div className="flex items-center justify-between text-[11px] font-medium tracking-[0.18em] uppercase text-muted-foreground">
          <span>
            Question {session.currentIndex + 1} of {session.quiz.questions.length}
          </span>
          {session.timerOption !== 'none' && (
            <span className="flex items-center gap-1.5 text-foreground">
              <Clock className="w-3.5 h-3.5" />
              {formatTime(session.timerSeconds)}
            </span>
          )}
        </div>
        <Progress value={progress} />
      </div>

      {/* Question & options */}
      <Card>
        <CardHeader>
          <CardTitle className="text-base md:text-lg leading-snug">
            {currentQuestion.question}
          </CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="space-y-2.5">
            {currentQuestion.options.map((option, index) => {
              const isSelected = index === selectedAnswer

              return (
                <button
                  key={index}
                  type="button"
                  onClick={() => handleOptionClick(index)}
                  className={`flex w-full items-start gap-3 rounded-md border px-3.5 py-3 text-left text-sm transition-colors cursor-pointer ${
                    isSelected
                      ? 'border-primary bg-primary/10 text-primary'
                      : 'border-border bg-background hover:bg-accent/40'
                  }`}
                >
                  <span className="mt-0.5 w-7 flex-shrink-0 text-xs font-mono font-semibold text-muted-foreground">
                    {LETTER_LABELS[index]}
                  </span>
                  <span className="flex-1 text-sm md:text-[0.94rem] leading-relaxed">
                    {option}
                  </span>
                </button>
              )
            })}
          </div>
        </CardContent>
      </Card>

      {/* Next button */}
      <div className="flex justify-center">
        <Button
          type="button"
          onClick={handleNext}
          disabled={!isAnswered}
          className="min-w-[180px]"
        >
          {isLastQuestion && isAnswered
            ? 'See results'
            : isAnswered
            ? 'Next question'
            : 'Select an answer'}
        </Button>
      </div>
    </div>
  )
}
