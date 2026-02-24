"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { getCategories, generateQuiz } from "@/lib/api";
import { useQuizStore } from "@/lib/quizStore";
import { Category, QuizSession } from "@/types";
import { Header } from "@/components/Header";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Separator } from "@/components/ui/separator";
import { AlertCircle, X, Loader2 } from "lucide-react";

const DIFFICULTIES = ["Easy", "Medium", "Hard"];

export default function HomePage() {
  const router = useRouter();
  const { setCurrentSession, setLoading, setError, error } = useQuizStore();

  const [categories, setCategories] = useState<Category[]>([]);
  const [selectedCategories, setSelectedCategories] = useState<string[]>([]);
  const [selectedDifficulty, setSelectedDifficulty] = useState<string | null>(
    null,
  );
  const [isLoadingCategories, setIsLoadingCategories] = useState(true);
  const [isGeneratingQuiz, setIsGeneratingQuiz] = useState(false);

  useEffect(() => {
    const fetchCategories = async () => {
      try {
        setIsLoadingCategories(true);
        setError(null);
        const data = await getCategories();
        setCategories(data);
      } catch {
        setError({ message: "Failed to load categories. Please try again." });
      } finally {
        setIsLoadingCategories(false);
      }
    };

    fetchCategories();
  }, [setError]);

  const toggleCategory = (id: string) => {
    setSelectedCategories((prev) =>
      prev.includes(id) ? prev.filter((c) => c !== id) : [...prev, id],
    );
  };

  const handleStartQuiz = async () => {
    if (!selectedCategories.length || !selectedDifficulty) return;

    try {
      setIsGeneratingQuiz(true);
      setLoading(true, "Generating your quiz...");
      setError(null);
      const rawQuiz: any = await generateQuiz(
        selectedCategories,
        selectedDifficulty.toLowerCase(),
      );

      const transformedQuestions = rawQuiz.questions.map(
        (q: any, index: number) => ({
          id: String(index),
          text: q.question,
          options: q.options,
          correctAnswer: q.answer,
          explanation: "",
        }),
      );

      const session: QuizSession = {
        id: String(rawQuiz.id ?? ""),
        categoryId: selectedCategories.join(","),
        difficulty: rawQuiz.difficulty ?? selectedDifficulty.toLowerCase(),
        questions: transformedQuestions,
        currentQuestionIndex: 0,
        answers: Array(transformedQuestions.length).fill(null),
        score: 0,
        completed: false,
        createdAt:
          rawQuiz.created_at ?? new Date().toISOString(),
      };

      setCurrentSession(session);
      setLoading(false);
      router.push("/quiz");
    } catch {
      setError({ message: "Failed to generate quiz. Please try again." });
      setLoading(false);
    } finally {
      setIsGeneratingQuiz(false);
    }
  };

  const canStart =
    selectedCategories.length > 0 && selectedDifficulty && !isGeneratingQuiz;

  return (
    <div className="min-h-screen bg-background">
      <Header />

      <main className="max-w-3xl mx-auto px-6 py-16">
        <div className="mb-10">
          <h1 className="text-3xl font-semibold tracking-tight mb-2">
            Test your knowledge
          </h1>
          <p className="text-muted-foreground">
            Pick a topic and difficulty to begin.
          </p>
        </div>

        <Separator className="mb-10" />

        {error && (
          <Alert variant="destructive" className="mb-8">
            <AlertCircle className="h-4 w-4" />
            <AlertDescription className="flex items-center justify-between">
              <span>{error.message}</span>
              <Button
                variant="ghost"
                size="icon"
                className="h-5 w-5 ml-4 shrink-0"
                onClick={() => setError(null)}
              >
                <X className="h-3 w-3" />
              </Button>
            </AlertDescription>
          </Alert>
        )}

        {isLoadingCategories ? (
          <div className="flex items-center gap-2 text-sm text-muted-foreground py-8">
            <Loader2 className="h-4 w-4 animate-spin" />
            <span>Loading categories…</span>
          </div>
        ) : (
          <>
            {/* Categories */}
            <section className="mb-12">
              <p className="text-xs font-medium tracking-widest uppercase text-muted-foreground mb-4">
                Category
              </p>

              <div className="flex flex-wrap gap-3">
                {categories.map((cat) => {
                  const id = String((cat as any).ID ?? (cat as any).id);
                  const name = (cat as any).Name ?? (cat as any).name;
                  const isSelected = selectedCategories.includes(id);

                  return (
                    <Card
                      key={id}
                      onClick={() => toggleCategory(id)}
                      className={`cursor-pointer transition-all hover:shadow-sm px-4 py-2 ${
                        isSelected ? "ring-2 ring-primary" : ""
                      }`}
                    >
                      <CardContent className="p-0">
                        <p className="text-sm font-medium whitespace-nowrap">
                          {name}
                        </p>
                      </CardContent>
                    </Card>
                  );
                })}
              </div>
            </section>

            {/* Difficulty */}
            <section className="mb-12">
              <p className="text-xs font-medium tracking-widest uppercase text-muted-foreground mb-4">
                Difficulty
              </p>

              <div className="flex gap-3 flex-wrap">
                {DIFFICULTIES.map((d) => (
                  <Button
                    key={d}
                    variant={selectedDifficulty === d ? "default" : "outline"}
                    size="sm"
                    onClick={() => setSelectedDifficulty(d)}
                  >
                    {d}
                  </Button>
                ))}
              </div>
            </section>

            <Button
              onClick={handleStartQuiz}
              disabled={!canStart}
              size="lg"
              className="w-full sm:w-auto"
            >
              {isGeneratingQuiz ? (
                <>
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                  Generating…
                </>
              ) : (
                "Start quiz"
              )}
            </Button>
          </>
        )}
      </main>
    </div>
  );
}
