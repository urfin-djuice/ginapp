package action

import (
	"oko/pkg/util"
	pb "oko/srv/proxy/proto"
	"time"
)

type View struct {
	ID        uint32   `json:"id"`
	Type      int32    `json:"type"`
	Params    string   `json:"params"`
	RuleID    []uint32 `json:"rule_id"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at,omitempty"`
	DeletedAt string   `json:"deleted_at,omitempty"`
}

func toView(act pb.Action) (res View) {
	res.ID = act.Id
	if act.UpdatedAt != nil {
		res.UpdatedAt = util.TimestampToString(*act.UpdatedAt, time.RFC3339)
	}
	if act.DeletedAt != nil {
		res.DeletedAt = util.TimestampToString(*act.DeletedAt, time.RFC3339)
	}
	if act.CreatedAt != nil {
		res.CreatedAt = util.TimestampToString(*act.CreatedAt, time.RFC3339)
	}
	res.RuleID = make([]uint32, 0, len(act.RuleIds))
	if act.RuleIds != nil {
		res.RuleID = append(res.RuleID, act.RuleIds...)
	}
	res.Type = int32(act.Type)
	res.Params = act.Params
	return
}
