'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { useParams } from 'next/navigation';
import { QuizResult, QuestionResult, ApiError } from '@/types';
import { getQuizResult } from '@/lib/api';
import { Header } from '@/components/Header';
import { LoadingOverlay } from '@/components/LoadingOverlay';
import { ErrorBanner } from '@/components/ErrorBanner';
import { useQuizStore } from '@/lib/quizStore';

export default function ResultsPage() {
  const router = useRouter();
  const params = useParams();
  const sessionId = params.sessionId as string;
  const { clearSession } = useQuizStore();

  const [result, setResult] = useState<QuizResult | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<ApiError | null>(null);
  const [expandedQuestion, setExpandedQuestion] = useState<string | null>(null);

  useEffect(() => {
    const fetchResult = async () => {
      try {
        setIsLoading(true);
        setError(null);
        const data = await getQuizResult(sessionId);
        setResult(data);
      } catch (err) {
        setError({
          message: 'Failed to load results. Please try again.',
        });
      } finally {
        setIsLoading(false);
      }
    };

    fetchResult();
  }, [sessionId]);

  if (isLoading) {
    return (
      <div>
        <Header />
        <div className="container py-12">
          <div className="flex justify-center">
            <div className="text-center">
              <div className="flex gap-1 justify-center mb-4">
                <div className="w-2 h-8 bg-accent rounded-full animate-bounce" style={{ animationDelay: '0ms' }} />
                <div className="w-2 h-8 bg-accent rounded-full animate-bounce" style={{ animationDelay: '150ms' }} />
                <div className="w-2 h-8 bg-accent rounded-full animate-bounce" style={{ animationDelay: '300ms' }} />
              </div>
              <p className="text-muted-foreground">Loading your results...</p>
            </div>
          </div>
        </div>
      </div>
    );
  }

  if (error || !result) {
    return (
      <div>
        <Header />
        <div className="container py-12">
          {error && (
            <ErrorBanner
              message={error.message}
              onRetry={() => window.location.reload()}
            />
          )}
        </div>
      </div>
    );
  }

  const scoreColor = result.percentage >= 70 ? 'text-success' : result.percentage >= 50 ? 'text-accent' : 'text-destructive';
  const performanceText =
    result.percentage >= 80
      ? 'Excellent work!'
      : result.percentage >= 60
        ? 'Good job!'
        : 'Keep practicing!';

  const handleRetakeQuiz = () => {
    clearSession();
    router.push('/');
  };

  return (
    <div>
      <Header />

      <div className="container py-12 md:py-16">
        <div className="max-w-2xl mx-auto">
          {/* Score Section */}
          <div className="text-center mb-12">
            <h1 className="text-4xl md:text-5xl font-serif font-bold mb-4">Quiz Complete!</h1>
            <div className={`text-7xl font-serif font-bold ${scoreColor} mb-4`}>
              {result.percentage.toFixed(0)}%
            </div>
            <h2 className="text-2xl font-serif mb-2">{performanceText}</h2>
            <p className="text-lg text-muted-foreground">
              You answered {result.score} out of {result.totalQuestions} questions correctly.
            </p>
          </div>

          {/* Quiz Info */}
          <div className="bg-muted p-6 rounded-[var(--radius)] mb-8">
            <div className="grid grid-cols-2 md:grid-cols-3 gap-4 text-center">
              <div>
                <p className="text-sm text-muted-foreground mb-1">Category</p>
                <p className="font-serif font-bold text-lg">{result.categoryName}</p>
              </div>
              <div>
                <p className="text-sm text-muted-foreground mb-1">Difficulty</p>
                <p className="font-serif font-bold text-lg capitalize">{result.difficulty}</p>
              </div>
              <div className="col-span-2 md:col-span-1">
                <p className="text-sm text-muted-foreground mb-1">Score</p>
                <p className="font-serif font-bold text-lg">
                  {result.score}/{result.totalQuestions}
                </p>
              </div>
            </div>
          </div>

          {/* Review Section */}
          <div className="mb-12">
            <h3 className="text-2xl font-serif font-bold mb-6">Review Your Answers</h3>
            <div className="space-y-4">
              {result.questions.map((question) => (
                <button
                  key={question.id}
                  onClick={() =>
                    setExpandedQuestion(expandedQuestion === question.id ? null : question.id)
                  }
                  className="w-full text-left p-4 rounded-[var(--radius)] border border-border hover:border-accent transition-colors"
                >
                  <div className="flex items-start gap-3">
                    <div
                      className={`flex-shrink-0 w-6 h-6 rounded-full flex items-center justify-center font-bold text-sm mt-1 ${
                        question.isCorrect
                          ? 'bg-success text-success-foreground'
                          : 'bg-destructive text-destructive-foreground'
                      }`}
                    >
                      {question.isCorrect ? '✓' : '✗'}
                    </div>
                    <div className="flex-1 min-w-0">
                      <p className="font-serif font-bold text-lg text-balance">{question.text}</p>
                      <p className="text-sm text-muted-foreground mt-1">
                        {question.isCorrect ? 'Correct' : 'Incorrect'}
                      </p>
                    </div>
                    <div className="flex-shrink-0 text-accent">
                      {expandedQuestion === question.id ? '−' : '+'}
                    </div>
                  </div>

                  {/* Expanded Details */}
                  {expandedQuestion === question.id && (
                    <div className="mt-4 pt-4 border-t border-border space-y-3">
                      <div>
                        <p className="text-sm font-medium text-muted-foreground mb-2">Your answer:</p>
                        <p className="p-3 bg-muted rounded-[var(--radius)]">
                          {question.userAnswer !== null
                            ? question.options[question.userAnswer]
                            : 'Not answered'}
                        </p>
                      </div>
                      <div>
                        <p className="text-sm font-medium text-muted-foreground mb-2">Correct answer:</p>
                        <p className="p-3 bg-success bg-opacity-10 text-success rounded-[var(--radius)]">
                          {question.options[question.correctAnswer]}
                        </p>
                      </div>
                      {question.explanation && (
                        <div>
                          <p className="text-sm font-medium text-muted-foreground mb-2">Explanation:</p>
                          <p className="text-base leading-relaxed">{question.explanation}</p>
                        </div>
                      )}
                    </div>
                  )}
                </button>
              ))}
            </div>
          </div>

          {/* CTA Buttons */}
          <div className="flex gap-4">
            <button onClick={handleRetakeQuiz} className="button-primary flex-1">
              Take Another Quiz
            </button>
            <button onClick={() => router.push('/')} className="button-secondary flex-1">
              Home
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
