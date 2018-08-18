package service

import (
	"fmt"
	"time"

	"github.com/tetsuyanh/c2c-demo/model"
	"github.com/tetsuyanh/c2c-demo/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	userService UserService

	initialPoint = 10000
)

type (
	UserService interface {
		CreateUserSession() (*model.Session, error)
		CreateSession(userID string) (*model.Session, error)
		CreateAuth(userID, email, password string) (*model.Authentication, error)
		GetUser(token string) (*model.User, error)
		GetAuth(email, password string) (*model.Authentication, error)
		EnableAuth(token string) error
	}

	userServiceImpl struct {
		repo     repository.Repo
		userRepo repository.UserRepo
	}
)

func GetUserService() UserService {
	if userService == nil {
		userService = &userServiceImpl{
			repo:     repository.GetRepository(),
			userRepo: repository.GetUserRepo(),
		}
	}
	return userService
}

func (us *userServiceImpl) CreateUserSession() (*model.Session, error) {
	u := model.DefaultUser()
	s := model.DefaultSession()
	s.UserId = &u.Id

	tx := us.repo.Transaction()
	if err := tx.Insert(u).Insert(s).Commit(); err != nil {
		return nil, err
	}
	return s, nil
}

func (us *userServiceImpl) CreateSession(userId string) (*model.Session, error) {
	s := model.DefaultSession()
	s.UserId = &userId
	if err := us.repo.Insert(s); err != nil {
		return nil, err
	}
	return s, nil
}

func (us *userServiceImpl) CreateAuth(userID, email, password string) (*model.Authentication, error) {
	encrypted := encrypt(password)
	a := model.DefaultAuthentication()
	a.UserId = &userID
	a.EMail = &email
	a.Password = &encrypted
	if err := us.repo.Insert(a); err != nil {
		return nil, err
	}
	return a, nil
}

func (us *userServiceImpl) GetUser(token string) (*model.User, error) {
	return us.userRepo.FindUserBySettionToken(token)
}

func (us *userServiceImpl) GetAuth(email, password string) (*model.Authentication, error) {
	a, errFind := us.userRepo.FindAuthByEmail(email)
	if errFind != nil {
		return nil, errFind
	}
	if !correctPassword(*a.Password, password) {
		return nil, fmt.Errorf("invalid password")
	}
	if !*a.Enabled {
		return nil, fmt.Errorf("not enabled")
	}
	return a, nil
}

func (us *userServiceImpl) EnableAuth(token string) error {
	au, err := us.userRepo.FindAuthByToken(token)
	if err != nil {
		return err
	}
	if *au.Enabled {
		return fmt.Errorf("already enabled")
	}
	tr := true
	t := time.Now()
	au.Enabled = &tr
	au.UpdatedAt = &t

	as := model.DefaultAsset()
	as.UserId = au.UserId
	as.Point = &initialPoint

	tx := us.repo.Transaction()
	return tx.Update(au).Insert(as).Commit()
}

func encrypt(password string) string {
	// bcrypt.MinCost = 4, bcrypt.MaxCost = 31, bcrypt.DefaultCost = 10
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash)
}

func correctPassword(encryped, raw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(encryped), []byte(raw)) == nil
}
