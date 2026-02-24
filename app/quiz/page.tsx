'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { useQuizStore } from '@/lib/quizStore';
import { submitQuiz } from '@/lib/api';
import { Header } from '@/components/Header';
import { LoadingOverlay } from '@/components/LoadingOverlay';
import { ErrorBanner } from '@/components/ErrorBanner';

export default function QuizPage() {
  const router = useRouter();
  const { currentSession, recordAnswer, setLoading, setError, error } = useQuizStore();
  const [isSubmitting, setIsSubmitting] = useState(false);

  // Redirect if no quiz session
  useEffect(() => {
    if (!currentSession) {
      router.push('/');
    }
  }, [currentSession, router]);

  if (!currentSession) {
    return null;
  }

  const currentQuestion = currentSession.questions[currentSession.currentQuestionIndex];
  const progressPercentage = ((currentSession.currentQuestionIndex + 1) / currentSession.questions.length) * 100;
  const userAnswer = currentSession.answers[currentSession.currentQuestionIndex];

  const handleAnswerClick = (optionIndex: number) => {
    recordAnswer(currentSession.currentQuestionIndex, optionIndex);
  };

  const handleNext = async () => {
    const isLastQuestion = currentSession.currentQuestionIndex === currentSession.questions.length - 1;

    if (isLastQuestion) {
      // Submit the quiz
      try {
        setIsSubmitting(true);
        setLoading(true, 'Calculating your score...');
        setError(null);

        const result = await submitQuiz(currentSession);

        setLoading(false);
        router.push(`/results/${result.sessionId}`);
      } catch (err) {
        setError({
          message: 'Failed to submit quiz. Please try again.',
        });
        setLoading(false);
      } finally {
        setIsSubmitting(false);
      }
    } else {
      // Move to next question
      useQuizStore.setState({
        currentSession: {
          ...currentSession,
          currentQuestionIndex: currentSession.currentQuestionIndex + 1,
        },
      });
    }
  };

  const handlePrevious = () => {
    if (currentSession.currentQuestionIndex > 0) {
      useQuizStore.setState({
        currentSession: {
          ...currentSession,
          currentQuestionIndex: currentSession.currentQuestionIndex - 1,
        },
      });
    }
  };

  const isLastQuestion = currentSession.currentQuestionIndex === currentSession.questions.length - 1;
  const canProceed = userAnswer !== null;

  return (
    <div>
      <Header showBackButton backHref="/" />
      <LoadingOverlay isVisible={isSubmitting} />

      <div className="container py-8 md:py-12">
        {/* Progress Bar */}
        <div className="mb-8">
          <div className="flex justify-between items-center mb-2">
            <span className="text-sm text-muted-foreground font-medium">
              Question {currentSession.currentQuestionIndex + 1} of {currentSession.questions.length}
            </span>
            <span className="text-sm font-medium text-accent">{Math.round(progressPercentage)}%</span>
          </div>
          <div className="w-full bg-muted rounded-full h-2 overflow-hidden">
            <div
              className="bg-accent h-full rounded-full transition-all duration-300"
              style={{ width: `${progressPercentage}%` }}
              role="progressbar"
              aria-valuenow={currentSession.currentQuestionIndex + 1}
              aria-valuemin={1}
              aria-valuemax={currentSession.questions.length}
            />
          </div>
        </div>

        {/* Error Banner */}
        {error && (
          <div className="mb-6">
            <ErrorBanner
              message={error.message}
              onDismiss={() => setError(null)}
            />
          </div>
        )}

        {/* Question Section */}
        <div className="max-w-2xl mx-auto mb-12">
          <h1 className="text-3xl md:text-4xl font-serif font-bold mb-8 text-balance">
            {currentQuestion.text}
          </h1>

          {/* Options */}
          <div className="space-y-3">
            {currentQuestion.options.map((option, index) => (
              <button
                key={index}
                onClick={() => handleAnswerClick(index)}
                className={`w-full p-4 text-left rounded-[var(--radius)] border-2 transition-all text-lg ${
                  userAnswer === index
                    ? 'border-accent bg-accent bg-opacity-5'
                    : 'border-border hover:border-accent'
                }`}
                aria-pressed={userAnswer === index}
                disabled={isSubmitting}
              >
                <div className="flex items-start gap-3">
                  <div
                    className={`flex-shrink-0 w-6 h-6 rounded-full border-2 flex items-center justify-center font-medium text-sm mt-0.5 transition-all ${
                      userAnswer === index
                        ? 'border-accent bg-accent text-accent-foreground'
                        : 'border-border'
                    }`}
                  >
                    {userAnswer === index && '✓'}
                  </div>
                  <span>{option}</span>
                </div>
              </button>
            ))}
          </div>
        </div>

        {/* Navigation Buttons */}
        <div className="max-w-2xl mx-auto flex gap-4">
          <button
            onClick={handlePrevious}
            disabled={currentSession.currentQuestionIndex === 0 || isSubmitting}
            className="button-secondary flex-1"
          >
            Previous
          </button>
          <button
            onClick={handleNext}
            disabled={!canProceed || isSubmitting}
            className="button-primary flex-1"
          >
            {isLastQuestion ? 'Submit Quiz' : 'Next Question'}
          </button>
        </div>
      </div>
    </div>
  );
}
