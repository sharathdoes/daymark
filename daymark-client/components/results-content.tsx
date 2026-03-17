'use client'

import { useRouter } from 'next/navigation'
import { useEffect, useRef, useState } from 'react'
import { useQuiz, useAuth } from '@/lib/store'
import { saveQuizResult } from '@/lib/api'

import { Button } from '@/components/ui/button'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { LinkIcon, Check } from 'lucide-react'

const LETTER_LABELS = ['A', 'B', 'C', 'D']

function getLabel(percentage: number): string {
  if (percentage === 100) return 'Perfect Score'
  if (percentage >= 67) return 'Great Work'
  if (percentage >= 34) return 'Good Effort'
  return 'Keep Reading'
}

export default function ResultsContent() {
  const router = useRouter()
  const { session, resetSession } = useQuiz()
  const { isAuthenticated } = useAuth()
  const [isSaving, setIsSaving] = useState(false)
  const savedRef = useRef(false)

  const score = session?.score ?? 0
  const total = session?.quiz?.questions?.length ?? 0
  const percentage = total > 0 ? Math.round((score / total) * 100) : 0
  const label = getLabel(percentage)

  useEffect(() => {
    if (session?.completed && isAuthenticated && !savedRef.current) {
      savedRef.current = true
      setIsSaving(true)
      
      const categoryIdsStr = session.quiz.category_ids.join(',')

      saveQuizResult({
        quiz_id: session.quiz.id,
        score: session.score,
        total_questions: session.quiz.questions.length,
        difficulty: session.quiz.difficulty,
        categories: categoryIdsStr,
      })
      .catch(err => console.error("Failed to save quiz result:", err))
      .finally(() => setIsSaving(false))
    }
  }, [session, isAuthenticated])
  
  const [copied, setCopied] = useState(false)

  if (!session || !session.completed) return null

  const handleTryAgain = () => {
    resetSession()
    const categoryIds = session.quiz.category_ids.join(',')
    const difficulty = session.quiz.difficulty
    router.push(`/?categories=${categoryIds}&difficulty=${difficulty}`)
  }

  const handleNewQuiz = () => {
    resetSession()
    router.push('/')
  }

  const handleShare = () => {
    const url = `${window.location.origin}/quiz/${session.quiz.id}`
    navigator.clipboard.writeText(url)
    setCopied(true)
    setTimeout(() => setCopied(false), 2000)
  }

  return (
    <div className="space-y-10">
      {/* Score summary */}

      <div className="flex flex-col sm:flex-row gap-3 justify-center">
        <Button
          type="button"
          onClick={handleTryAgain}
          className="min-w-[150px]"
        >
          Try again with same settings
        </Button>
        <Button
          type="button"
          variant="secondary"
          onClick={handleShare}
          className="min-w-[150px]"
        >
          {copied ? (
            <Check className="w-4 h-4 mr-2" />
          ) : (
            <LinkIcon className="w-4 h-4 mr-2" />
          )}
          {copied ? "Copied!" : "Share Quiz"}
        </Button>
        <Button
          type="button"
          variant="outline"
          onClick={handleNewQuiz}
          className="min-w-[150px]"
        >
          Start a new quiz
        </Button>
      </div>
      
      <Card>
        <CardHeader className="text-center space-y-3">
          <CardTitle className="text-4xl md:text-5xl font-semibold">
            {score} / {total}
          </CardTitle>
          <CardDescription className="text-sm md:text-base">
            {label} — you answered {percentage}% of questions correctly.
          </CardDescription>
        </CardHeader>
      </Card>

      {/* Question review */}
      <section className="space-y-5">
        <div className="flex items-center justify-between">
          <h2 className="text-sm font-medium tracking-[0.18em] uppercase text-muted-foreground">
            Review
          </h2>
          <p className="text-xs text-muted-foreground">
            Correct answers are highlighted in green.
          </p>
        </div>

        <div className="space-y-5">
          {session.quiz.questions.map((question, qIndex) => {
            const userAnswer = session.answers[qIndex];

            return (
              <Card key={qIndex} className="border-border/80">
                <CardContent className="pt-5 space-y-4">
                  <div className="flex items-start gap-3">
                    <span className="mt-0.5 text-xs font-mono font-semibold text-muted-foreground w-5 flex-shrink-0">
                      {qIndex + 1}
                    </span>
                    <p className="text-sm text-foreground leading-relaxed">
                      {question.question}
                    </p>
                  </div>

                  <div className="space-y-1.5 pl-8">
                    {question.options.map((option, optIndex) => {
                      const isQuestionCorrect = optIndex === question.answer;
                      const isUserSelected = optIndex === userAnswer;
                      const isUserWrong = isUserSelected && !isQuestionCorrect;

                      return (
                        <div
                          key={optIndex}
                          className={`flex items-start gap-2 rounded-md border px-3 py-1.5 text-xs md:text-sm ${
                            isQuestionCorrect
                              ? "border-emerald-500 bg-emerald-500/20 text-emerald-800 dark:text-emerald-100"
                              : isUserWrong
                                ? "border-destructive/60 bg-destructive/10 text-destructive"
                                : "border-border/70 text-muted-foreground"
                          }`}
                        >
                          <span className="mt-0.5 font-mono font-semibold">
                            {LETTER_LABELS[optIndex]}
                          </span>
                          <span>{option}</span>
                        </div>
                      );
                    })}
                  </div>

                  <p className="pl-8 text-[11px] text-muted-foreground mt-4 inline-block bg-accent/20 px-3 py-1.5 rounded-md border border-border/50 font-medium">
                    View the source article for this answer:{" "}
                    <a
                      href={question.article_url}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="text-foreground underline underline-offset-4 font-semibold"
                    >
                      {new URL(question.article_url).hostname}
                    </a>
                  </p>
                </CardContent>
              </Card>
            );
          })}
        </div>
      </section>

      {/* Actions */}
    </div>
  );
}
