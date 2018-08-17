package repository

import (
	"database/sql"
	"fmt"

	"github.com/go-gorp/gorp"
	_ "github.com/lib/pq"

	"github.com/tetsuyanh/c2c-demo/conf"
	"github.com/tetsuyanh/c2c-demo/model"
)

var (
	repo *repoImpl
)

type (
	Repo interface {
		Get(i interface{}, id string) (interface{}, error)
		Insert(i interface{}) error
		Update(i interface{}) error
		Delete(i interface{}) error
	}

	repoImpl struct {
		db    *sql.DB
		dbMap *gorp.DbMap
	}
)

func Setup(c conf.Postgres) error {
	if repo != nil {
		return nil
	}

	// sslmode disable, because of internal connect
	connStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", c.DbUser, c.DbPassword, c.DbHost, c.DbName)
	db, errOpen := sql.Open("postgres", connStr)
	if errOpen != nil {
		return errOpen
	}
	if errPing := db.Ping(); errPing != nil {
		return errPing
	}

	models := make(map[string]interface{}, 64)
	models[UserTableName] = model.User{}
	models[SessionTableName] = model.Session{}
	models[AuthenticationTableName] = model.Authentication{}
	models[ItemTableName] = model.Item{}
	models[DealTableName] = model.Deal{}
	dbMap := &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}
	for k, v := range models {
		dbMap.AddTableWithName(v, k).SetKeys(false, "id")
	}

	repo = &repoImpl{
		db:    db,
		dbMap: dbMap,
	}
	return nil
}

func TearDown() {
	if repo != nil {
		repo.dbMap = nil
		repo.db.Close()
		repo.db = nil
		repo = nil
	}
}

func GetRepository() Repo {
	return repo
}

func (r *repoImpl) Get(i interface{}, id string) (interface{}, error) {
	return r.dbMap.Get(i, id)
}

func (r *repoImpl) Insert(i interface{}) error {
	return r.dbMap.Insert(i)
}

func (r *repoImpl) Update(i interface{}) error {
	if cnt, err := r.dbMap.Update(i); err != nil {
		return fmt.Errorf("Update invalid query: %s", err)
	} else if cnt == 0 {
		return fmt.Errorf("Update no target")
	}
	return nil
}

func (r *repoImpl) Delete(i interface{}) error {
	if cnt, err := r.dbMap.Delete(i); err != nil {
		return fmt.Errorf("delete invalid query: %s", err)
	} else if cnt == 0 {
		return fmt.Errorf("delete no target")
	}
	return nil
}

func resultExec(result sql.Result, err error) error {
	if err != nil {
		return fmt.Errorf("query error: %s", err)
	}
	// see: https://golang.org/pkg/database/sql/#Result
	// Not every database or database driver may support this.
	if cnt, errRow := result.RowsAffected(); errRow != nil {
		return fmt.Errorf("result.RowsAffected: %s", errRow)
	} else if cnt == 0 {
		return fmt.Errorf("query no target")
	}
	return nil
}
