package workers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	neturl "net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/joho/godotenv"

	"github.com/go-shiori/go-readability"
	"github.com/mmcdole/gofeed"
)

///////////////////////////////////////////////////////////
// TYPES
///////////////////////////////////////////////////////////

type Article struct {
	Title     string
	Link      string
	Published time.Time
	Source    string
	Content   string
}

type Quiz struct {
	Title     string     `json:"title"`
	Source    string     `json:"source"`
	Date      time.Time  `json:"date"`
	Questions []Question `json:"questions"`
}

type Question struct {
	Question string   `json:"question"`
	Options  []string `json:"options"`
	Answer   int      `json:"answer"`
}

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

///////////////////////////////////////////////////////////
// CONFIG
///////////////////////////////////////////////////////////

var rssFeeds = map[string]string{
	"TOI":           "https://timesofindia.indiatimes.com/rssfeedstopstories.cms",
	"TheHindu":      "https://www.thehindu.com/feeder/default.rss",
	"IndianExpress": "https://indianexpress.com/section/trending/top-10-listing/feed/",
}

var groqURL = "https://api.groq.com/openai/v1/chat/completions"

///////////////////////////////////////////////////////////
// MAIN LOGIC
///////////////////////////////////////////////////////////

func extractArticleContent(rawURL string) (string, error) {
	client := &http.Client{Timeout: 15 * time.Second}

	req, err := http.NewRequest("GET", rawURL, nil)
	if err != nil {
		return "", fmt.Errorf("invalid URL %s: %w", rawURL, err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP %d for URL: %s", resp.StatusCode, rawURL)
	}

	parsedURL, err := neturl.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse URL: %w", err)
	}

	article, err := readability.FromReader(resp.Body, parsedURL)
	if err != nil {
		return "", fmt.Errorf("readability parse failed: %w", err)
	}

	text := cleanText(article.TextContent)
	if len(text) > 4000 {
		text = text[:4000]
	}

	return text, nil
}


func GenerateQuiz() (*Quiz, error) {
	_ = godotenv.Load()

	// Load API key from environment variable (never hardcode secrets)
	apiKey := os.Getenv("GROQ_API_KEY")
	if apiKey == "" {
		return nil, errors.New("GROQ_API_KEY environment variable is not set")
	}

	articles, err := fetchArticlesFromFeeds()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch articles: %w", err)
	}

	if len(articles) == 0 {
		return nil, errors.New("no articles found from RSS feeds (check date filter or feed URLs)")
	}

	log.Printf("Fetched %d articles from feeds", len(articles))

	var allQuestions []Question

	for i, article := range articles {
		if i >= 5 { // Try up to 5 articles to get enough questions
			break
		}

		log.Printf("Processing article %d: %s", i+1, article.Title)

		text, err := extractArticleContent(article.Link)
		if err != nil {
			log.Printf("Failed to extract content from %s: %v", article.Link, err)
			continue
		}

		if len(text) < 200 {
			log.Printf("Article too short (%d chars), skipping: %s", len(text), article.Title)
			continue
		}

		questions, err := generateQuestionsWithGroq(text, apiKey)
		if err != nil {
			log.Printf("Groq generation failed for '%s': %v", article.Title, err)
			continue
		}

		log.Printf("Generated %d questions from article: %s", len(questions), article.Title)
		allQuestions = append(allQuestions, questions...)

		// Stop once we have enough questions
		if len(allQuestions) >= 10 {
			break
		}
	}

	if len(allQuestions) == 0 {
		return nil, errors.New("no questions generated — check article content or Groq API key")
	}

	return &Quiz{
		Title:     "Daily Current Affairs Quiz",
		Source:    "Multiple News Sources",
		Date:      time.Now(),
		Questions: allQuestions,
	}, nil
}

///////////////////////////////////////////////////////////
// RSS FETCH
///////////////////////////////////////////////////////////

