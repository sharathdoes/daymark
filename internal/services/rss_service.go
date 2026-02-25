package services

import (
	"daymark/internal/models"
	"fmt"
	"log"
	"net/http"
	neturl "net/url"
	"strings"
	"time"

	"github.com/go-shiori/go-readability"
	"github.com/mmcdole/gofeed"
)

func FetchArticlesFromFeeds(feedSources []models.FeedSource) ([]models.Article, error) {
	parser := gofeed.NewParser()
	parser.UserAgent = "Mozilla/5.0 (compatible; QuizBot/1.0)"

	today := time.Now().Format("2006-01-02")
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")

	var articles []models.Article

	for i := 0; i < len(feedSources); i++ {
		feed, err := parser.ParseURL(feedSources[i].URL)
		if err != nil {
			return nil, err
		}

		itemCount := 0
		for _, item := range feed.Items {
			if itemCount >= 5 {
				break
			}

			if item.Link == "" {
				continue
			}

			if item.PublishedParsed != nil {
				pubDate := item.PublishedParsed.Format("2006-01-02")
				if pubDate != today && pubDate != yesterday {
					continue
				}
			}

			pub := time.Now()
			if item.PublishedParsed != nil {
				pub = *item.PublishedParsed
			}

			
			content, err := extractArticleContent(item.Link)
			if err != nil {
				log.Printf("Skipping article due to scrape error: %v", err)
				continue
			}
			var categoryID uint
			if len(feedSources[i].Categories) > 0 {
				categoryID = feedSources[i].Categories[0].ID
			}

			articles = append(articles, models.Article{
				Title:        strings.ToValidUTF8(item.Title, ""),
				Link:         item.Link,
				Source:       strings.ToValidUTF8(feedSources[i].Name, ""),
				Content:      strings.ToValidUTF8(content, ""),
				PublishedAt:  pub,
				CategoryID:   categoryID,
				FeedSourceID: feedSources[i].ID,
			})
			itemCount++
		}
	}

	return articles, nil
}

func extractArticleContent(rawURL string) (string, error) {
	client := &http.Client{Timeout: 15 * time.Second}

	req, err := http.NewRequest("GET", rawURL, nil)
	if err != nil {
		return "", err
	}
req.Header.Set("Referer", "https://google.com")
req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/121.0.0.0 Safari/537.36")
req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
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
