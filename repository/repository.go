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
		Exec(q string) error
		Transaction() Transaction
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
	models[AssetTableName] = model.Asset{}
	models[SessionTableName] = model.Session{}
	models[AuthenticationTableName] = model.Authentication{}
	models[ItemTableName] = model.Item{}
	models[DealTableName] = model.Deal{}
	dbMap := &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}
	for k, v := range models {
		t := dbMap.AddTableWithName(v, k)
		t.SetKeys(false, "id")
		// table using auto locking by gorp have to call this function
		// see https://github.com/go-gorp/gorp#optimistic-locking
		if k == AssetTableName || k == ItemTableName {
			t.SetVersionCol("version")
		}
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
	obj, err := r.dbMap.Get(i, id)
	if err != nil {
		return nil, err
	} else if obj == nil && err == nil {
		// when row is not found, gorp returns nil, nil
		// see https://github.com/go-gorp/gorp/blob/master/gorp.go#L428
		return nil, fmt.Errorf("not found")
	}
	return obj, err
}

func (r *repoImpl) Insert(i interface{}) error {
	return r.dbMap.Insert(i)
}

func (r *repoImpl) Update(i interface{}) error {
	if cnt, err := r.dbMap.Update(i); err != nil {
		return fmt.Errorf("invalid query to update: %s", err)
	} else if cnt == 0 {
		return fmt.Errorf("no target to update")
	}
	return nil
}

func (r *repoImpl) Delete(i interface{}) error {
	if cnt, err := r.dbMap.Delete(i); err != nil {
		return fmt.Errorf("invalid query to delete: %s", err)
	} else if cnt == 0 {
		return fmt.Errorf("no target to delete")
	}
	return nil
}

// Exec is to exec
func (r *repoImpl) Exec(q string) error {
	_, err := r.dbMap.Exec(q)
	if err != nil {
		return err
	}
	return nil
}

func (r *repoImpl) Transaction() Transaction {
	tx, err := r.dbMap.Begin()
	return &transactionImpl{
		tx:  tx,
		err: err,
	}
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
