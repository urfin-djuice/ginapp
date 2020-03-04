package links

import (
	"fmt"
	"oko/pkg/log"
	"strings"

	"github.com/jinzhu/gorm"
)

type Repository interface {
	List([]uint) (models []*Link, err error)
	Get(id uint) (*Link, error)
	Update(id uint, values *Link) error
	GetForDownloaderOld() []Link
	GetForCache(filter CacheFilter) ([]string, error)
	GetForContentParser() []Link
	BulkCreateRecords(links []Link) error
	Create(link *Link) error
}

type linkRepository struct {
	db *gorm.DB
}

func NewLinkRepository(db *gorm.DB) Repository {
	return &linkRepository{
		db: db,
	}
}

func (r *linkRepository) Create(link *Link) error {
	return r.db.Create(link).Error
}

func (r *linkRepository) Get(id uint) (model *Link, err error) {
	model = &Link{
		Model: gorm.Model{
			ID: id,
		},
	}

	if err = r.db.Preload("Domain").First(model).Error; err != nil && err != gorm.ErrRecordNotFound {
		log.Println("Error in LinkRepository.Get", err)
		return model, err
	}

	return
}

func (r *linkRepository) List(ids []uint) (models []*Link, err error) {
	if err = r.db.Preload("Domain").Where("id in (?)", ids).Find(&models).Error; err != nil {
		log.Println("Error in LinkRepository.List", err)
		return
	}

	return
}
func (r *linkRepository) Update(id uint, values *Link) error {
	err := r.db.
		Model(&Link{
			Model: gorm.Model{
				ID: id,
			},
		}).
		Updates(values).Error
	if err != nil {
		if err1 := r.Update(id, &Link{
			Error: err.Error(),
		}); err1 != nil {
			log.Println(err)
		}
		return err
	}
	return err
}

func (r *linkRepository) GetForDownloaderOld() []Link {
	var l []Link

	r.db.Table("links").
		Where("download_path is Null").
		Where("error is Null").
		Order("created_at asc").
		Limit(2000).
		Find(&l)

	return l
}

func (r *linkRepository) GetForCache(filter CacheFilter) (urls []string, err error) {
	q := r.db.Table("links").
		Select("url")
	if filter.SiteMapID != nil {
		q = q.Where("sitemap_id = ?")
	}
	if filter.DomainID != nil {
		q = q.Where("domain_id = ?")
	}

	rows, err := q.Rows()
	if err != nil {
		return
	}
	urls = make([]string, 0)
	for rows.Next() {
		var url string
		if err1 := rows.Scan(&url); err1 != nil {
			log.Println(err1)
			continue
		}
		urls = append(urls, url)
	}
	if err1 := rows.Close(); err1 != nil {
		log.Println(err1)
	}
	return
}

func (r *linkRepository) GetForContentParser() []Link {
	var l []Link
	r.db.Table("links").
		Where("links.has_content is False").
		Where("links.domain_id is not Null").
		Where("links.download_path is not Null").
		Where("links.error is Null").
		Joins("inner join domains on domains.id = links.domain_id and domains.cache = true").
		Order("created_at asc").
		Limit(1000).
		Find(&l)

	return l
}

func (r *linkRepository) BulkCreateRecords(links []Link) error {
	var valueStrings []string
	var valueArgs []interface{}

	for _, l := range links {
		valueStrings = append(valueStrings, "(?, ?, ?, ?)")

		valueArgs = append(valueArgs, l.URL)
		valueArgs = append(valueArgs, l.DomainID)
		valueArgs = append(valueArgs, l.SitemapID)
		valueArgs = append(valueArgs, l.PublishedAt)
	}

	smt := `INSERT INTO links(url, domain_id, sitemap_id, published_at) VALUES %s ON CONFLICT DO NOTHING`

	smt = fmt.Sprintf(smt, strings.Join(valueStrings, ","))
	tx := r.db.Begin()
	if err := tx.Exec(smt, valueArgs...).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}
