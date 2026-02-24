"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { QuizResult, ApiError } from "@/types";
import { useQuizStore } from "@/lib/quizStore";
import { Header } from "@/components/Header";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Separator } from "@/components/ui/separator";
import { Badge } from "@/components/ui/badge";
import {
  AlertCircle,
  Loader2,
  ChevronDown,
  ChevronUp,
  CheckCircle2,
  XCircle,
} from "lucide-react";

export default function ResultsPage() {
  const router = useRouter();
  const { lastResult, clearSession } = useQuizStore();

  const [result, setResult] = useState<QuizResult | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<ApiError | null>(null);
  const [expandedQuestion, setExpandedQuestion] = useState<string | null>(null);

  useEffect(() => {
    setIsLoading(true);
    if (!lastResult) {
      setError({ message: "No results found. Please take a quiz first." });
      setIsLoading(false);
      return;
    }
    setResult(lastResult);
    setIsLoading(false);
  }, [lastResult]);

  const handleRetakeQuiz = () => {
    clearSession();
    router.push("/");
  };

  if (isLoading) {
    return (
      <div className="min-h-screen bg-background">
        <Header />
        <div className="flex items-center gap-2 text-sm text-muted-foreground max-w-2xl mx-auto px-6 py-16">
          <Loader2 className="h-4 w-4 animate-spin" />
          <span>Loading results…</span>
        </div>
      </div>
    );
  }

  if (error || !result) {
    return (
      <div className="min-h-screen bg-background">
        <Header />
        <div className="max-w-2xl mx-auto px-6 py-16">
          <Alert variant="destructive">
            <AlertCircle className="h-4 w-4" />
            <AlertDescription className="flex items-center justify-between">
              <span>{error?.message ?? "Something went wrong."}</span>
              <Button
                variant="ghost"
                size="sm"
                onClick={() => router.push("/")}
              >
                Go home
              </Button>
            </AlertDescription>
          </Alert>
        </div>
      </div>
    );
  }

  const performanceText =
    result.percentage >= 80
      ? "Excellent work!"
      : result.percentage >= 60
        ? "Good job!"
        : "Keep practicing!";

  const scoreVariant =
    result.percentage >= 70
      ? "default"
      : result.percentage >= 50
        ? "secondary"
        : "destructive";

  return (
    <div className="min-h-screen bg-background">
      <Header />

      <main className="max-w-2xl mx-auto px-6 py-16">
        {/* Score */}
        <div className="mb-10">
          <p className="text-xs font-medium tracking-widest uppercase text-muted-foreground mb-3">
            Result
          </p>
          <div className="flex items-end gap-4 mb-2">
            <span className="text-6xl font-semibold tracking-tight">
              {result.percentage.toFixed(0)}%
            </span>
            <span className="text-muted-foreground text-lg mb-2">
              {performanceText}
            </span>
          </div>
          <p className="text-sm text-muted-foreground">
            {result.score} of {result.totalQuestions} correct
          </p>
        </div>

        {/* Meta */}
        <Card className="mb-10">
          <CardContent className="px-4 py-4 flex gap-6 flex-wrap">
            <div>
              <p className="text-xs text-muted-foreground mb-1">Category</p>
              <p className="text-sm font-medium">{result.categoryName}</p>
            </div>
            <div>
              <p className="text-xs text-muted-foreground mb-1">Difficulty</p>
              <p className="text-sm font-medium capitalize">
                {result.difficulty}
              </p>
            </div>
            <div>
              <p className="text-xs text-muted-foreground mb-1">Score</p>
              <Badge variant={scoreVariant} className="text-xs">
                {result.score}/{result.totalQuestions}
              </Badge>
            </div>
          </CardContent>
        </Card>

        <Separator className="mb-10" />

        {/* Review */}
        <section className="mb-12">
          <p className="text-xs font-medium tracking-widest uppercase text-muted-foreground mb-4">
            Review
          </p>
          <div className="flex flex-col gap-2">
            {result.questions.map((question) => {
              const isExpanded = expandedQuestion === question.id;
              return (
                <Card
                  key={question.id}
                  className="cursor-pointer hover:border-foreground/40 transition-colors"
                  onClick={() =>
                    setExpandedQuestion(isExpanded ? null : question.id)
                  }
                >
                  <CardContent className="px-4 py-3">
                    <div className="flex items-start gap-3">
                      {question.isCorrect ? (
                        <CheckCircle2 className="h-4 w-4 text-green-500 shrink-0 mt-0.5" />
                      ) : (
                        <XCircle className="h-4 w-4 text-destructive shrink-0 mt-0.5" />
                      )}
                      <span className="text-sm flex-1">{question.text}</span>
                      {isExpanded ? (
                        <ChevronUp className="h-4 w-4 text-muted-foreground shrink-0" />
                      ) : (
                        <ChevronDown className="h-4 w-4 text-muted-foreground shrink-0" />
                      )}
                    </div>

                    {isExpanded && (
                      <div className="mt-4 pt-4 border-t border-border space-y-3 text-sm">
                        <div>
                          <p className="text-xs text-muted-foreground mb-1">
                            Your answer
                          </p>
                          <p className="text-foreground">
                            {question.userAnswer !== null
                              ? question.options[question.userAnswer]
                              : "Not answered"}
                          </p>
                        </div>
                        <div>
                          <p className="text-xs text-muted-foreground mb-1">
                            Correct answer
                          </p>
                          <p className="text-green-600 dark:text-green-400">
                            {question.options[question.correctAnswer]}
                          </p>
                        </div>
                        {question.explanation && (
                          <div>
                            <p className="text-xs text-muted-foreground mb-1">
                              Explanation
                            </p>
                            <p className="text-foreground leading-relaxed">
                              {question.explanation}
                            </p>
                          </div>
                        )}
                      </div>
                    )}
                  </CardContent>
                </Card>
              );
            })}
          </div>
        </section>

        {/* Actions */}
        <div className="flex gap-3">
          <Button onClick={handleRetakeQuiz} size="lg">
            Take another quiz
          </Button>
          <Button variant="outline" size="lg" onClick={() => router.push("/")}>
            Home
          </Button>
        </div>
      </main>
    </div>
  );
}
