import { Category, QuizSession, QuizResult, ApiError, Question } from '@/types';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api';

async function handleResponse<T>(response: Response): Promise<T> {
  if (!response.ok) {
    const error: ApiError = {
      message: `API Error: ${response.statusText}`,
      code: response.status.toString(),
    };
    throw error;
  }
  return response.json();
}

export async function getCategories(): Promise<Category[]> {
  try {
    const response = await fetch(`${API_BASE_URL}/categories`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
    });
    return handleResponse<Category[]>(response);
  } catch (error) {
    throw error as ApiError;
  }
}

export async function generateQuiz(
  categoryId: string,
  difficulty: string,
  numberOfQuestions: number = 5
): Promise<QuizSession> {
  try {
    const response = await fetch(`${API_BASE_URL}/quizzes/generate`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        category_id: categoryId,
        difficulty,
        number_of_questions: numberOfQuestions,
      }),
    });
    return handleResponse<QuizSession>(response);
  } catch (error) {
    throw error as ApiError;
  }
}

export async function submitQuiz(session: QuizSession): Promise<QuizResult> {
  try {
    const response = await fetch(`${API_BASE_URL}/quizzes/${session.id}/submit`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        answers: session.answers,
      }),
    });
    return handleResponse<QuizResult>(response);
  } catch (error) {
    throw error as ApiError;
  }
}

export async function getQuizResult(quizId: string): Promise<QuizResult> {
  try {
    const response = await fetch(`${API_BASE_URL}/quizzes/${quizId}/results`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
    });
    return handleResponse<QuizResult>(response);
  } catch (error) {
    throw error as ApiError;
  }
}
