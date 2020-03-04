package rule

import (
	"oko/pkg/util"
	pb "oko/srv/proxy/proto"
	"time"
)

type View struct {
	ID        uint32        `json:"id"`
	Host      string        `json:"host"`
	Status    uint32        `json:"status"`
	CreatedAt string        `json:"created_at"`
	UpdatedAt string        `json:"updated_at"`
	DeletedAt string        `json:"deleted_at,omitempty"`
	Triggers  []TriggerView `json:"triggers"`
	Actions   []ActionView  `json:"actions"`
}
type ActionView struct {
	ID        uint32 `json:"id"`
	Type      int32  `json:"type"`
	Params    string `json:"params"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at,omitempty"`
}
type TriggerView struct {
	URL       string `json:"url"`
	ID        uint32 `json:"id"`
	Type      int    `json:"type"`
	Params    string `json:"params"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

func toView(rule *pb.Rule) View {
	view := View{
		ID:     rule.Id,
		Host:   rule.Host,
		Status: rule.Status,
	}
	if rule.CreatedAt != nil {
		view.CreatedAt = util.TimestampToString(*rule.CreatedAt, time.RFC3339)
	}
	if rule.UpdatedAt != nil {
		view.UpdatedAt = util.TimestampToString(*rule.UpdatedAt, time.RFC3339)
	}

	view.Triggers = make([]TriggerView, 0, len(rule.Triggers))
	view.Actions = make([]ActionView, 0, len(rule.Actions))

	for _, t := range rule.Triggers {
		view.Triggers = append(view.Triggers, TriggerView{
			URL:       t.Url,
			ID:        t.Id,
			Type:      int(t.Type),
			Params:    t.Params,
			CreatedAt: util.TimestampToString(*t.CreatedAt, time.RFC3339),
			UpdatedAt: util.TimestampToString(*t.UpdatedAt, time.RFC3339),
		})
	}
	for _, a := range rule.Actions {
		view.Actions = append(view.Actions, ActionView{
			ID:        a.Id,
			Type:      int32(a.Type),
			Params:    a.Params,
			CreatedAt: util.TimestampToString(*a.CreatedAt, time.RFC3339),
			UpdatedAt: util.TimestampToString(*a.UpdatedAt, time.RFC3339),
		})
	}

	return view
}
