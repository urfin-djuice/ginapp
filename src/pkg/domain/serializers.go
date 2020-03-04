package domain

import (
	"time"

	"github.com/thoas/go-funk"
)

type Serializer struct {
	Domain Domain
}

type ListSerializer struct {
	Domains []*Domain
}

type RssResponse struct {
	ID        uint       `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeleteAt  *time.Time `json:"delete_at,omitempty"`
	Link      string     `json:"link"`
}

type Response struct {
	ID               uint            `json:"id"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at"`
	Name             string          `json:"name"`
	RssLinks         *[]*RssResponse `json:"rss_link,omitempty"`
	TelegramUsername string          `json:"telegram_username"`
	Error            string          `json:"error"`
}

func toRssResponse(dom *Domain) *[]*RssResponse {
	if dom == nil {
		return nil
	}
	if dom.Rss == nil {
		return nil
	}
	res := make([]*RssResponse, 0, len(dom.Rss))
	for _, rss := range dom.Rss {
		if rss != nil {
			res = append(res, &RssResponse{
				ID:        rss.ID,
				CreatedAt: rss.CreatedAt,
				UpdatedAt: rss.UpdatedAt,
				DeleteAt:  rss.DeletedAt,
				Link:      rss.Link,
			})
		}
	}
	return &res
}

func (s *Serializer) To() Response {
	return Response{
		ID:               s.Domain.ID,
		CreatedAt:        s.Domain.CreatedAt,
		UpdatedAt:        s.Domain.UpdatedAt,
		Name:             s.Domain.Name,
		TelegramUsername: s.Domain.TelegramUsername,
		Error:            s.Domain.Error,
		RssLinks:         toRssResponse(&s.Domain),
	}
}

func (s *ListSerializer) To() []*Response {
	data := funk.Map(s.Domains, func(model *Domain) *Response {
		return &Response{
			ID:               model.ID,
			CreatedAt:        model.CreatedAt,
			UpdatedAt:        model.UpdatedAt,
			Name:             model.Name,
			TelegramUsername: model.TelegramUsername,
			Error:            model.Error,
			RssLinks:         toRssResponse(model),
		}
	}).([]*Response)

	return data
}
