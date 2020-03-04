package rss

import (
	"errors"
	"oko/pkg/log"

	"github.com/jinzhu/gorm"
)

type rssRepository struct {
	db *gorm.DB
}

func NewRssRepository(db *gorm.DB) *rssRepository { //nolint
	return &rssRepository{
		db: db,
	}
}

func (r rssRepository) CreateRss(domainID uint, link string) (err error) {
	if err = r.db.Create(&Rss{
		Link:     link,
		DomainID: domainID,
	}).Error; err != nil {
		log.Error("Error in RssRepository.CreateRss", err)
	}
	return
}

func (r rssRepository) UpdateRss(domainID uint, link string, rssID uint) (err error) {
	rss := &Rss{
		Model: gorm.Model{
			ID: rssID,
		},
	}
	model := &Rss{
		DomainID: domainID,
		Link:     link,
	}
	res := r.db.First(rss).Updates(model)
	if err = res.Error; err != nil {
		log.Print("Error in RssRepository.UpdateRss", err)
		return
	}

	if res.RowsAffected == 0 {
		log.Printf("Error in RssRepository.Update, rowsAffected: %v", res.RowsAffected)
		err = errors.New("no records updated, No match was found")
	}
	return
}

func (r rssRepository) DeleteRss(rssID uint) (err error) {
	if err = r.db.Delete(&Rss{
		Model: gorm.Model{
			ID: rssID,
		},
	}).Error; err != nil {
		log.Print("Error in RssRepository.DeleteRss", err)
	}
	return err
}

func (r rssRepository) List() (res []*Rss, err error) {
	res = []*Rss{}
	err = r.db.Model(&res).Find(&res).Error
	if err != nil {
		log.Print("Error in RssRepository.List", err)
	}
	return
}

func (r rssRepository) Get(id uint) (res *Rss, err error) {
	res = &Rss{}
	err = r.db.Find(res, "id = ?", id).Error
	if err != nil {
		log.Print("Error in RssRepository.Get", err)
	}
	return
}
