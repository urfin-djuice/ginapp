package links

import (
	"oko/pkg/domain"
	"oko/pkg/ginapp/types"
	"time"
)

type Serializer struct {
	Link *Link
}

type Response struct {
	ID                  uint            `json:"id"`
	URL                 string          `json:"url"`
	Domain              domain.Response `json:"domain"`
	SentimentalScore    *float32        `json:"sentimental_score"`
	SentimentalPositive *float32        `json:"sentimental_positive"`
	SentimentalNegative *float32        `json:"sentimental_negative"`
	Content             string          `json:"content"`
	CreatedAt           time.Time       `json:"created_at"`
	PublishedAt         *time.Time      `json:"published_at"`
}

type MetaListResponse struct {
	types.PaginationResponse
	NegativeCount uint32  `json:"negative_count"`
	PositiveCount uint32  `json:"positive_count"`
	NeutralCount  uint32  `json:"neutral_count"`
	Image         *string `json:"image"`
}

type LinkResponse struct {
	URL         string    `json:"url"`
	PublishedAt time.Time `json:"published_at"`
	CountRePost int       `json:"count_repost"`
}

type RePostResponse struct {
	Origin string         `json:"origin"`
	Links  []LinkResponse `json:"links"`
}

func (s *Serializer) To(content string) *Response {
	serializer := domain.Serializer{Domain: s.Link.Domain}

	return &Response{
		Domain:              serializer.To(),
		ID:                  s.Link.ID,
		URL:                 s.Link.URL,
		SentimentalScore:    s.Link.SentimentalScore,
		SentimentalPositive: s.Link.SentimentalPositive,
		SentimentalNegative: s.Link.SentimentalNegative,
		Content:             content,
		PublishedAt:         s.Link.PublishedAt,
		CreatedAt:           *s.Link.CreatedAt,
	}
}
