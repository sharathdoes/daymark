package utils

import (
	"daymark/internal/models"
	"fmt"
	"log"
	"net/http"
	neturl "net/url"
	"regexp"
	"strings"
	"time"

	"github.com/go-shiori/go-readability"
	"github.com/mmcdole/gofeed"
)

func FetchArticlesFromFeeds(
	rssFeeds []models.FeedSource,
	linkExists func(link string) (bool, error),
) ([]models.Article, error){
	parser := gofeed.NewParser()
	parser.UserAgent = "Mozilla/5.0 (compatible; QuizBot/1.0)"
	var articles []models.Article
	today := time.Now().Format("2006-01-02")
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	for _, feed := range rssFeeds {
		f, err := parser.ParseURL(feed.URL)
		if err != nil {
			log.Printf("Failed to fetch RSS feed '%s' (%s): %v", feed.Name, feed.URL, err)
			continue
		}
		for _, item := range f.Items {

			if item.PublishedParsed != nil {
				pubDate := item.PublishedParsed.Format("2006-01-02")
				if pubDate != today && pubDate != yesterday {
					continue
				}
			}

			if item.Link == "" {
				continue
			}
			exists, err := linkExists(item.Link)
			if err != nil {
				log.Printf("Failed to check duplicate for '%s': %v", item.Link, err)
				continue
			}
			if exists {
				continue // 🚀 skip completely — no HTTP, no parsing
			}

			pub := time.Now()
			if item.PublishedParsed != nil {
				pub = *item.PublishedParsed
			}

			content, err := extractArticleContent(item.Link)
			if err != nil {
				log.Printf("Failed to extract article content for '%s': %v", item.Link, err)
				content = ""
			}

			articles = append(articles, models.Article{
				Title:       item.Title,
				Link:        item.Link,
				Source:      feed.Name,
				Category:    strings.Join(feed.Category, ","),
				PublishedAt: pub,
				Content:     content,
				FeedSourceID: feed.ID,
			})
		}
	}
	return articles, nil
}

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

func cleanText(text string) string {
	text = strings.TrimSpace(text)
	re := regexp.MustCompile(`\s+`)
	return re.ReplaceAllString(text, " ")
}
