package repo

import (
	"database/sql"
	"fmt"

	"github.com/go-gorp/gorp"
	_ "github.com/lib/pq"

	"github.com/tetsuyanh/c2c-demo/conf"
)

var (
	repo *Repo
)

type Repo struct {
	db    *sql.DB
	dbMap *gorp.DbMap
	// models    map[string]interface{}
}

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
	// models[UserRepoName] = model.User{}
	dbMap := &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}
	for k, v := range models {
		dbMap.AddTableWithName(v, k).SetKeys(false, "id")
	}

	repo = &Repo{
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
