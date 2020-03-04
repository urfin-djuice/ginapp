package rule

import (
	"oko/pkg/ginapp/types"
)

type CreateRequest struct {
	Host   string `json:"host,omitempty"`
	Status uint32 `json:"status,omitempty" binding:"CheckRuleStatus"`
}

type UpdateRequest struct {
	Host   string `json:"host,omitempty"`
	Status int    `json:"status,omitempty" binding:"CheckRuleStatus"`
}

type ListResponse struct {
	types.StdResponse
	Data []View `json:"data"`
}

type GetResponse struct {
	types.StdResponse
	Data View `json:"data"`
}