func fetchArticlesFromFeeds() ([]Article, error) {
	parser := gofeed.NewParser()

	var articles []Article

	today := time.Now().Format("2006-01-02")
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")

	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	for source, url := range rssFeeds {

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Printf("Failed to create request for '%s' (%s): %v", source, url, err)
			continue
		}

		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")
		req.Header.Set("Accept", "application/rss+xml, application/xml;q=0.9, */*;q=0.8")
		req.Header.Set("Accept-Language", "en-US,en;q=0.9")

		resp, err := client.Do(req)
		if err != nil {
			log.Printf("HTTP request failed for '%s' (%s): %v", source, url, err)
			continue
		}

		if resp.StatusCode != http.StatusOK {
			log.Printf("HTTP status %d for '%s' (%s)", resp.StatusCode, source, url)
			resp.Body.Close()
			continue
		}

		feed, err := parser.Parse(resp.Body)
		resp.Body.Close()

		if err != nil {
			log.Printf("Failed to parse RSS feed '%s' (%s): %v", source, url, err)
			continue
		}

		log.Printf("Feed '%s': %d items found", source, len(feed.Items))

		for _, item := range feed.Items {

			// Accept articles with no date or from today/yesterday
			if item.PublishedParsed != nil {
				pubDate := item.PublishedParsed.Format("2006-01-02")
				if pubDate != today && pubDate != yesterday {
					continue
				}
			}

			if item.Link == "" {
				continue
			}

			pub := time.Now()
			if item.PublishedParsed != nil {
				pub = *item.PublishedParsed
			}

			articles = append(articles, Article{
				Title:     item.Title,
				Link:      item.Link,
				Published: pub,
				Source:    source,
			})
		}
	}

	return articles, nil
}

///////////////////////////////////////////////////////////
// ARTICLE EXTRACTION
///////////////////////////////////////////////////////////

func cleanText(text string) string {
	text = strings.TrimSpace(text)
	re := regexp.MustCompile(`\s+`)
	return re.ReplaceAllString(text, " ")
}

///////////////////////////////////////////////////////////
// GROQ GENERATION
///////////////////////////////////////////////////////////

func generateQuestionsWithGroq(articleText string, apiKey string) ([]Question, error) {
	prompt := fmt.Sprintf(`You must return ONLY valid JSON. No markdown, no explanation, no extra text.

Return exactly this format:
[
  {
    "question": "Question text here?",
    "options": ["Option A", "Option B", "Option C", "Option D"],
    "answer": 0
  }
]

The "answer" field is the zero-based index of the correct option (0, 1, 2, or 3).

Generate 3 UPSC-style multiple choice questions based on the following article. Focus on facts, events, people, places, and policies mentioned.

Article:
%s`, articleText)

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
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(
		context.Background(),
		"POST",
		groqURL,
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Groq API request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Groq API returned HTTP %d: %s", resp.StatusCode, string(body))
	}

	var groqResp GroqResponse
	if err := json.Unmarshal(body, &groqResp); err != nil {
		return nil, fmt.Errorf("failed to parse Groq response: %w\nBody: %s", err, string(body))
	}

	// Check for API-level error in response
	if groqResp.Error != nil {
		return nil, fmt.Errorf("Groq API error: %s (type: %s)", groqResp.Error.Message, groqResp.Error.Type)
	}

	if len(groqResp.Choices) == 0 {
		return nil, fmt.Errorf("empty choices from Groq, full response: %s", string(body))
	}

	raw := cleanLLMResponse(groqResp.Choices[0].Message.Content)
	if raw == "" {
		return nil, errors.New("LLM response was empty after cleaning")
	}

	var questions []Question
	if err := json.Unmarshal([]byte(raw), &questions); err != nil {
		log.Printf("JSON Parse Error: %v", err)
		log.Printf("RAW LLM OUTPUT: %s", raw)
		return nil, fmt.Errorf("failed to parse questions JSON: %w", err)
	}

	// Validate questions
	questions = validateQuestions(questions)
	if len(questions) == 0 {
		return nil, errors.New("all generated questions were invalid")
	}

	return questions, nil
}

///////////////////////////////////////////////////////////
// VALIDATE QUESTIONS
///////////////////////////////////////////////////////////

func validateQuestions(questions []Question) []Question {
	var valid []Question
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

///////////////////////////////////////////////////////////
// CLEAN LLM RESPONSE
///////////////////////////////////////////////////////////

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
