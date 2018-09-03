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

	testAuthPass  = "password"
	testItemPrice = 100
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

func createSession(userId string) *model.Session {
	s := model.DefaultSession()
	s.UserId = userId
	repo.Insert(s)
	return s
}

func createPerfectUser() (*model.User, *model.Authentication, *model.Asset, *model.Item) {
	u := createAnonymousUser()
	au := createAuthentication(u.Id, true)
	as := createAsset(u.Id, initialPoint)
	i := createItem(u.Id, model.ItemStatusSale)
	return u, au, as, i
}

func createEmail(userId string) string {
	return fmt.Sprintf("%s@test.com", userId)
}

func createAuthentication(userId string, enable bool) *model.Authentication {
	au := model.DefaultAuthentication()
	email := createEmail(userId)
	encrypted := encrypt(testAuthPass)
	au.UserId = userId
	au.EMail = email
	au.Password = encrypted
	au.Enabled = enable
	repo.Insert(au)
	return au
}

func createAsset(userId string, point int) *model.Asset {
	as := model.DefaultAsset()
	as.UserId = userId
	as.Point = point
	repo.Insert(as)
	return as
}

func createItem(userId string, status string) *model.Item {
	i := model.DefaultItem()
	label := "label"
	price := testItemPrice
	i.UserId = userId
	i.Label = label
	i.Price = price
	i.Status = status
	repo.Insert(i)
	return i
}

func createDeal(sellerId, buyerId string) *model.Deal {
	i := createItem(sellerId, model.ItemStatusSale)
	repo.Insert(i)
	d := model.DefaultDeal()
	d.ItemId = i.Id
	d.SellerId = sellerId
	d.BuyerId = buyerId
	repo.Insert(d)
	return d
}
