package services

import (
	"daymark/internal/models"
	"log"
	"time"

	"github.com/mmcdole/gofeed"
)

func FetchArticlesFromFeed(feedSources []models.FeedSource) []models.Article {
	parser := gofeed.NewParser()
	parser.UserAgent = "Mozilla/5.0 (compatible; QuizBot/1.0)"

	today := time.Now().Format("2006-01-02")
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")

	var articles []models.Article

	for _, fs := range feedSources {
		feed, err := parser.ParseURL(fs.URL)
		if err != nil {
			log.Printf("Failed to fetch RSS feed '%s': %v", fs.Name, err)
			continue
		}

		for _, item := range feed.Items {
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

			articles = append(articles, models.Article{
				Title:        item.Title,
				Link:         item.Link,
				Source:       fs.Name,
				PublishedAt:  pub,
				CategoryID:   fs.CategoryId,
				FeedSourceID: fs.ID,
			})
		}
	}

	return articles
}