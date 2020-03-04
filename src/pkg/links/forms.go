package links

import "oko/pkg/ginapp/types"

type ListRequest struct {
	types.PaginationRequest
	Query    string `json:"query" form:"query" binding:"required"`
	From     string `json:"from" form:"from"`
	To       string `json:"to" form:"to"`
	DomainID string `json:"domain_id" form:"domain_id"`
}

type RePostRequest struct {
	URL string `json:"url" form:"url" binding:"required"`
}
