'use client'

import { useEffect, useState } from 'react'
import { useRouter, useSearchParams } from 'next/navigation'
import Link from 'next/link'
import { useAuth, useQuiz } from '@/lib/store'
import { getCategories, generateQuiz } from '@/lib/api'
import { Category, QuizSession } from '@/lib/types'
import Header from '@/components/header'
import LoadingOverlay from '@/components/loading-overlay'
import { Button } from '@/components/ui/button'
import { Highlighter } from '@/components/ui/highlighter'
import { AnimatedShinyText } from '@/components/ui/animated-shiny-text'
import { ArrowRightIcon } from '@radix-ui/react-icons'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'

const DIFFICULTIES = [
  { id: 'easy', label: 'Easy', description: 'Perfect for warming up' },
  { id: 'medium', label: 'Medium', description: 'Test your knowledge' },
  { id: 'hard', label: 'Hard', description: 'Challenge yourself' },
]

const QUESTION_COUNTS = [5, 10, 15]

const TIMER_OPTIONS = [
  { id: 'none', label: 'No timer' },
  { id: '5', label: '5 min' },
  { id: '10', label: '10 min' },
]

export default function HomePage() {
  const router = useRouter()
  const searchParams = useSearchParams()
  const { user, isAuthenticated } = useAuth()
  const { setSession } = useQuiz()

  const [categories, setCategories] = useState<Category[]>([])
  const [selectedCategoryIds, setSelectedCategoryIds] = useState<number[]>([])
  const [selectedDifficulty, setSelectedDifficulty] = useState<string | null>(null)
  const [questionCount, setQuestionCount] = useState<number>(10)
  const [timerOption, setTimerOption] = useState<string>('none')
  const [isLoading, setIsLoading] = useState(true)
  const [isGenerating, setIsGenerating] = useState(false)
  const [error, setError] = useState('')

  // Load categories on mount
  useEffect(() => {
    const loadCategories = async () => {
      try {
        const data = await getCategories()
        setCategories(data)
        
        // Check for query params from results page
        const categoriesParam = searchParams.get('categories')
        const difficultyParam = searchParams.get('difficulty')
        
        if (categoriesParam) {
          const ids = categoriesParam.split(',').map(Number)
          setSelectedCategoryIds(ids)
        }
        if (difficultyParam) {
          setSelectedDifficulty(difficultyParam)
        }
      } catch (err) {
        setError('Failed to load categories')
        console.error(err)
      } finally {
        setIsLoading(false)
      }
    }

    loadCategories()
  }, [searchParams])

  const toggleCategory = (id: number) => {
    setSelectedCategoryIds(prev =>
      prev.includes(id) ? prev.filter(cid => cid !== id) : [...prev, id]
    )
  }

  const canStart = selectedCategoryIds.length > 0 && selectedDifficulty !== null

  const handleStart = async () => {
    if (!canStart) {
      setError('Please select at least one category and a difficulty level')
      return
    }

    if (!isAuthenticated) {
      const params = new URLSearchParams()
      params.set('categories', selectedCategoryIds.join(','))
      if (selectedDifficulty) params.set('difficulty', selectedDifficulty)
      params.set('count', String(questionCount))
      params.set('timer', timerOption)

      // Persist in localStorage so OAuth redirects don't lose the selection
      localStorage.setItem('pendingQuiz', JSON.stringify({
        categories: selectedCategoryIds,
        difficulty: selectedDifficulty,
        count: questionCount,
        timer: timerOption,
      }))

      router.push(`/login?${params.toString()}`)
      return
    }

    setIsGenerating(true)
    setError('')

    try {
      const quiz = await generateQuiz({
        category_ids: selectedCategoryIds,
        difficulty: selectedDifficulty!,
        number_of_questions: questionCount,
      })

      const session: QuizSession = {
        quiz,
        currentIndex: 0,
        answers: Array(quiz.questions.length).fill(null),
        score: 0,
        completed: false,
        timerOption: timerOption,
        timerSeconds: timerOption === 'none' ? 0 : parseInt(timerOption, 10) * 60,
      }

      setSession(session)
      router.push('/quiz')
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to generate quiz')
      console.error(err)
    } finally {
      setIsGenerating(false)
    }
  }

  return (
    <div className="min-h-screen flex flex-col bg-background text-foreground">
      {isGenerating && <LoadingOverlay />}

      <main className="flex-1">
        <div className="max-w-5xl mx-auto px-4 py-10 md:py-14">
          {/* Hero */}
          <section className="mb-10 md:mb-12 text-center">
            <Link href="/daily" className="mb-4 inline-block">
              <div className="group rounded-full border border-black/5 bg-neutral-100 text-base text-white transition-all ease-in hover:cursor-pointer hover:bg-neutral-200 dark:border-white/5 dark:bg-neutral-900 dark:hover:bg-neutral-800">
                <AnimatedShinyText className="inline-flex items-center justify-center px-4 py-1 transition ease-out hover:text-neutral-600 hover:duration-300 hover:dark:text-neutral-400">
                  <span>✨ Quiz of the Day</span>
                  <ArrowRightIcon className="ml-1 size-3 transition-transform duration-300 ease-in-out group-hover:translate-x-0.5" />
                </AnimatedShinyText>
              </div>
            </Link>

            <h1 className="text-3xl md:text-4xl lg:text-5xl font-semibold tracking-tight text-balance mb-3">
              You read the news{" "}
              <Highlighter action="underline" color="#FF9800" animationDuration={800}>
                today
              </Highlighter>
              ?{" "}<br />
              Let&apos;s see how much stuck in.
            </h1>
            <p className="text-sm md:text-base text-muted-foreground max-w-xl mx-auto">
              Every question is pulled from a real article published today. No
              trivia, no recycled facts — just{" "}
              <Highlighter action="highlight" color="#FF980033" animationDuration={1000}>
                today&apos;s news
              </Highlighter>
              {" "}and{" "}
              <Highlighter action="underline" color="#6366f1" animationDuration={1200}>
                how well you followed it
              </Highlighter>
              .
            </p>
          </section>

          {error && (
            <div className="mb-6">
              <Card className="bg-destructive/5 border-destructive/30 text-destructive">
                <CardContent className="py-3 text-sm">{error}</CardContent>
              </Card>
            </div>
          )}


          {isLoading ? (
            <div className="flex justify-center py-16">
              <p className="text-sm text-muted-foreground">
                Loading categories…
              </p>
            </div>
          ) : (
            <div className="grid gap-6 md:grid-cols-[minmax(0,1.4fr)_minmax(0,1fr)] items-start">
              {/* Categories */}
              <Card>
                <CardHeader>
                  <CardTitle>Select categories</CardTitle>
                  <CardDescription>
                    Choose the topics you&apos;d like today&apos;s quiz to focus
                    on.
                  </CardDescription>
                </CardHeader>
                <CardContent>
                  <div className="flex flex-wrap gap-2.5">
                    {categories.map((category) => {
                      const selected = selectedCategoryIds.includes(
                        category.id,
                      );
                      return (
                        <Button
                          key={category.id}
                          type="button"
                          variant={selected ? "default" : "outline"}
                          size="sm"
                          className="justify-start truncate"
                          onClick={() => toggleCategory(category.id)}
                        >
                          {category.name}
                        </Button>
                      );
                    })}
                  </div>

                  <div className=" flex gap-4 mt-6 space-y-4">
                    <div>
                      <p className="mb-1.5 text-xs font-medium text-muted-foreground">
                        Number of questions
                      </p>
                      <div className="flex flex-wrap gap-2">
                        {QUESTION_COUNTS.map((count) => (
                          <Button
                            key={count}
                            type="button"
                            size="sm"
                            variant={
                              questionCount === count ? "default" : "outline"
                            }
                            onClick={() => setQuestionCount(count)}
                          >
                            {count}
                          </Button>
                        ))}
                      </div>
                    </div>

                    <div>
                      <p className="mb-1.5 text-xs font-medium text-muted-foreground">
                        Timer (optional)
                      </p>
                      <div className="flex flex-wrap gap-2">
                        {TIMER_OPTIONS.map((option) => (
                          <Button
                            key={option.id}
                            type="button"
                            size="sm"
                            variant={
                              timerOption === option.id ? "default" : "outline"
                            }
                            onClick={() => setTimerOption(option.id)}
                          >
                            {option.label}
                          </Button>
                        ))}
                      </div>
                    </div>
                  </div>
                </CardContent>
              </Card>

              {/* Difficulty */}
              <Card className="h-full">
                <CardHeader>
                  <CardTitle>Choose difficulty</CardTitle>
                  <CardDescription>
                    Adjust how challenging the questions should be.
                  </CardDescription>
                </CardHeader>
                <CardContent className="space-y-3">
                  {DIFFICULTIES.map((difficulty) => {
                    const active = selectedDifficulty === difficulty.id;
                    return (
                      <button
                        key={difficulty.id}
                        type="button"
                        onClick={() => setSelectedDifficulty(difficulty.id)}
                        className={`w-full rounded-md border px-4 py-3 text-left text-sm transition-colors ${
                          active
                            ? "border-primary bg-primary/5"
                            : "border-border bg-background hover:bg-accent/40"
                        }`}
                      >
                        <div className="font-medium mb-0.5">
                          {difficulty.label}
                        </div>
                        <p className="text-xs text-muted-foreground">
                          {difficulty.description}
                        </p>
                      </button>
                    );
                  })}
                </CardContent>
              </Card>
            </div>
          )}

          {/* Primary CTA */}
          <section className="mt-10 flex flex-col items-center gap-4 text-center">
            <div className="space-y-1">
              <p className="text-sm md:text-base text-muted-foreground max-w-md">
                {selectedCategoryIds.length === 0
                  ? "Choose a few topics and a difficulty to begin."
                  : `You\'re set for a ${selectedDifficulty ?? "news"} quiz across ${selectedCategoryIds.length} topic${selectedCategoryIds.length > 1 ? "s" : ""}.`}
              </p>
            </div>

            <Button
              type="button"
              size="lg"
              onClick={handleStart}
              disabled={!canStart || isGenerating}
              className="min-w-[200px] shadow-sm"
            >
              {isGenerating ? "Preparing your quiz…" : "Start today's quiz"}
            </Button>

            {!isAuthenticated && (
              <p className="text-xs md:text-sm text-muted-foreground">
                <Link href="/login" className="underline underline-offset-4">
                  Sign in
                </Link>{" "}
                to save your streaks and past quizzes.
              </p>
            )}
          </section>
        </div>
      </main>
    </div>
  );
}
