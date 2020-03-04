package domain

import (
	"errors"
	"oko/pkg/log"
	"time"

	"github.com/jinzhu/gorm"
)

type Repository interface {
	List(filter ...Filter) (models []*Domain, err error)
	Get(id uint) (*Domain, error)
	Create(model *Domain) error
	Update(id uint, model *Domain) error
	Delete(model *Domain) error
	GetForCacheJob(limit int) []Domain
	GetByName(name string) (*Domain, bool)
}

type domainRepository struct {
	db *gorm.DB
}

func NewDomainRepository(db *gorm.DB) Repository {
	return &domainRepository{
		db: db,
	}
}

func (r *domainRepository) Get(id uint) (model *Domain, err error) {
	model = &Domain{
		Model: gorm.Model{
			ID: id,
		},
	}

	if err = r.db.First(model).Error; err != nil && err != gorm.ErrRecordNotFound {
		log.Println("Error in DomainRepository.Get", err)
		return model, err
	}
	return
}

func (r *domainRepository) Create(model *Domain) error {
	if err := r.db.Create(model).Error; err != nil {
		log.Println("Error in DomainRepository.Create", err)
		return err
	}

	return nil
}

func (r *domainRepository) Update(id uint, model *Domain) error {
	domain := &Domain{
		Model: gorm.Model{
			ID: id,
		},
	}

	result := r.db.Model(domain).Updates(model)
	if err := result.Error; err != nil {
		log.Println("Error in DomainRepository.Update", err)
		return err
	}
	if rowsAffected := result.RowsAffected; rowsAffected == 0 {
		log.Printf("Error in DomainRepository.Update, rowsAffected: %v", rowsAffected)
		return errors.New("no records updated, No match was found")
	}
	return nil
}

func (r *domainRepository) Delete(model *Domain) error {
	result := r.db.Delete(model)
	if err := result.Error; err != nil {
		log.Println("Error in DomainRepository.Delete", err)
		return err
	}

	if rowsAffected := result.RowsAffected; rowsAffected == 0 {
		log.Printf("Error in DomainRepository.Delete, rowsAffected: %v", rowsAffected)
		return errors.New("no records deleted, No match was found")
	}
	return nil
}

func (r *domainRepository) List(f ...Filter) (models []*Domain, err error) {
	q := r.db.Model(&Domain{})

	if len(f) > 0 {
		if f[0].ForTelegram {
			q = q.Where("telegram_username is not null")
		}
		if f[0].ExistRss {
			q = q.Joins("join rss_links on domains.id = rss_links.id").
				Select("Distinct(domains.*)")
		}
		if f[0].Last {
			q = q.Where("(updated_at < ? or updated_at is null)", time.Now().AddDate(0, 0, -1)).
				Limit(3)
		}
	}

	if len(f) == 0 || f[0].preloadRss() {
		q = q.Preload("Rss")
	}

	if err = q.Order("domains.updated_at").Find(&models).Error; err != nil {
		log.Println("Error in DomainRepository.List", err)
		return
	}

	return
}

func (r domainRepository) GetForCacheJob(limit int) []Domain {
	var domains []Domain

	r.db.Table("domains").
		Where("domains.cache = ?", "false").
		Joins("inner join links on links.domain_id = domains.id and links.download_path is not null").
		Group("domains.id").
		Having("count(links.id) > ?", limit).
		Order("updated_at asc").
		Find(&domains)

	return domains
}

func (r domainRepository) GetByName(name string) (*Domain, bool) {
	dom := &Domain{}
	notFound := r.db.Model(dom).Where("name = ?", name).First(dom).RecordNotFound()
	return dom, notFound
}
