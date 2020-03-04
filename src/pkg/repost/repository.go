package repost

import (
	"errors"
	"oko/pkg/account"
	"oko/pkg/log"
	"oko/pkg/util"
	"time"

	"github.com/jinzhu/gorm"
)

type Repository interface {
	List(limit uint32, page uint32, maxLevel uint32, accountID int, hasProcessed string) (models []*Request,
		count uint32, err error)
	ListForParser(limit, page, maxLevel uint32, hasProcessed string) (models []*Request, count uint32, err error)
	Get(model *Request) (*Request, error)
	Create(model *Request) error
	Exist(model *Request) (bool, error)
	Update(id uint, model interface{}) error
	GetWithDate(url string, accountID int, dateFrom, dateTo *time.Time) (*Request, error)
	GetForExport(url string, accID int, from, to *time.Time) ([]RecordForExport, error)
	GetOrNil(*Request) (*Request, error)
	CreateAndAssign(*Request, account.Account) error
}

const trueStr = "true"
const falseStr = "false"
const nullStr = "null"

type requestRepository struct {
	db *gorm.DB
}

func NewRequestRepository(db *gorm.DB) Repository {
	return &requestRepository{
		db: db,
	}
}

func (repo *requestRepository) Exist(model *Request) (ok bool, err error) {
	err = repo.db.Where(model).First(model).Error
	if gorm.IsRecordNotFoundError(err) {
		return false, nil
	}

	if err != nil {
		log.Println("Error in RequestRepository.Exist", err)
		return false, err
	}

	return true, err
}

func (repo *requestRepository) GetOrNil(model *Request) (res *Request, err error) {
	err = repo.db.
		Preload("Links").
		Joins("join repost_link on repost_link.repost_id = repost_request.id").
		Preload("Accounts").
		Where(model).First(model).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
	}
	return model, err
}

func (repo *requestRepository) Get(model *Request) (*Request, error) {
	var err error

	if err = repo.db.Preload("Links").First(model).Error; err != nil {
		log.Println("Error in RequestRepository.Get", err)
		return model, err
	}

	return model, err
}
func (repo *requestRepository) ListForParser(limit, page, maxLevel uint32, hasProcessed string) (
	models []*Request, count uint32, err error) {
	var offset uint32
	if page > 1 {
		offset = (page - 1) * limit
	} else {
		offset = 0
	}

	query := repo.db.
		Preload("Links").
		Order("created_at asc").
		Where("level <= ?", maxLevel).
		Offset(offset).
		Limit(limit)

	switch hasProcessed {
	case nullStr:
		query = query.Where("has_processed is null")
	case trueStr, falseStr:
		query = query.Where("has_processed is ?", hasProcessed)
	}

	repo.db.Model(&Request{}).Where("level <= ?", maxLevel).Count(&count)

	err = query.Find(&models).Error

	if err != nil {
		log.Println("Error in RequestRepository.List", err)
		return
	}

	return
}
func (repo *requestRepository) List(limit, page, maxLevel uint32, accountID int, hasProcessed string) (
	models []*Request, count uint32, err error) {
	var offset uint32
	if page > 1 {
		offset = (page - 1) * limit
	} else {
		offset = 0
	}

	query := repo.db.
		Preload("Links").
		Preload("Accounts").
		Order("created_at asc").
		Where("level <= ?", maxLevel).
		Joins("join account_repost_request arr on arr.request_id = repost_request.id").
		Where("arr.account_id = ?", accountID).
		Offset(offset).
		Limit(limit)

	switch hasProcessed {
	case nullStr:
		query = query.Where("has_processed is null")
	case trueStr, falseStr:
		query = query.Where("has_processed is ?", hasProcessed)
	}

	repo.db.Model(&Request{}).Where("level <= ?", maxLevel).Count(&count)

	err = query.Find(&models).Error

	if err != nil {
		log.Println("Error in RequestRepository.List", err)
		return
	}

	return
}

func (repo *requestRepository) Create(model *Request) error {
	if err := repo.db.Create(model).Error; err != nil {
		log.Println("Error in RequestRepository.Create", err)
		return err
	}

	return nil
}

func (repo *requestRepository) Update(id uint, model interface{}) error {
	result := repo.db.
		Table(Request{}.TableName()).
		Where("id = ?", id).
		Updates(model)

	if err := result.Error; err != nil {
		log.Println("Error in RequestRepository.Update", err)
		return err
	}

	if rowsAffected := result.RowsAffected; rowsAffected == 0 {
		log.Printf("Error in DomainRepository.Update, rowsAffected: %v", rowsAffected)
		return errors.New("no records updated, No match was found")
	}
	return nil
}

