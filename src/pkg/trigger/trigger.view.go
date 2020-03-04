package trigger

import (
	"oko/pkg/util"
	pb "oko/srv/proxy/proto"
	"time"
)

type View struct {
	ID        uint32 `json:"id"`
	URL       string `json:"url"`
	Params    string `json:"params"`
	Type      int    `json:"type"`
	RuleID    uint32 `json:"rule_id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at,omitempty"`
	DeletedAt string `json:"deleted_at,omitempty"`
}

func toView(trigger *pb.Trigger) (res View) {
	res = View{
		ID:     trigger.Id,
		URL:    trigger.Url,
		Params: trigger.Params,
		Type:   int(trigger.Type),
		RuleID: trigger.RuleId,
	}
	if trigger.CreatedAt != nil {
		res.CreatedAt = util.TimestampToString(*trigger.CreatedAt, time.RFC3339)
	}
	if trigger.UpdatedAt != nil {
		res.UpdatedAt = util.TimestampToString(*trigger.UpdatedAt, time.RFC3339)
	}
	if trigger.DeletedAt != nil {
		res.DeletedAt = util.TimestampToString(*trigger.DeletedAt, time.RFC3339)
	}
	return
}
