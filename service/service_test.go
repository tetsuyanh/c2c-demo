package service

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/tetsuyanh/c2c-demo/conf"
	"github.com/tetsuyanh/c2c-demo/model"
	"github.com/tetsuyanh/c2c-demo/repository"
)

var (
	repo repository.Repo
)

func TestMain(m *testing.M) {
	cp := conf.Postgres{
		DbUser:     "c2c_test",
		DbPassword: "",
		DbHost:     "localhost",
		DbName:     "c2c_test",
	}
	if eSetup := repository.Setup(cp); eSetup != nil {
		log.Printf("failed to setup: %s", eSetup)
		os.Exit(1)
	}
	defer repository.TearDown()

	repo = repository.GetRepository()
	clearAll(repo)
	code := m.Run()
	os.Exit(code)
}

func clearAll(repo repository.Repo) {
	// order to delete for constraint
	repo.Exec(clear(repository.DealTableName))
	repo.Exec(clear(repository.ItemTableName))
	repo.Exec(clear(repository.AssetTableName))
	repo.Exec(clear(repository.AuthenticationTableName))
	repo.Exec(clear(repository.SessionTableName))
	repo.Exec(clear(repository.UserTableName))
}

func clear(m string) string {
	return fmt.Sprintf("delete from %s", m)
}

func createAnonymousUser() *model.User {
	u := model.DefaultUser()
	repo.Insert(u)
	return u
}

func createPerfectUser() (*model.User, *model.Authentication, *model.Asset, *model.Item) {
	u := createAnonymousUser()
	au := createAuthentication(u)
	as := createAsset(u)
	i := createItem(u, model.ItemStatusSold)
	return u, au, as, i
}

func createAuthentication(u *model.User) *model.Authentication {
	au := model.DefaultAuthentication()
	email := fmt.Sprintf("%s@test.com", u.Id)
	pass := "hogehoge"
	t := true
	au.UserId = &u.Id
	au.EMail = &email
	au.Password = &pass
	au.Enabled = &t
	repo.Insert(au)
	return au
}

func createAsset(u *model.User) *model.Asset {
	as := model.DefaultAsset()
	as.UserId = &u.Id
	as.Point = &initialPoint
	repo.Insert(as)
	return as
}

func createItem(u *model.User, status string) *model.Item {
	i := model.DefaultItem()
	label := "label"
	price := 100
	i.UserId = &u.Id
	i.Label = &label
	i.Price = &price
	i.Status = &status
	repo.Insert(i)
	return i
}
