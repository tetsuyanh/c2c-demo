package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tetsuyanh/c2c-demo/model"
)

func TestCreateUser(t *testing.T) {
	ur := GetUserRepo()
	ast := assert.New(t)

	u := model.NewUser()
	err := ur.CreateUser(u)
	ast.Nil(err)
}

func TestCreateSession(t *testing.T) {
	ur := GetUserRepo()
	ast := assert.New(t)

	u := model.NewUser()
	ur.CreateUser(u)

	s1 := model.NewSession(u.ID)
	err1 := ur.CreateSession(s1)
	ast.Nil(err1)

	// duplicate userID is ok for multi-login
	s2 := model.NewSession(u.ID)
	err2 := ur.CreateSession(s2)
	ast.Nil(err2)

	// invalid constraint of reference
	s3 := model.NewSession("--------")
	err3 := ur.CreateSession(s3)
	ast.NotNil(err3)

	clearAll()
}
