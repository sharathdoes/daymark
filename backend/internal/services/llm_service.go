package services

import (
	"bytes"
	"context"
	"daymark/internal/models"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type GroqMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type GroqRequest struct {
	Model    string        `json:"model"`
	Messages []GroqMessage `json:"messages"`
}

type GroqChoice struct {
	Message GroqMessage `json:"message"`
}

type GroqResponse struct {
	Choices []GroqChoice `json:"choices"`
	Error   *GroqError   `json:"error,omitempty"`
}

type GroqError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}

type llmQuestion struct {
	Title    string   `json:"title"`
	Question string   `json:"question"`
	Options  []string `json:"options"`
	Answer   int      `json:"answer"`
}

var groqURL = "https://api.groq.com/openai/v1/chat/completions"

func cleanText(text string) string {
	text = strings.TrimSpace(text)
	re := regexp.MustCompile(`\s+`)
	return re.ReplaceAllString(text, " ")
}

func GenerateQuiz(NumberOfQuestions int, categoryIds []uint, apiKey string, difficulty string, articles []models.Article) (*models.Quiz, error) {
	var allQuestions []models.Question

	var quizTitle string
	r := rand.New(rand.NewSource(rand.Int63()))
	r.Shuffle(len(articles), func(i, j int) {
		articles[i], articles[j] = articles[j], articles[i]
	})

	for i, article := range articles {
		if i >= 8 {
			break
		}

		generatedTitle, questions, err := generateQuestionsWithGroq(article.Link, article.Title, article.Content, apiKey)
		if err != nil {
			log.Printf("Groq generation failed for '%s': %v", article.Title, err)
			continue
		}

		if generatedTitle != "" && quizTitle == "" {
			quizTitle = generatedTitle
		}
		allQuestions = append(allQuestions, questions...)

		if len(allQuestions) >= NumberOfQuestions {
			break
		}
	}

	if len(allQuestions) == 0 {
		return &models.Quiz{}, errors.New("no questions generated — check article content or Groq API key")
	}

	if quizTitle == "" && len(articles) > 0 {
		quizTitle = fmt.Sprintf("Quiz on %s", articles[0].Title)
	}

	return &models.Quiz{
		Title:       quizTitle,
		CategoryIDs: categoryIds,
		Difficulty:  difficulty,
		Questions:   allQuestions,
	}, nil
}

func generateQuestionsWithGroq(articleURL string, articleTitle string, articleText string, apiKey string) (string, []models.Question, error) {
	prompt := fmt.Sprintf(`You must return ONLY valid JSON. No markdown, no explanation, no extra text.

Return exactly this format:
[
	{
		"title": "Catchy category-based quiz title (e.g., Polity Pulse, Economy Check, World Watch, Science Snap, Environment Brief)",
		"question": "Direct factual question from the article?",
		"options": ["Option A", "Option B", "Option C", "Option D"],
		"answer": 0
	}
]

The "answer" field must be the zero-based index of the correct option (0, 1, 2, or 3).

Generate EXACTLY 3 simple UPSC-style multiple choice questions based ONLY on clearly stated facts in the article.

IMPORTANT RULES:
- Do NOT ask deep analytical, opinion-based, or conceptual questions.
- Do NOT frame questions on theoretical background beyond what is mentioned.
- Keep questions straightforward and factual.
- Each question must focus on a DIFFERENT fact, event, person, place, scheme, date, or policy mentioned in the article.
- Cover different parts of the article quickly. Do not stay on the same sub-topic.
- The title must be catchy and category-based, NOT about the specific subject of the article.
- Keep questions concise and exam-oriented.

Article title: %s
Article content: %s`, articleTitle, articleText)

	requestBody := GroqRequest{
		Model: "llama-3.3-70b-versatile",
		Messages: []GroqMessage{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(
		context.Background(),
		"POST",
		groqURL,
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return "", nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", nil, fmt.Errorf("Groq API request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", nil, fmt.Errorf("Groq API returned HTTP %d", resp.StatusCode)
	}

	var groqResp GroqResponse
	if err := json.Unmarshal(body, &groqResp); err != nil {
		return "", nil, fmt.Errorf("failed to parse Groq response: %w", err)
	}

	// Check for API-level error in response
	if groqResp.Error != nil {
		return "", nil, fmt.Errorf("Groq API error: %s (type: %s)", groqResp.Error.Message, groqResp.Error.Type)
	}

	if len(groqResp.Choices) == 0 {
		return "", nil, fmt.Errorf("empty choices from Groq")
	}

	raw := cleanLLMResponse(groqResp.Choices[0].Message.Content)
	if raw == "" {
		return "", nil, errors.New("LLM response was empty after cleaning")
	}

	var generated []llmQuestion
	if err := json.Unmarshal([]byte(raw), &generated); err != nil {
		log.Printf("JSON Parse Error: %v", err)
		return "", nil, fmt.Errorf("failed to parse questions JSON: %w", err)
	}

	// Validate questions
	valid := validateQuestions(generated)
	if len(valid) == 0 {
		return "", nil, errors.New("all generated questions were invalid")
	}

	quizTitle := strings.TrimSpace(valid[0].Title)
	if quizTitle == "" {
		quizTitle = articleTitle
	}

	questions := make([]models.Question, 0, len(valid))
	for _, q := range valid {
		questions = append(questions, models.Question{
			Question:   q.Question,
			Options:    q.Options,
			Answer:     q.Answer,
			ArticleURL: articleURL,
		})
	}

	return quizTitle, questions, nil
}

func validateQuestions(questions []llmQuestion) []llmQuestion {
	var valid []llmQuestion
	for _, q := range questions {
		if q.Question == "" {
			log.Println("Skipping question with empty text")
			continue
		}
		if len(q.Options) != 4 {
			log.Printf("Skipping question with %d options (expected 4): %s", len(q.Options), q.Question)
			continue
		}
		if q.Answer < 0 || q.Answer > 3 {
			log.Printf("Skipping question with invalid answer index %d: %s", q.Answer, q.Question)
			continue
		}
		valid = append(valid, q)
	}
	return valid
}

func cleanLLMResponse(resp string) string {
	resp = strings.TrimSpace(resp)

	// Most robust approach: extract the JSON array directly
	start := strings.Index(resp, "[")
	end := strings.LastIndex(resp, "]")

	if start != -1 && end != -1 && end > start {
		return strings.TrimSpace(resp[start : end+1])
	}

	// Fallback: strip markdown code fences
	if strings.HasPrefix(resp, "```") {
		resp = strings.TrimPrefix(resp, "```json")
		resp = strings.TrimPrefix(resp, "```")
		if idx := strings.LastIndex(resp, "```"); idx != -1 {
			resp = resp[:idx]
		}
	}

	return strings.TrimSpace(resp)
}
