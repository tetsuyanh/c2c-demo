package repository

import (
	"github.com/tetsuyanh/c2c-demo/model"
)

const (
	UserTableName           = "users"
	SessionTableName        = "sessions"
	AuthenticationTableName = "authentications"
)

var (
	userRepo UserRepo
)

type (
	UserRepo interface {
		FindUserBySettionToken(token string) (*model.User, error)
		FindAuthByUserId(userId string) (*model.Authentication, error)
		FindAuthByEmail(email string) (*model.Authentication, error)
		FindAuthByToken(token string) (*model.Authentication, error)
		FindAssetByUserId(userId string) (*model.Asset, error)
	}

	userRepoImpl struct{}
)

func GetUserRepo() UserRepo {
	if userRepo == nil {
		userRepo = &userRepoImpl{}
	}
	return userRepo
}

func (ur *userRepoImpl) FindUserBySettionToken(token string) (*model.User, error) {
	u := &model.User{}
	if err := repo.dbMap.SelectOne(u, `
		select
			u.id,
			u.created_at
		from
			users u
			left join sessions s on u.id = s.user_id
		where
			s.token = $1
		`, token); err != nil {
		return nil, err
	}
	return u, nil
}

func (ur *userRepoImpl) FindAuthByUserId(userId string) (*model.Authentication, error) {
	u := &model.Authentication{}
	if err := repo.dbMap.SelectOne(u, `
		select
			*
		from
			authentications
		where
			user_id = $1
		`, userId); err != nil {
		return nil, err
	}
	return u, nil
}

func (ur *userRepoImpl) FindAuthByEmail(email string) (*model.Authentication, error) {
	u := &model.Authentication{}
	if err := repo.dbMap.SelectOne(u, `
		select
			*
		from
			authentications
		where
			email = $1
		`, email); err != nil {
		return nil, err
	}
	return u, nil
}

func (ur *userRepoImpl) FindAuthByToken(token string) (*model.Authentication, error) {
	u := &model.Authentication{}
	if err := repo.dbMap.SelectOne(u, `
		select
			*
		from
			authentications
		where
			token = $1
		`, token); err != nil {
		return nil, err
	}
	return u, nil
}

func (ur *userRepoImpl) FindAssetByUserId(userId string) (*model.Asset, error) {
	u := &model.Asset{}
	if err := repo.dbMap.SelectOne(u, `
		select
			*
		from
			assets
		where
			user_id = $1
		`, userId); err != nil {
		return nil, err
	}
	return u, nil
}
