package service

import (
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/tetsuyanh/c2c-demo/model"
	"github.com/tetsuyanh/c2c-demo/repository"
)

var (
	userService UserService

	initialPoint = 10000
)

type (
	UserService interface {
		GetUser(token string) (*model.User, error)
		GetAsset(userId string) (*model.Asset, error)
		Start() (*model.Session, error)
		Login(email, password string) (*model.Session, error)
		PublishAuth(userID, email, password string) (*model.Authentication, error)
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

func (us *userServiceImpl) GetUser(token string) (*model.User, error) {
	return us.userRepo.FindUserBySettionToken(token)
}

func (us *userServiceImpl) GetAsset(userId string) (*model.Asset, error) {
	return us.userRepo.FindAssetByUserId(userId)
}

func (us *userServiceImpl) Start() (*model.Session, error) {
	u := model.DefaultUser()
	s := model.DefaultSession()
	s.UserId = &u.Id

	tx := us.repo.Transaction()
	if err := tx.Insert(u).Insert(s).Commit(); err != nil {
		return nil, err
	}
	return s, nil
}

func (us *userServiceImpl) Login(email, password string) (*model.Session, error) {
	a, err := us.userRepo.FindAuthByEmail(email)
	if err != nil {
		return nil, err
	}
	if !*a.Enabled {
		return nil, fmt.Errorf("not enabled")
	}
	if !correctPassword(*a.Password, password) {
		return nil, fmt.Errorf("invalid password")
	}

	s := model.DefaultSession()
	s.UserId = a.UserId
	if err := us.repo.Insert(s); err != nil {
		return nil, err
	}
	return s, nil
}

func (us *userServiceImpl) PublishAuth(userID, email, password string) (*model.Authentication, error) {
	tx := us.repo.Transaction()

	exist, _ := us.userRepo.FindAuthByEmail(email)
	if exist != nil {
		if *exist.Enabled {
			return nil, fmt.Errorf("email already enabled")
		}
		tx.Delete(exist)
	}

	encrypted := encrypt(password)
	a := model.DefaultAuthentication()
	a.UserId = &userID
	a.EMail = &email
	a.Password = &encrypted
	if err := tx.Insert(a).Commit(); err != nil {
		return nil, err
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