func (repo *requestRepository) GetWithDate(url string, accountID int, dateFrom, dateTo *time.Time) (*Request, error) {
	var model Request
	var err error
	if !repo.haveAccess(url, accountID) {
		log.Printf("Error in RequestRepository.GetWithDate access denied to %s for %d", url, accountID)
		return nil, errors.New("not found")
	}
	preload := repo.db.Preload("Links").Where("(repost_request.url = ? or repost_request.url = ?)", url,
		util.URLEncoded(url))
	if dateFrom != nil || dateTo != nil {
		preload = preload.Joins("join repost_link on repost_link.repost_id = repost_request.id")
		if dateFrom != nil {
			preload = preload.Where("repost_link.published_at >= ?", dateFrom)
		}
		if dateTo != nil {
			preload = preload.Where("repost_link.published_at <= ?", dateTo)
		}
	}
	err = preload.
		First(&model).
		Error

	if err != nil {
		log.Println("Error in RequestRepository.Get", err)
		return nil, err
	}

	return &model, nil
}

func (repo *requestRepository) haveAccess(url string, accID int) bool {
	row := repo.db.Raw(`WITH RECURSIVE
    starting (id, url, parent_id) AS
        (
            SELECT t.id, t.url, t.parent_id
            FROM repost_request AS t
            WHERE t.url = ? or t.url = ? 
        ),
    descendants (id, url, parent_id) AS
        (
            SELECT s.id, s.url, s.parent_id
            FROM starting AS s
            UNION ALL
            SELECT t.id, t.url, t.parent_id
            FROM repost_request AS t
                     JOIN descendants AS d ON t.parent_id = d.id
        ),
    ancestors (id, url, parent_id) AS
        (
            SELECT t.id, t.url, t.parent_id
            FROM repost_request AS t
            WHERE t.id IN (SELECT parent_id FROM starting)
            UNION ALL
            SELECT t.id, t.url, t.parent_id
            FROM repost_request AS t
                     JOIN ancestors AS a ON t.id = a.parent_id
        )
select distinct account_id
from (TABLE ancestors
      UNION ALL
      TABLE descendants) as res
         join account_repost_request on account_repost_request.request_id = res.id and account_id = ?`, url,
		util.URLEncoded(url), accID).Row()
	tmp := -1
	if err := row.Scan(&tmp); err != nil {
		return false
	}
	return true
}

func (repo *requestRepository) GetForExport(url string, accID int, from, to *time.Time) ([]RecordForExport, error) {
	args := make([]interface{}, 0, 3)
	args = append(args, url)
	if !repo.haveAccess(url, accID) {
		log.Printf("Error in RequestRepository.GetForExport access denied to %s for %d", url, accID)
		return nil, errors.New("access denied")
	}
	sql := `WITH RECURSIVE nodes(id, parent_id) AS (
    SELECT s1.id, s1.parent_id
    FROM repost_request s1
    WHERE url = ?
    UNION
    SELECT s2.id, s2.parent_id
    FROM repost_request s2,
         nodes s1
    WHERE s2.parent_id = s1.id
)
SELECT
    r1.title        as title,
    r1.url          as repost_url,
    rp.url          as parent_url,
    rr.level        as repost_level,
    r1.published_at as published_at
FROM nodes n
         join repost_link r1 on n.id = r1.repost_id
         join repost_request rr on n.id = rr.id
         join repost_request rp on rp.id = n.parent_id `
	if from != nil || to != nil {
		sql += "where published_at is not null "
		if from != nil {
			sql += "and published_at > ? "
			args = append(args, from)
		}
		if to != nil {
			sql += "and published_at < ?"
			args = append(args, to)
		}
	}
	raw := repo.db.Raw(sql, args...)

	exportRecs := make([]RecordForExport, 0, 100)
	err := raw.Scan(&exportRecs).Error
	return exportRecs, err
}
func (repo *requestRepository) CreateAndAssign(req *Request, acc account.Account) error {
	tx := repo.db.Begin()
	if err := tx.Error; err != nil {
		return err
	}
	if req.ID == 0 {
		if err := tx.Create(req).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	if err := tx.Model(req).Association("Accounts").Append(acc).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
