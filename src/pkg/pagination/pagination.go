package pagination

import (
	"math"

	"github.com/jinzhu/gorm"
)

type Param struct {
	DB          *gorm.DB
	CurrentPage int
	PerPage     int
	OrderBy     []string
}

type Paginator struct {
	TotalRecords int `json:"total_records"`
	TotalPages   int `json:"total_pages"`
	PerPage      int `json:"per_page"`
	CurrentPage  int `json:"current_page"`
}

func Paging(p *Param, result interface{}) (paginator Paginator, err error) {
	db := p.DB

	if p.CurrentPage < 1 {
		p.CurrentPage = 1
	}
	if p.PerPage == 0 {
		p.PerPage = 10
	}
	if len(p.OrderBy) > 0 {
		for _, o := range p.OrderBy {
			db = db.Order(o)
		}
	}

	done := make(chan bool, 1)
	var count, offset int

	go countRecords(db, result, done, &count)

	if p.CurrentPage == 1 {
		offset = 0
	} else {
		offset = (p.CurrentPage - 1) * p.PerPage
	}

	err = db.Limit(p.PerPage).Offset(offset).Find(result).Error
	if err != nil {
		return
	}
	<-done

	paginator.TotalRecords = count
	paginator.CurrentPage = p.CurrentPage

	paginator.PerPage = p.PerPage
	paginator.TotalPages = int(math.Ceil(float64(count) / float64(p.PerPage)))

	return
}

func countRecords(db *gorm.DB, anyType interface{}, done chan bool, count *int) {
	db.Model(anyType).Count(count)
	done <- true
}
