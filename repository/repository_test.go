package repository

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tetsuyanh/c2c-demo/conf"
	"github.com/tetsuyanh/c2c-demo/model"
)

func TestMain(m *testing.M) {
	cp := conf.Postgres{
		DbUser:     "c2c_test",
		DbPassword: "",
		DbHost:     "localhost",
		DbName:     "c2c_test",
	}
	if eSetup := Setup(cp); eSetup != nil {
		log.Printf("failed to setup: %s", eSetup)
		os.Exit(1)
	}
	defer TearDown()

	clearAll()
	code := m.Run()
	os.Exit(code)
}

func TestUpdate(t *testing.T) {
	ast := assert.New(t)

	u := model.DefaultUser()
	err := repo.dbMap.Insert(u)
	ast.Nil(err)

	a1 := model.DefaultAsset()
	a1.UserId = u.Id
	a1.Point = 100
	err = repo.dbMap.Insert(a1) // Version is now 1
	ast.Nil(err)

	obj, err := repo.dbMap.Get(model.Asset{}, a1.Id)
	ast.NotNil(obj)
	ast.Nil(err)

	a2 := obj.(*model.Asset)
	a2.Point = 200
	_, err = repo.dbMap.Update(a2) // Version is now 2
	ast.Nil(err)

	a1.Point = 300
	count, err := repo.dbMap.Update(a1)
	ast.NotNil(err)
	ast.Equal(int64(-1), count)
}

func clearAll() {
	// order to delete for constraint
	repo.dbMap.Exec(clear(DealTableName))
	repo.dbMap.Exec(clear(ItemTableName))
	repo.dbMap.Exec(clear(AssetTableName))
	repo.dbMap.Exec(clear(AuthenticationTableName))
	repo.dbMap.Exec(clear(SessionTableName))
	repo.dbMap.Exec(clear(UserTableName))
}

func clear(m string) string {
	return fmt.Sprintf("delete from %s", m)
}
