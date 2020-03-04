package links

import (
	"oko/pkg/author"
	"oko/pkg/domain"
	"time"

	"github.com/jinzhu/gorm"
)

type Link struct {
	gorm.Model
	URL          string
	PublishedAt  *time.Time        `gorm:"column:published_at;default:'null'"`
	CreatedAt    *time.Time        `gorm:"column:created_at;default:'null'"`
	DomainID     uint              `gorm:"column:domain_id"`
	DownloadPath string            `gorm:"column:download_path;default:'null'"`
	HasContent   bool              `gorm:"column:has_content;default:'false'"`
	HasIndex     bool              `gorm:"column:has_index;default:'false'"`
	Error        string            `gorm:"column:error;default:'null'"`
	SitemapID    *uint             `gorm:"column:sitemap_id;default:'null'"`
	Authors      *[]*author.Author `gorm:"many2many:links_author"`

	SentimentalScore    *float32 `gorm:"column:sentimental_score;default:'null'"`
	SentimentalPositive *float32 `gorm:"column:sentimental_positive;default:'null'"`
	SentimentalNegative *float32 `gorm:"column:sentimental_negative;default:'null'"`

	Domain domain.Domain
}

func (Link) TableName() string {
	return "links"
}

type CacheFilter struct {
	SiteMapID *uint
	DomainID  *uint
}
