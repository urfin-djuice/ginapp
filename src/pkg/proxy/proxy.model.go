package proxy

import (
	"oko/pkg/ginapp/types"
	"time"
)

type CreateProxyRequest struct {
	Host  string `json:"host"`
	HTTPS bool   `json:"https"`
	HTTP  bool   `json:"http"`
}

type State struct {
	Host     string     `json:"host,omitempty"`
	StartBan *time.Time `json:"start_ban,omitempty"`
	EndBan   *time.Time `json:"end_ban,omitempty"`
}

type Response struct {
	//nolint
	types.StdResponse
	Host      string     `json:"host"`
	HTTP      bool       `json:"http"`
	HTTPS     bool       `json:"https"`
	ID        uint32     `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	State     []State    `json:"state"`
}
