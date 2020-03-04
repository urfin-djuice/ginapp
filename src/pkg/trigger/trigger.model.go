package trigger

import (
	"oko/pkg/ginapp/types"
	pb "oko/srv/proxy/proto"
)

type CreateRequest struct {
	URL    string `json:"url"`
	Type   int    `json:"type" binding:"required"`
	Params string `json:"params" binding:"required"`
	RuleID uint32 `json:"rule_id" binding:"required"`
}

type UpdateRequest struct {
	Type   *int    `json:"type,omitempty"`
	Params *string `json:"params,omitempty"`
	RuleID *uint32 `json:"rule_id,omitempty"`
	URL    *string `json:"url"`
}

type ListFilter struct {
	RuleID *uint32 `json:"rule_id" form:"rule_id"`
}

var triggerTypes = map[int]pb.TriggerType{
	1: pb.TriggerType_status,
	2: pb.TriggerType_body,
	3: pb.TriggerType_xpath_include,
	4: pb.TriggerType_xpath_not_include,
}

func ToTriggerType(typeInt int) pb.TriggerType {
	tp, ok := triggerTypes[typeInt]
	if !ok {
		return pb.TriggerType_unset
	}
	return tp
}

type GetResponse struct {
	types.StdResponse
	Data View `json:"data"`
}

type ListResponse struct {
	types.StdResponse
	Data []View `json:"data"`
}
