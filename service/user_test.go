package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetUser(t *testing.T) {
	ast := assert.New(t)
	userSrv := GetUserService()

	// invalid token
	{
		u, e := userSrv.GetUser("hogehogeToken")
		ast.Nil(u)
		ast.NotNil(e)
	}

	// success
	{
		an := createAnonymousUser()
		s := createSession(an.Id)
		u, e := userSrv.GetUser(*s.Token)
		ast.NotNil(u)
		ast.Nil(e)
	}
}

func TestGetAsset(t *testing.T) {
	ast := assert.New(t)
	userSrv := GetUserService()

	// invalid userId
	{
		as, e := userSrv.GetAsset("hogehogeId")
		ast.Nil(as)
		ast.NotNil(e)
	}

	// success
	{
		pu, _, _, _ := createPerfectUser()
		as, e := userSrv.GetAsset(pu.Id)
		ast.NotNil(as)
		ast.Nil(e)
	}
}

func TestStart(t *testing.T) {
	ast := assert.New(t)
	userSrv := GetUserService()

	// success
	{
		s, e := userSrv.Start()
		ast.NotNil(s)
		ast.Nil(e)
	}
}

func TestLogin(t *testing.T) {
	ast := assert.New(t)
	userSrv := GetUserService()

	// invalid email
	{
		s, e := userSrv.Login("hoge@hoge.com", testAuthPass)
		ast.Nil(s)
		ast.NotNil(e)
	}

	// not enabled
	{
		an := createAnonymousUser()
		au := createAuthentication(an, false)
		s, e := userSrv.Login(*au.EMail, testAuthPass)
		ast.Nil(s)
		ast.NotNil(e)
	}

	// invalid password
	{
		_, au, _, _ := createPerfectUser()
		s, e := userSrv.Login(*au.EMail, "hogehoge")
		ast.Nil(s)
		ast.NotNil(e)
	}

	// success
	{
		_, au, _, _ := createPerfectUser()
		s, e := userSrv.Login(*au.EMail, testAuthPass)
		ast.NotNil(s)
		ast.Nil(e)
	}
}

func TestPublishAuth(t *testing.T) {
	ast := assert.New(t)
	userSrv := GetUserService()

	// already enabled
	{
		pu, au, _, _ := createPerfectUser()
		a, e := userSrv.PublishAuth(pu.Id, *au.EMail, testAuthPass)
		ast.Nil(a)
		ast.NotNil(e)
	}

	// success, republish when auth is not enable
	{
		an := createAnonymousUser()
		au := createAuthentication(an, false)
		s, e := userSrv.PublishAuth(an.Id, *au.EMail, testAuthPass)
		ast.NotNil(s)
		ast.Nil(e)
	}

	// success, publish new auth
	{
		an := createAnonymousUser()
		mail := createEmail(an.Id)
		s, e := userSrv.PublishAuth(an.Id, mail, testAuthPass)
		ast.NotNil(s)
		ast.Nil(e)
	}
}

func TestEnableAuth(t *testing.T) {
	ast := assert.New(t)
	userSrv := GetUserService()

	// invalid token
	{
		e := userSrv.EnableAuth("hogehogeToken")
		ast.NotNil(e)
	}

	// already enabled
	{
		_, au, _, _ := createPerfectUser()
		e := userSrv.EnableAuth(*au.Token)
		ast.NotNil(e)
	}

	// success
	{
		an := createAnonymousUser()
		au := createAuthentication(an, false)
		e := userSrv.EnableAuth(*au.Token)
		ast.Nil(e)
	}
}

func TestEncryptAndCorrectPasswor(t *testing.T) {
	ast := assert.New(t)

	srcs := []string{
		"hogehoge",
		"fugafuga",
		"1111",
	}
	for _, src := range srcs {
		ast.True(correctPassword(encrypt(src), src))
	}
}
