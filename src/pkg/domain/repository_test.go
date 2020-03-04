package domain

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite
	mock sqlmock.Sqlmock
	db   *gorm.DB
	repo Repository
}

func (s *Suite) SetupSuite() {
	db, sqlMock, err := sqlmock.New()
	require.NoError(s.T(), err)

	s.db, err = gorm.Open("postgres", db)
	s.db = s.db.LogMode(true)
	require.NoError(s.T(), err)

	s.mock = sqlMock

	s.repo = NewDomainRepository(s.db)
}

func TestDomainRepository(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) TestList() {
	id := 1
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "domains" WHERE "domains"."deleted_at" IS NULL ORDER BY domains.updated_at`)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(id))
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "rss_links"  WHERE "rss_links"."deleted_at" IS NULL AND (("domain_id" IN ($1)))`)).
		WithArgs(id).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(id))
	_, err := s.repo.List()
	require.NoError(s.T(), err)
}

//nolint
func (s *Suite) TestListDomainsForTelegram() {
	id := 1
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "domains"  WHERE "domains"."deleted_at" IS NULL AND ((telegram_username is not null)) ORDER BY domains.updated_at`)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(id))
	_, err := s.repo.List(Filter{ForTelegram: true})
	require.NoError(s.T(), err)
}

//nolint
func (s *Suite) TestListDomainsWithRss() {
	id := 1
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT Distinct(domains.*) FROM "domains" join rss_links on domains.id = rss_links.id WHERE "domains"."deleted_at" IS NULL ORDER BY domains.updated_at`)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(id))
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "rss_links"  WHERE "rss_links"."deleted_at" IS NULL AND (("domain_id" IN ($1)))`)).
		WithArgs(id).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(id))
	_, err := s.repo.List(Filter{ExistRss: true})
	require.NoError(s.T(), err)
}

//nolint
func (s *Suite) TestListDomainsLast() {
	id := 1
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "domains" WHERE "domains"."deleted_at" IS NULL AND (((updated_at < $1 or updated_at is null))) ORDER BY domains.updated_at LIMIT 3`)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(id))
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "rss_links"  WHERE "rss_links"."deleted_at" IS NULL AND (("domain_id" IN ($1)))`)).
		WithArgs(id).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(id))
	_, err := s.repo.List(Filter{Last: true})
	require.NoError(s.T(), err)
}

func (s *Suite) AfterTest(_, _ string) {
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}
