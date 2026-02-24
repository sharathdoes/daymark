"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { useQuizStore } from "@/lib/quizStore";
import { Header } from "@/components/Header";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Progress } from "@/components/ui/progress";
import { AlertCircle, X, Loader2 } from "lucide-react";

export default function QuizPage() {
  const router = useRouter();
  const { currentSession, recordAnswer, setLoading, setError, computeResult, error } =
    useQuizStore();
  const [isSubmitting, setIsSubmitting] = useState(false);

  useEffect(() => {
    if (!currentSession) router.push("/");
  }, [currentSession, router]);

  if (!currentSession) return null;

  const { questions, currentQuestionIndex, answers } = currentSession;
  const question = questions[currentQuestionIndex];
  const userAnswer = answers[currentQuestionIndex];
  const isLast = currentQuestionIndex === questions.length - 1;
  const progress = ((currentQuestionIndex + 1) / questions.length) * 100;

  const handleNext = async () => {
    if (isLast) {
      try {
        setIsSubmitting(true);
        setLoading(true, "Calculating your score...");
        setError(null);
        const result = computeResult();
        setLoading(false);
        if (result) {
          router.push(`/results`);
        } else {
          setError({ message: "Unable to compute result." });
        }
      } catch {
        setError({ message: "Failed to submit quiz. Please try again." });
        setLoading(false);
      } finally {
        setIsSubmitting(false);
      }
    } else {
      useQuizStore.setState({
        currentSession: {
          ...currentSession,
          currentQuestionIndex: currentQuestionIndex + 1,
        },
      });
    }
  };

  const handlePrevious = () => {
    if (currentQuestionIndex > 0) {
      useQuizStore.setState({
        currentSession: {
          ...currentSession,
          currentQuestionIndex: currentQuestionIndex - 1,
        },
      });
    }
  };

  return (
    <div className="min-h-screen bg-background">
      <Header showBackButton backHref="/" />

      <main className="max-w-2xl mx-auto px-6 py-12">
        {/* Progress */}
        <div className="mb-10">
          <div className="flex justify-between text-xs text-muted-foreground mb-2">
            <span>
              Question {currentQuestionIndex + 1} of {questions.length}
            </span>
            <span>{Math.round(progress)}%</span>
          </div>
          <Progress value={progress} className="h-1" />
        </div>

        {/* Error */}
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

        {/* Question */}
        <h1 className="text-2xl font-semibold tracking-tight mb-8 leading-snug">
          {question.text}
        </h1>

        {/* Options */}
        <div className="flex flex-col gap-2 mb-10">
          {question.options.map((option, i) => (
            <Card
              key={i}
              onClick={() =>
                !isSubmitting && recordAnswer(currentQuestionIndex, i)
              }
              className={`cursor-pointer transition-colors ${
                userAnswer === i
                  ? "border-foreground bg-foreground/5"
                  : "hover:border-foreground/40"
              }`}
            >
              <CardContent className="px-4 py-3 flex items-center gap-3">
                <span
                  className={`text-xs font-medium shrink-0 ${userAnswer === i ? "text-foreground" : "text-muted-foreground"}`}
                >
                  {String.fromCharCode(65 + i)}
                </span>
                <span className="text-sm">{option}</span>
              </CardContent>
            </Card>
          ))}
        </div>

        {/* Navigation */}
        <div className="flex gap-3">
          <Button
            variant="outline"
            onClick={handlePrevious}
            disabled={currentQuestionIndex === 0 || isSubmitting}
          >
            Previous
          </Button>
          <Button
            onClick={handleNext}
            disabled={userAnswer == null || isSubmitting}
          >
            {isSubmitting ? (
              <>
                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                Submitting…
              </>
            ) : isLast ? (
              "Submit"
            ) : (
              "Next"
            )}
          </Button>
        </div>
      </main>
    </div>
  );
}
