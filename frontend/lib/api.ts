import { Category, QuizSession, ApiError } from '@/types';

// Backend base URL (Gin server). Can be overridden with NEXT_PUBLIC_API_URL.
const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

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
  const response = await fetch(`${API_BASE_URL}/category/`, {
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
  categoryIds: string[],
  difficulty: string,
  numberOfQuestions: number = 5
): Promise<QuizSession> {
  try {
  const response = await fetch(`${API_BASE_URL}/quiz/generate`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
    category_ids: categoryIds.map((id) => Number(id)),
    difficulty,
    number_of_questions: numberOfQuestions,
      }),
    });
    console.log(response)
    return handleResponse<QuizSession>(response);
  } catch (error) {
    throw error as ApiError;
  }
}

// submitQuiz and getQuizResult are now handled entirely on the frontend
// via quizStore.computeResult and the /results page.
