package domain

import (
	"oko/pkg/rss"

	"github.com/jinzhu/gorm"
)

type Domain struct {
	gorm.Model
	Name             string
	TelegramUsername string `gorm:"column:telegram_username;default:'null'"`
	Cache            bool   `gorm:"column:cache;default:'false'"`
	Error            string `gorm:"column:error;default:'null'"`
	Rss              []*rss.Rss
}

func (Domain) TableName() string {
	return "domains"
}

type Filter struct {
	Last        bool
	ForTelegram bool
	ExistRss    bool
}

func (f Filter) preloadRss() bool {
	return f.Last || f.ExistRss
}
