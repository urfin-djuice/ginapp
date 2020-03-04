package repost

import (
	"oko/pkg/ginapp/types"
	"time"
)

type ListRequest struct {
	types.PaginationRequest
}

type RequestForm struct {
	URL      string     `json:"url" form:"url" binding:"required,ExistsRepostRequest"`
	DateFrom *time.Time `json:"date_from" form:"date_from"`
	DateTo   *time.Time `json:"date_to" form:"date_to"`
}

type NewRequestForm struct {
	URL string `json:"url" form:"url" binding:"required"`
}
