package rss

import (
	"github.com/jinzhu/gorm"
)

type Rss struct {
	gorm.Model
	Link     string
	DomainID uint `gorm:"foreignkey:DomainID"`
}

func (r *Rss) TableName() string {
	return "rss_links"
}
