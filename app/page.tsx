'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { Category, ApiError } from '@/types';
import { getCategories, generateQuiz } from '@/lib/api';
import { useQuizStore } from '@/lib/quizStore';
import { Header } from '@/components/Header';
import { LoadingOverlay } from '@/components/LoadingOverlay';
import { ErrorBanner } from '@/components/ErrorBanner';

const DIFFICULTIES = ['Easy', 'Medium', 'Hard'];

export default function HomePage() {
  const router = useRouter();
  const { setCurrentSession, setLoading, setError, error } = useQuizStore();

  const [categories, setCategories] = useState<Category[]>([]);
  const [selectedCategory, setSelectedCategory] = useState<string | null>(null);
  const [selectedDifficulty, setSelectedDifficulty] = useState<string | null>(null);
  const [isLoadingCategories, setIsLoadingCategories] = useState(true);
  const [isGeneratingQuiz, setIsGeneratingQuiz] = useState(false);

  useEffect(() => {
    const fetchCategories = async () => {
      try {
        setIsLoadingCategories(true);
        setError(null);
        const data = await getCategories();
        setCategories(data);
      } catch (err) {
        const apiError: ApiError = {
          message: 'Failed to load categories. Please try again.',
        };
        setError(apiError);
      } finally {
        setIsLoadingCategories(false);
      }
    };

    fetchCategories();
  }, [setError]);

  const handleStartQuiz = async () => {
    if (!selectedCategory || !selectedDifficulty) return;

    try {
      setIsGeneratingQuiz(true);
      setLoading(true, 'Generating your quiz...');
      setError(null);

      const session = await generateQuiz(
        selectedCategory,
        selectedDifficulty.toLowerCase()
      );

      setCurrentSession(session);
      setLoading(false);
      router.push('/quiz');
    } catch (err) {
      const apiError: ApiError = {
        message: 'Failed to generate quiz. Please try again.',
      };
      setError(apiError);
      setLoading(false);
    } finally {
      setIsGeneratingQuiz(false);
    }
  };

  const canStartQuiz = selectedCategory && selectedDifficulty && !isGeneratingQuiz;

  return (
    <div>
      <Header />
      <LoadingOverlay isVisible={isGeneratingQuiz} />

      <div className="container py-12 md:py-16">
        <div className="max-w-2xl mx-auto">
          {/* Hero Section */}
          <div className="mb-12 text-center">
            <h1 className="mb-4 text-balance">Test Your News Knowledge</h1>
            <p className="text-xl text-muted-foreground">
              Choose a category and difficulty level to start your daily brief quiz.
            </p>
          </div>

          {/* Error Banner */}
          {error && (
            <div className="mb-6">
              <ErrorBanner
                message={error.message}
                onRetry={() => window.location.reload()}
                onDismiss={() => setError(null)}
              />
            </div>
          )}

          {/* Loading State */}
          {isLoadingCategories ? (
            <div className="flex justify-center py-12">
              <div className="text-center">
                <div className="flex gap-1 justify-center mb-4">
                  <div className="w-2 h-8 bg-accent rounded-full animate-bounce" style={{ animationDelay: '0ms' }} />
                  <div className="w-2 h-8 bg-accent rounded-full animate-bounce" style={{ animationDelay: '150ms' }} />
                  <div className="w-2 h-8 bg-accent rounded-full animate-bounce" style={{ animationDelay: '300ms' }} />
                </div>
                <p className="text-muted-foreground">Loading categories...</p>
              </div>
            </div>
          ) : (
            <>
              {/* Categories Section */}
              <div className="mb-12">
                <h2 className="text-2xl font-serif font-bold mb-6">Choose a Category</h2>
                <div className="grid gap-3">
                  {categories.map((category) => (
                    <button
                      key={category.id}
                      onClick={() => setSelectedCategory(category.id)}
                      className={`p-4 text-left rounded-[var(--radius)] border-2 transition-all ${
                        selectedCategory === category.id
                          ? 'border-accent bg-accent bg-opacity-5'
                          : 'border-border hover:border-accent'
                      }`}
                      aria-pressed={selectedCategory === category.id}
                    >
                      <div className="font-serif font-bold text-lg">{category.name}</div>
                      <div className="text-sm text-muted-foreground mt-1">{category.description}</div>
                    </button>
                  ))}
                </div>
              </div>

              {/* Difficulty Section */}
              <div className="mb-12">
                <h2 className="text-2xl font-serif font-bold mb-6">Select Difficulty</h2>
                <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                  {DIFFICULTIES.map((difficulty) => (
                    <button
                      key={difficulty}
                      onClick={() => setSelectedDifficulty(difficulty)}
                      className={`p-6 rounded-[var(--radius)] border-2 transition-all text-center min-h-[120px] flex items-center justify-center ${
                        selectedDifficulty === difficulty
                          ? 'border-accent bg-accent bg-opacity-5'
                          : 'border-border hover:border-accent'
                      }`}
                      aria-pressed={selectedDifficulty === difficulty}
                    >
                      <div className="font-serif font-bold text-2xl">{difficulty}</div>
                    </button>
                  ))}
                </div>
              </div>

              {/* Start Button */}
              <button
                onClick={handleStartQuiz}
                disabled={!canStartQuiz}
                className="button-primary w-full text-lg"
              >
                Start Quiz
              </button>
            </>
          )}
        </div>
      </div>
    </div>
  );
}
