package action

import (
	"oko/pkg/ginapp/types"
)

type UpdateRequest struct {
	Type   *int32  `json:"type,omitempty"`
	ID     uint32  `json:"id"`
	Params *string `json:"params,omitempty"`
}

type UpdateRequestBody struct {
	Type   *int32  `json:"type,omitempty"`
	Params *string `json:"params,omitempty"`
}

type CreateRequest struct {
	Type   uint32 `json:"type" binding:"required"`
	Params string `json:"params" binding:"required"`
	RuleID uint32 `json:"rule_id"`
}

type AddRuleRequest struct {
	ID     uint32 `json:"action_id"`
	RuleID uint32 `json:"rule_id"`
}

type AddRuleRequestBody struct {
	RuleID uint32 `json:"rule_id"`
}

type ListFilter struct {
	RuleID *uint32 `json:"rule_id" form:"rule_id"`
}

type ListResponse struct {
	types.StdResponse
	Data []View `json:"data"`
}

type GetResponse struct {
	types.StdResponse
	Data View `json:"data"`
}
