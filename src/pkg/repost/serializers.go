package repost

import (
	"net/url"
	"oko/pkg/ginapp/types"
	"time"

	"github.com/thoas/go-funk"
)

type Serializer struct {
	Request Request
}

type ListSerializer struct {
	Requests []*Request
}

type RequestResponse struct {
	ID           uint      `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Error        string    `json:"error"`
	URL          string    `json:"url"`
	HasProcessed *bool     `json:"has_processed"`
	Domain       string    `json:"domain"`
	CountRePost  int       `json:"count_repost"`

	Links []*LinkResponse `json:"links"`
}

type LinkResponse struct {
	ID          uint      `json:"id"`
	URL         string    `json:"url"`
	PublishedAt time.Time `json:"published_at"`
	Title       string    `json:"title"`
	Domain      string    `json:"domain"`
}

type NewRequestResponse struct {
	types.StdResponse
	Data RequestResponse `json:"data"`
}

func (s Serializer) To() RequestResponse {
	data := funk.Map(s.Request.Links, func(model Link) *LinkResponse {
		tmp := &LinkResponse{
			ID:          model.ID,
			URL:         model.URL,
			PublishedAt: model.PublishedAt,
			Title:       model.Title,
		}

		u, err := url.Parse(model.URL)
		if err == nil {
			tmp.Domain = u.Host
		}

		return tmp
	}).([]*LinkResponse)

	tmp := RequestResponse{
		ID:           s.Request.ID,
		CreatedAt:    s.Request.CreatedAt,
		UpdatedAt:    s.Request.UpdatedAt,
		URL:          s.Request.URL,
		Error:        s.Request.Error,
		CountRePost:  len(data),
		Links:        data,
		HasProcessed: s.Request.HasProcessed,
	}

	u, err := url.Parse(s.Request.URL)
	if err == nil {
		tmp.Domain = u.Host
	}

	return tmp
}

func (s *ListSerializer) To() []RequestResponse {
	data := funk.Map(s.Requests, func(model *Request) RequestResponse {
		return Serializer{*model}.To()
	}).([]RequestResponse)

	return data
}
