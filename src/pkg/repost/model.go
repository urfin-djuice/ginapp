package repost

import (
	"oko/pkg/account"
	"time"

	"github.com/jinzhu/gorm"
)

type Request struct {
	gorm.Model
	HasProcessed *bool `gorm:"column:has_processed;default:'null'"`
	Level        uint
	ParentID     uint `gorm:"column:parent_id;default:'null'"`
	URL          string
	Error        string             `gorm:"column:error;default:'null'"`
	Links        []Link             `gorm:"foreignkey:repost_id"`
	Accounts     []*account.Account `gorm:"many2many:account_repost_request"`
}

func (Request) TableName() string {
	return "repost_request"
}

type Link struct {
	ID          uint `gorm:"primary_key"`
	RepostID    uint `gorm:"column:repost_id"`
	URL         string
	PublishedAt time.Time `gorm:"column:published_at;default:'null'"`
	Title       string    `gorm:"column:title;default:'null'"`
}

func (Link) TableName() string {
	return "repost_link"
}

type RecordForExport struct {
	Title       string
	RepostURL   string
	ParentURL   string
	RepostLevel string
	PublishedAt *time.Time
}

func (er RecordForExport) fillCsvRec(rec []string) {
	rec[0] = er.Title
	rec[1] = er.RepostURL
	rec[2] = er.ParentURL
	rec[3] = er.RepostLevel
	if er.PublishedAt != nil {
		rec[4] = er.PublishedAt.String()
	} else {
		rec[4] = ""
	}
}
