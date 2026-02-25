package services

import (
	"daymark/internal/models"
	"fmt"
	"log"
	"net/http"
	"net/url"
	neturl "net/url"
	"os"
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
	 apiURL := fmt.Sprintf(
        "https://app.scrapingbee.com/api/v1/?api_key=%s&url=%s&render_js=false",
        os.Getenv("SCRAPINGBEE_KEY"),
        url.QueryEscape(rawURL),
    )
    resp, err := http.Get(apiURL)
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
