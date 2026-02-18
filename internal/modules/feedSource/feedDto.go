package feedSource

import "github.com/lib/pq"

type CreateFeedDTO struct {
	Name     string         `json:"name"`
	URL      string         `json:"url"`
	Category pq.StringArray `json:"category"`
}
