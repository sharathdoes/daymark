 package utils

// import (
// 	"bytes"
// 	"context"
// 	"daymark/internal/models"
// 	"encoding/json"
// 	"errors"
// 	"fmt"
// 	"io"
// 	"log"
// 	"net/http"
// 	"strings"
// 	"time"
// )

// type GroqMessage struct {
// 	Role    string `json:"role"`
// 	Content string `json:"content"`
// }

// type GroqRequest struct {
// 	Model    string        `json:"model"`
// 	Messages []GroqMessage `json:"messages"`
// }

// type GroqChoice struct {
// 	Message GroqMessage `json:"message"`
// }

// type GroqResponse struct {
// 	Choices []GroqChoice `json:"choices"`
// 	Error   *GroqError   `json:"error,omitempty"`
// }

// type GroqError struct {
// 	Message string `json:"message"`
// 	Type    string `json:"type"`
// }

// type llmQuestion struct {
// 	Question string   `json:"question"`
// 	Options  []string `json:"options"`
// 	Answer   int      `json:"answer"`
// }

// var groqURL = "https://api.groq.com/openai/v1/chat/completions"

// func GenerateQuestions(articleText string, apiKey string) ([]models.Question, error) {
// 	if strings.TrimSpace(articleText) == "" {
// 		return nil, errors.New("article text cannot be empty")
// 	}

// 	if strings.TrimSpace(apiKey) == "" {
// 		return nil, errors.New("api key cannot be empty")
// 	}

// 	prompt := fmt.Sprintf(`You must return ONLY valid JSON. No markdown, no explanation, no extra text.

// Return exactly this format:
// [
//   {
//     "question": "Question text here?",
//     "options": ["Option A", "Option B", "Option C", "Option D"],
//     "answer": 0
//   }
// ]

// The "answer" field is the zero-based index of the correct option (0, 1, 2, or 3).

// Generate 3 UPSC-style multiple choice questions based on the following article. Focus on facts, events, people, places, and policies mentioned.

// Article:
// %s`, articleText)

// 	requestBody := GroqRequest{
// 		Model: "llama-3.3-70b-versatile",
// 		Messages: []GroqMessage{
// 			{
// 				Role:    "user",
// 				Content: prompt,
// 			},
// 		},
// 	}

// 	jsonData, err := json.Marshal(requestBody)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to marshal request: %w", err)
// 	}

// 	req, err := http.NewRequestWithContext(
// 		context.Background(),
// 		"POST",
// 		groqURL,
// 		bytes.NewBuffer(jsonData),
// 	)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to create request: %w", err)
// 	}

// 	req.Header.Set("Authorization", "Bearer "+apiKey)
// 	req.Header.Set("Content-Type", "application/json")

// 	client := &http.Client{Timeout: 30 * time.Second}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return nil, fmt.Errorf("Groq API request failed: %w", err)
// 	}
// 	defer resp.Body.Close()

// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to read response body: %w", err)
// 	}

// 	if resp.StatusCode != http.StatusOK {
// 		return nil, fmt.Errorf("Groq API returned HTTP %d: %s", resp.StatusCode, string(body))
// 	}

// 	var groqResp GroqResponse
// 	if err := json.Unmarshal(body, &groqResp); err != nil {
// 		return nil, fmt.Errorf("failed to parse Groq response: %w\nBody: %s", err, string(body))
// 	}

// 	// Check for API-level error in response
// 	if groqResp.Error != nil {
// 		return nil, fmt.Errorf("Groq API error: %s (type: %s)", groqResp.Error.Message, groqResp.Error.Type)
// 	}

// 	if len(groqResp.Choices) == 0 {
// 		return nil, fmt.Errorf("empty choices from Groq, full response: %s", string(body))
// 	}

// 	raw := cleanLLMResponse(groqResp.Choices[0].Message.Content)
// 	if raw == "" {
// 		return nil, errors.New("LLM response was empty after cleaning")
// 	}

// 	var generated []llmQuestion
// 	if err := json.Unmarshal([]byte(raw), &generated); err != nil {
// 		log.Printf("JSON Parse Error: %v", err)
// 		log.Printf("RAW LLM OUTPUT: %s", raw)
// 		return nil, fmt.Errorf("failed to parse questions JSON: %w", err)
// 	}

// 	valid := validateQuestions(generated)
// 	if len(valid) == 0 {
// 		return nil, errors.New("all generated questions were invalid")
// 	}

// 	questions := make([]models.Question, 0, len(valid))
// 	for _, q := range valid {
// 		questions = append(questions, models.Question{
// 			Text:    q.Question,
// 			OptionA: q.Options[0],
// 			OptionB: q.Options[1],
// 			OptionC: q.Options[2],
// 			OptionD: q.Options[3],
// 			Answer:  q.Answer,
// 		})
// 	}

// 	return questions, nil
// }

// func validateQuestions(questions []llmQuestion) []llmQuestion {
// 	valid := make([]llmQuestion, 0, len(questions))

// 	for _, q := range questions {
// 		if strings.TrimSpace(q.Question) == "" {
// 			log.Println("Skipping question with empty text")
// 			continue
// 		}

// 		if len(q.Options) != 4 {
// 			log.Printf("Skipping question with %d options (expected 4): %s", len(q.Options), q.Question)
// 			continue
// 		}

// 		if q.Answer < 0 || q.Answer > 3 {
// 			log.Printf("Skipping question with invalid answer index %d: %s", q.Answer, q.Question)
// 			continue
// 		}

// 		valid = append(valid, q)
// 	}

// 	return valid
// }

// func cleanLLMResponse(resp string) string {
// 	resp = strings.TrimSpace(resp)

// 	start := strings.Index(resp, "[")
// 	end := strings.LastIndex(resp, "]")

// 	if start != -1 && end != -1 && end > start {
// 		return strings.TrimSpace(resp[start : end+1])
// 	}

// 	if strings.HasPrefix(resp, "```") {
// 		resp = strings.TrimPrefix(resp, "```json")
// 		resp = strings.TrimPrefix(resp, "```")
// 		if idx := strings.LastIndex(resp, "```"); idx != -1 {
// 			resp = resp[:idx]
// 		}
// 	}

// 	return strings.TrimSpace(resp)
// }
