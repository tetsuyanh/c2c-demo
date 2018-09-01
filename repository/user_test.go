package repository

import (
	"fmt"

	"github.com/tetsuyanh/c2c-demo/model"
)

func CreateAuthenticatedUser() {
	u := model.DefaultUser()
	repo.Insert(u)

	a := model.DefaultAuthentication()
	a.UserId = u.Id
	a.EMail = fmt.Sprintf("%s@test.com", u.Id)
	a.Password = "hogehoge"
	a.Enabled = true
	repo.Insert(a)
}
